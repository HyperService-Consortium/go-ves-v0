package service

import (
	"golang.org/x/net/context"

	uiprpc "github.com/Myriad-Dreamin/go-ves/grpc"
	types "github.com/Myriad-Dreamin/go-ves/types"
)

type UserRegisterService struct {
	types.VESDB
	context.Context
	*uiprpc.UserRegisterRequest
}

func (s UserRegisterService) Serve() (*uiprpc.UserRegisterReply, error) {
	if err := s.InsertAccount(s.GetUserName(), s.GetAccount()); err != nil {
		return nil, err
	} else {
		return &uiprpc.UserRegisterReply{
			Ok: true,
		}, nil
	}
}
