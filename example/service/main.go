package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/nats-io/nats"
	"github.com/plimble/micro"
	"github.com/plimble/micro/encode/json"
	"github.com/plimble/micro/service"
)

type CalService struct{}

type AddReq struct {
	X int
	Y int
}
type AddRes struct {
	Result int
}

func (s *CalService) Add(req *AddReq, res *AddRes) error {
	fmt.Println("call add")
	return errors.New("get error")
}

type SubReq struct {
	X int
	Y int
}
type SubRes struct {
	Result int
}

func (s *CalService) Sub(req *SubReq, res *SubRes) error {
	fmt.Println("call sub")
	res.Result = req.X - req.Y
	return nil
}

func main() {
	conn, err := nats.Connect("nats://localhost:4222")
	defer conn.Close()
	if err != nil {
		panic(err)
	}

	enc := json.New()

	m := micro.New(conn, enc)

	m.HandleError(func(ctx *micro.Context, err error) error {
		fmt.Println(err)
		return err
	})

	calService := &CalService{}
	service.QueueSubscribe(m, "example", calService)

	m.RegisterSubscribe()
	m.RegisterQueueSubscribe()

	time.Sleep(1 * time.Second)

	req := &SubReq{10, 2}
	res := &SubRes{}

	err = m.Request("example.Sub", req, res, micro.DefaultTimeout)
	if err != nil {
		panic(err)
	}

	fmt.Println(res.Result)
}
