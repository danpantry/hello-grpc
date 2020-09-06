module github.com/danpantry/hello-grpc/client

go 1.15

replace github.com/danpantry/hello-grpc/protocol => ../protocol

require (
	github.com/danpantry/hello-grpc/protocol v0.0.0
	github.com/golang/protobuf v1.4.2 // indirect
	google.golang.org/grpc v1.33.0-dev
	google.golang.org/protobuf v1.25.0 // indirect
)
