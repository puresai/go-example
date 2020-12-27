package endpoint

import (
	"context"
	"demo6/service"
	"fmt"

	"github.com/go-kit/kit/endpoint"
)

type Request struct {
	A int `json:"a" form:"a"`
	B int `json:"b" form:"b"`
}

type Res struct {
	Res int   `json:"res"`
	Err error `json:"err"`
}

func MakeAddEndpoint(s service.CalculateService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(Request)
		return Res{Res: s.Add(req.A, req.B)}, nil
	}
}

func MakeReduceEndpoint(s service.CalculateService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(Request)
		fmt.Println("ajfskja")
		return Res{Res: s.Reduce(req.A, req.B)}, nil
	}
}

func MakeMultiEndpoint(s service.CalculateService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(Request)
		fmt.Println("ajfskja")
		return Res{Res: s.Multi(req.A, req.B)}, nil
	}
}
