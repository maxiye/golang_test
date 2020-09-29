package main

import (
	"context"
	"encoding/json"
	api "golang_test/grpc/proto_go"
	"google.golang.org/grpc"
	"io/ioutil"
	"log"
	"net"
	"net/http"
)

type ApiServer struct {
}

func (a ApiServer) GetModTrends(_ context.Context, req *api.TrendRequest) (*api.TrendReply, error) {
	period := req.Period
	log.Println("https://goproxy.cn/stats/trends/" + period)
	var itemList []*api.TrendReply_ModItem
	if res, err := http.Get("https://goproxy.cn/stats/trends/" + period); err == nil && res != nil {
		defer func() {
			_ = res.Body.Close()
		}()
		body, _ := ioutil.ReadAll(res.Body)
		log.Println(string(body))
		if err = json.Unmarshal(body, &itemList); err == nil {
			return &api.TrendReply{Period: period, ModItem: itemList}, nil
		}
	} else {
		log.Println("err: ", err, res)
	}

	return &api.TrendReply{Period: period, ModItem: itemList}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50002")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	api.RegisterApiServer(s, &ApiServer{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
