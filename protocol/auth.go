package protocol

import "google.golang.org/grpc/metadata"

// TODO: Instead of accepting byte arrays, we should accept something more specific
// Like a struct with all the claims, which we then set, or something.
func WithJWT(md metadata.MD, jwt []byte) metadata.MD {
	md2 := metadata.Pairs("jwt", string(jwt))
	return metadata.Join(md, md2)
}

func GetJWT(md metadata.MD) ([]byte, bool) {
	jwts := md.Get("jwt")
	if len(jwts) == 0 {
		return nil, false
	}

	return []byte(jwts[0]), true
}
