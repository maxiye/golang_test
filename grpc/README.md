# grpc
_安装组件_
```
go get -u google.golang.org/grpc
go get -u github.com/golang/protobuf/protoc-gen-go
go mod vendor
```
_protoc编译器(根据proto文件生成gRPC服务代码)_

下载地址：[https://github.com/protocolbuffers/protobuf/releases][1]
运行路径放到path环境变量

_新建helloworld.proto_
```
syntax = "proto3";

option objc_class_prefix = "HLW";

package helloworld;

// The greeting service definition.
service Nihao {
  // Sends a greeting
  rpc SayHello (HelloRequest) returns (HelloReply) {}
}

// The request message containing the user's name.
message HelloRequest {
  string name = 1;
}

// The response message containing the greetings
message HelloReply {
  string message = 1;
}
```
_编译proto为go代码_
```
cd golang_test/grpc/proto
protoc * --go_out=plugins=grpc:../proto_go
```

_编写服务器代码_
```
package main

import (
	"context"
	helloworld "golang_test/grpc/proto_go"
	"google.golang.org/grpc"
	"log"
	"net"
)

const (
	PORT = ":50001"
)

type server struct {

}

func (s *server) SayHello(ctx context.Context, in *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
	return &helloworld.HelloReply{Message:"Hello " + in.Name}, nil
}

func main() {
	lis, err := net.Listen("tcp", PORT)
	if  err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	helloworld.RegisterNihaoServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

```

_编写客户端代码_
```
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
	r, err := c.SayHello(context.Background(), &helloworld.HelloRequest{Name:name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	log.Println(r.Message)
}

```

_运行服务端代码_
```
go run grpc\service\hello_server.go
```


_运行客户端代码_
```
go run grpc\client\hello_client.go
```

[1]: https://github.com/protocolbuffers/protobuf/releases