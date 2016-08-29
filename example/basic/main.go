package main

import (
	"fmt"

	"github.com/nats-io/nats"
	"github.com/plimble/micro"
	"github.com/plimble/micro/encode/protobuf"
	proto "github.com/plimble/micro/example/proto"
)

// protoc --go_out=plugins=micro:. *.proto
type helloService struct{}

func (s *helloService) Hello(ctx *micro.Context, req *proto.HelloReq, res *proto.HelloRes) error {
	fmt.Println("start test", ctx.Reply)
	fmt.Println("header", ctx.Header.GetDefault("header1", "default"))
	res.Result = "Hello " + req.Name

	return nil
}

func main() {
	conn, err := nats.Connect("nats://localhost:4222")
	defer conn.Close()
	if err != nil {
		panic(err)
	}

	enc := protobuf.New()

	m := micro.New(conn, enc)

	m.Use(func(ctx *micro.Context) error {
		fmt.Println("start mid")
		err := ctx.Next()
		fmt.Println("end mid")

		return err
	})

	m.HandleError(func(ctx *micro.Context, err error) error {
		fmt.Println(err)
		return err
	})

	hs := &helloService{}

	proto.RegisterHelloService(m, "example", hs)

	m.RegisterSubscribe()
	m.RegisterQueueSubscribe()

	// Client

	client := proto.NewHelloServiceClient("example", m)

	req := &proto.HelloReq{
		Name: "test",
	}

	res, err := client.HelloRequest(req, micro.WithHeader(micro.Header{
		"header1": "h1",
	}))
	if err != nil {
		panic(err)
	}

	fmt.Println(res.Result)

	req = &proto.HelloReq{
		Name: "test",
	}

	res, err = client.HelloRequest(req)
	if err != nil {
		panic(err)
	}

	fmt.Println(res.Result)
}
