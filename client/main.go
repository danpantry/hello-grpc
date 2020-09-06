package main

import (
	"context"

	"github.com/danpantry/hello-grpc/protocol"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var serverAddr = "localhost:5051"

func main() {
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	client := protocol.NewGreeterClient(conn)
	md := metadata.New(nil)
	md = protocol.WithJWT(md, "Hello, world!")
	ctx := context.Background()
	ctx = metadata.NewOutgoingContext(ctx, md)
	greeting, err := client.GetGreeting(ctx, &protocol.GreetingParams{})
	if err != nil {
		panic(err)
	}

	println(greeting.GetGreeting())
}
