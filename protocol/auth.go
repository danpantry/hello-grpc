package protocol

import (
	"google.golang.org/grpc/metadata"
)

// WithJWT adds the given JWT to the correct Metadata key to be parsed by the server side.
func WithJWT(md metadata.MD, jwt string) metadata.MD {
	md2 := metadata.Pairs("jwt", string(jwt))
	return metadata.Join(md, md2)
}

// GetJWT returns a JWT encoded in the given metadata and a flag indicating if the JWT was present.
func GetJWT(md metadata.MD) (string, bool) {
	var s string
	jwts := md.Get("jwt")
	if len(jwts) == 0 {
		return s, false
	}

	return jwts[0], true
}
