package main

import (
	"context"
	helloworld "golang_test/grpc/proto_go"
	"google.golang.org/grpc"
	"log"
)

const (
	ADDRESS = "localhost:50001"
)

func main() {
	conn, err := grpc.Dial(ADDRESS, grpc.WithInsecure())

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()

	c := helloworld.NewNihaoClient(conn)
	name := "zhangyl"
	res, err := c.SayHello(context.Background(), &helloworld.HelloRequest{Name: name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	log.Println(res.Message)
}
