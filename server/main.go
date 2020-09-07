package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"

	"github.com/danpantry/hello-grpc/protocol"
	"github.com/dgrijalva/jwt-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"gopkg.in/yaml.v2"
)

type permissionSet []string
type callPermissionMap map[string]permissionSet
type servicePermissionMap map[string]callPermissionMap

func (m *servicePermissionMap) AddPermission(service, procedure, permission string) {
	if *m == nil {
		*m = servicePermissionMap{}
	}

	if (*m)[service] == nil {
		(*m)[service] = callPermissionMap{}
	}

	if (*m)[service][procedure] == nil {
		(*m)[service][procedure] = permissionSet{}
	}

	(*m)[service][procedure] = append((*m)[service][procedure], permission)
}

func (m servicePermissionMap) GetEntry(service, procedure string) []string {
	s, ok := m[service]
	if !ok {
		return nil
	}

	entry, ok := s[procedure]
	return entry
}

func (m servicePermissionMap) HasEntry(service, procedure string) bool {
	s, ok := m[service]
	if !ok {
		return false
	}

	_, ok = s[procedure]
	return ok
}

// Validate ensures that all service/rpc calls within the given server have a corresponding permissions entry.
// This is useful to make sure that no calls are blocked by accident.
//
// Services/Calls without corresponding permissions will be written to the log.
// A better implementation would return them.
func (m servicePermissionMap) Validate(s *grpc.Server) {
	type pair struct {
		service string
		call    string
	}

	uninitialized := []pair{}
	for name, value := range s.GetServiceInfo() {
		for _, method := range value.Methods {
			if !m.HasEntry(name, method.Name) {
				uninitialized = append(uninitialized, pair{name, method.Name})
			}
		}
	}

	for _, p := range uninitialized {
		log.Printf("%s.%s has no permissions - all calls will be blocked!\n", p.service, p.call)
	}
}

// CanExecute returns true if the given claim set has permission to execute the given call.
func (m servicePermissionMap) CanExecute(claims protocol.Claims, info *grpc.UnaryServerInfo) bool {
	// If the permission set gets too large this could be a problem and we should use a map instead
	// The first character will be a / so skip that.
	// Claims is signed and checked by us but it still comes from the user so be wary.
	ident := strings.Split(info.FullMethod[1:], "/")
	for _, permission := range m.GetEntry(ident[0], ident[1]) {
		if claims.HasPermission(permission) {
			return true
		}
	}

	return false
}

type greetingService struct{}

func (*greetingService) GetGreeting(ctx context.Context, params *protocol.GreetingParams) (*protocol.Greeting, error) {
	greeting := "Hello, world!"
	msg := protocol.Greeting{
		Greeting: &greeting,
	}

	return &msg, nil
}

var (
	errUnauthenticated = fmt.Errorf("unauthenticated")
	errUnauthorized    = fmt.Errorf("unauthorized")
)

func authInterceptor(kf jwt.Keyfunc, permissions servicePermissionMap) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// Verify that the use has a signed JWT
		// You would also want to do authorization here by checking that the claims include permission to execute the required handler
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, errUnauthenticated
		}

		var claims protocol.Claims
		_, err := protocol.ParseJWTWithClaims(md, &claims, kf)
		if err != nil {
			if err == protocol.ErrNotPresent {
				err = errUnauthenticated
			}

			return nil, err
		}

		if !permissions.CanExecute(claims, info) {
			return nil, errUnauthorized
		}

		return handler(ctx, req)
	}
}

func loadPermissions(p *servicePermissionMap) error {
	b, err := ioutil.ReadFile("claims.yml")
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(b, p); err != nil {
		// An error may be returned but the YAML fail may still parse otherwise.
		// Not sure how we want to handle this, for now I suggest we just bomb out.
		// One could check the returned error against yaml.TypeError.
		return fmt.Errorf("error parsing claims.yml: %w", err)
	}

	return nil
}

var exitCodeConfigError = 0x1

func hs256(secret []byte) jwt.Keyfunc {
	return func(tok *jwt.Token) (interface{}, error) {
		if tok.Method.Alg() != jwt.SigningMethodHS256.Name {
			return nil, jwt.ErrInvalidKey
		}

		return secret, nil
	}
}

func main() {
	var permissions servicePermissionMap
	if err := loadPermissions(&permissions); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(exitCodeConfigError)
	}

	// In reality, you would probably retrieve the signing secret from some asynchronous source not stored with each instance of the gRPC server.
	kf := hs256([]byte(os.Getenv("JWT_SIGNING_SECRET")))
	interceptor := authInterceptor(kf, permissions)
	server := grpc.NewServer(grpc.UnaryInterceptor(interceptor))
	s := protocol.NewGreeterService(&greetingService{})
	protocol.RegisterGreeterService(server, s)
	permissions.Validate(server)

	lis, err := net.Listen("tcp", ":5051")
	if err != nil {
		panic(err)
	}

	server.Serve(lis)
}
