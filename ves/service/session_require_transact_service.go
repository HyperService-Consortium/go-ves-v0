package service

import (
	"golang.org/x/net/context"

	uiprpc "github.com/Myriad-Dreamin/go-ves/grpc/uiprpc"
	types "github.com/Myriad-Dreamin/go-ves/types"
)

type SessionRequireTransactService struct {
	types.VESDB
	context.Context
	*uiprpc.SessionRequireTransactRequest
}

func (s SessionRequireTransactService) Serve() (*uiprpc.SessionRequireTransactReply, error) {
	// todo errors.New("TODO")
	var err error
	if err != nil {
		return nil, err
	} else {
		return &uiprpc.SessionRequireTransactReply{
			Ok: true,
		}, nil
	}
}
