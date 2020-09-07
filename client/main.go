package main

import (
	"context"
	"os"

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

	jwt := os.Getenv("AUTHENTICATION_JWT")
	client := protocol.NewGreeterClient(conn)
	md := protocol.WithJWT(metadata.New(nil), jwt)
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	greeting, err := client.GetGreeting(ctx, &protocol.GreetingParams{})
	if err != nil {
		panic(err)
	}

	println(greeting.GetGreeting())
}
