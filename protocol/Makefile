hello.pb.go hello_grpc.pb.go: ../schema/hello.proto go.sum go.mod
	bin/protoc -I ../schema ../schema/*.proto --go-grpc_out=. --go_out=.
	cp -r github.com/danpantry/hello-grpc/protocol/* .
	rm -r github.com
	go build

all: hello.pb.go hello_grpc.pb.go
