package main

import (
	"errors"
	"fmt"

	"github.com/nats-io/nats"
	"github.com/plimble/micro"
	"github.com/plimble/micro/encode/protobuf"
	proto "github.com/plimble/micro/example/proto"
)

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

	m.QueueSubscribe("test", "q", func(ctx *micro.Context) error {
		fmt.Println("start test", ctx.Reply)
		req := &proto.HelloReq{}
		if err := ctx.Decode(ctx.Data, req); err != nil {
			return err
		}

		if req.Name == "error" {
			return errors.New("error name")
		}

		res := &proto.HelloRes{
			Result: "Hello " + req.Name,
		}

		return ctx.Publish(ctx.Reply, res)
	})

	m.RegisterSubscribe()
	m.RegisterQueueSubscribe()

	req := &proto.HelloReq{
		Name: "test",
	}
	res := &proto.HelloRes{}

	err = m.Request("test", req, res, micro.DefaultTimeout)
	if err != nil {
		panic(err)
	}

	fmt.Println(res.Result)

	req = &proto.HelloReq{
		Name: "test",
	}
	res = &proto.HelloRes{}

	err = m.Request("test", req, res, micro.DefaultTimeout)
	if err != nil {
		panic(err)
	}

	fmt.Println(res.Result)
}
