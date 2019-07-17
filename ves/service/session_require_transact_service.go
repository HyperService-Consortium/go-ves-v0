package service

import (
	"errors"

	"golang.org/x/net/context"

	uiprpc "github.com/Myriad-Dreamin/go-ves/grpc/uip-rpc"
	types "github.com/Myriad-Dreamin/go-ves/types"
)

type SessionRequireTransactService struct {
	types.VESDB
	context.Context
	*uiprpc.SessionRequireTransactRequest
}

func (s SessionRequireTransactService) Serve() (*uiprpc.SessionRequireTransactReply, error) {
	if err := errors.New("TODO"); err != nil {
		return nil, err
	} else {
		return &uiprpc.SessionRequireTransactReply{
			Ok: true,
		}, nil
	}
}
