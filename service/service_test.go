package service

import (
	"errors"
	"fmt"
	"testing"
)

type CalService struct{}

type AddReq struct {
	X int
}
type AddRes struct{}

func (s *CalService) Add(req *AddReq, res *AddRes) error {
	fmt.Println("call add")
	return errors.New("get error")
}

type SubReq struct{}
type SubRes struct{}

func (s *CalService) Sub(req *SubReq, res *SubRes) error {
	fmt.Println("call sub")
	return nil
}

func TestService(t *testing.T) {
	QueueService(nil, "com.test", &CalService{})
}
