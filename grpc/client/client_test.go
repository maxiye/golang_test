package main

import (
	"context"
	grpcService "golang_test/grpc/proto_go"
	"google.golang.org/grpc"
	"log"
	"testing"
)

const (
	ADDRESS = "localhost:50001"
)

func TestNihao(t *testing.T) {
	conn, err := grpc.Dial(ADDRESS, grpc.WithInsecure())

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer func() {
		_ = conn.Close()
	}()

	c := grpcService.NewNihaoClient(conn)
	name := "zhangyl"
	res, err := c.SayHello(context.Background(), &grpcService.HelloRequest{Name: name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	log.Println(res.Message)
}

func TestApi(t *testing.T) {
	conn, err := grpc.Dial("localhost:50002", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer func() {
		_ = conn.Close()
	}()

	c := grpcService.NewApiClient(conn)

	if res, err := c.GetModTrends(context.Background(), &grpcService.TrendRequest{Period: "latest"}); err == nil {
		t.Log(res.ModItem[0])
	}
	if res, err := c.GetModTrends(context.Background(), &grpcService.TrendRequest{Period: "last-7-days"}); err == nil {
		t.Log(res.ModItem[0])
	}
	if res, err := c.GetModTrends(context.Background(), &grpcService.TrendRequest{Period: "last-30-days"}); err == nil {
		t.Log(res.ModItem)
	}

}
