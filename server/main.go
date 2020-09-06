package main

import (
	"context"
	"net"

	"github.com/danpantry/hello-grpc/protocol"
	"google.golang.org/grpc"
)

func new() protocol.GreeterService {
	server := protocol.GreeterService{}
	server.GetGreeting = func(ctx context.Context, params *protocol.GreetingParams) (*protocol.Greeting, error) {
		greeting := "Hello, world!"
		msg := protocol.Greeting{
			Greeting: &greeting,
		}

		return &msg, nil
	}

	return server
}

func main() {
	lis, err := net.Listen("tcp", ":5051")
	if err != nil {
		panic(err)
	}

	server := grpc.NewServer()
	s := new()
	protocol.RegisterGreeterService(server, &s)
	server.Serve(lis)
}
