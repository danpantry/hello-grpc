package protocol

import (
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"google.golang.org/grpc/metadata"
)

// WithJWT adds the given JWT to the correct Metadata key to be parsed by the server side.
func WithJWT(md metadata.MD, jwt string) metadata.MD {
	md2 := metadata.Pairs("jwt", string(jwt))
	return metadata.Join(md, md2)
}

// GetJWT returns a JWT encoded in the given metadata and a flag indicating if the JWT was present.
func getJWT(md metadata.MD) (string, bool) {
	var s string
	jwts := md.Get("jwt")
	if len(jwts) == 0 {
		return s, false
	}

	return jwts[0], true
}

// Claims is the set of claims encoded within an authorization JWT.
type Claims struct {
	Permissions []string `json:"permissions"`
}

// HasPermission indicates whether the claim set has the given permission.
func (c *Claims) HasPermission(perm string) bool {
	// Our usage of permissions as being a []string will make this slower as this will probably be called in a loop!
	// This is fine as long as the claim permission set and the set of permissions required for calls stays low.
	for _, permission := range c.Permissions {
		if permission == perm {
			return true
		}
	}

	return false
}

// Valid indicates whether or not the claim set is valid.
func (c *Claims) Valid() error {
	return nil
}

// ErrNotPresent is returned when the authentication JWT is not included in the gRPC metadata.
var ErrNotPresent = fmt.Errorf("jwt not present in metadata")

// ParseJWTWithClaims extracts a JWT from the given Metadata.
// Refer to github.com/dgrijalva/jwt-go for information on the kf parameter.
func ParseJWTWithClaims(md metadata.MD, c *Claims, kf jwt.Keyfunc) (*jwt.Token, error) {
	str, ok := getJWT(md)
	if !ok {
		return nil, ErrNotPresent
	}

	return jwt.ParseWithClaims(str, c, kf)
}
