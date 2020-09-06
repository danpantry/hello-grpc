package main

import (
	"context"
	"fmt"
	"net"

	"github.com/danpantry/hello-grpc/protocol"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type greetingService struct {
}

func (*greetingService) GetGreeting(ctx context.Context, params *protocol.GreetingParams) (*protocol.Greeting, error) {
	greeting := "Hello, world!"
	msg := protocol.Greeting{
		Greeting: &greeting,
	}

	return &msg, nil
}

var errUnauthenticated = fmt.Errorf("unauthenticated")

func authInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// Verify that the use has a signed JWT
	// You would also want to do authorization here by checking that the claims include permission to execute the required handler
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errUnauthenticated
	}

	_, ok = protocol.GetJWT(md)
	if !ok {
		return nil, errUnauthenticated
	}

	// TODO: check jwt
	return handler(ctx, req)
}

func main() {
	lis, err := net.Listen("tcp", ":5051")
	if err != nil {
		panic(err)
	}

	server := grpc.NewServer(grpc.UnaryInterceptor(authInterceptor))
	s := protocol.NewGreeterService(&greetingService{})
	protocol.RegisterGreeterService(server, s)
	server.Serve(lis)
}
