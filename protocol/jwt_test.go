package protocol

import (
	"testing"

	"github.com/dgrijalva/jwt-go"
	"google.golang.org/grpc/metadata"
)

func TestParsesClaimsCorrectly(t *testing.T) {
	str := `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyLCJwZXJtaXNzaW9ucyI6WyJncmVldGVyIl19.1BqKP_6iHZ6h1hVHZDRZo2KkqADN9S3VdbtFr6Y9-3Q`
	md := WithJWT(metadata.New(nil), str)

	var claims Claims
	_, err := ParseJWTWithClaims(md, &claims, func(*jwt.Token) (interface{}, error) {
		return []byte("ninjas"), nil
	})

	if err != nil {
		t.Fatal(err)
	}

	if claims.HasPermission("greeter") == false {
		t.Error("HasPermission(): expected true, received false")
	}
}
