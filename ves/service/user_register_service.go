package service

import (
	"github.com/HyperService-Consortium/go-ves/ves/vs"
	"golang.org/x/net/context"

	"github.com/HyperService-Consortium/go-ves/grpc/uiprpc"
)

type UserRegisterService struct {
	*vs.VServer
	context.Context
	*uiprpc.UserRegisterRequest
}

func NewUserRegisterService(server *vs.VServer, context context.Context, userRegisterRequest *uiprpc.UserRegisterRequest) UserRegisterService {
	return UserRegisterService{VServer: server, Context: context, UserRegisterRequest: userRegisterRequest}
}

func (s UserRegisterService) Serve() (*uiprpc.UserRegisterReply, error) {
	if err := s.DB.InsertAccount(s.GetUserName(), s.GetAccount()); err != nil {
		s.Logger.Error("error", "error", err)
		return nil, err
	} else {
		return &uiprpc.UserRegisterReply{
			Ok: true,
		}, nil
	}
}
