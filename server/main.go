package main

import (
	"context"
	"net"

	"github.com/danpantry/hello-grpc/protocol"
	"google.golang.org/grpc"
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

func main() {
	lis, err := net.Listen("tcp", ":5051")
	if err != nil {
		panic(err)
	}

	server := grpc.NewServer()
	s := protocol.NewGreeterService(&greetingService{})
	protocol.RegisterGreeterService(server, s)
	server.Serve(lis)
}
