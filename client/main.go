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

	// TODO: In reality, this JWT would be issued using some other service credentials.
	// For now, we will use a shared secret.
	jwt := os.Getenv("AUTHENTICATION_JWT")
	client := protocol.NewGreeterClient(conn)
	md := metadata.New(nil)
	md = protocol.WithJWT(md, jwt)
	ctx := context.Background()
	ctx = metadata.NewOutgoingContext(ctx, md)
	greeting, err := client.GetGreeting(ctx, &protocol.GreetingParams{})
	if err != nil {
		panic(err)
	}

	println(greeting.GetGreeting())
}
