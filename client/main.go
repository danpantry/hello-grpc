package main

import (
	"context"

	"github.com/danpantry/hello-grpc/protocol"
	"google.golang.org/grpc"
)

var serverAddr = "localhost:5051"

func main() {
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	client := protocol.NewGreeterClient(conn)
	greeting, err := client.GetGreeting(context.Background(), &protocol.GreetingParams{})
	if err != nil {
		panic(err)
	}

	println(greeting.GetGreeting())
}
