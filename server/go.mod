module github.com/danpantry/hello-grpc/server

go 1.15

replace github.com/danpantry/hello-grpc/protocol => ../protocol

require (
	github.com/danpantry/hello-grpc/protocol v0.0.0
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	google.golang.org/grpc v1.33.0-dev
	gopkg.in/yaml.v2 v2.3.0
)
