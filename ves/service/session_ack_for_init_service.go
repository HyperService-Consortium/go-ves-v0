package service

import (
	"errors"

	"golang.org/x/net/context"

	uiprpc "github.com/Myriad-Dreamin/go-ves/grpc/uip-rpc"
	types "github.com/Myriad-Dreamin/go-ves/types"
)

type SessionAckForInitService struct {
	types.VESDB
	context.Context
	*uiprpc.SessionAckForInitRequest
}

func (s SessionAckForInitService) Serve() (*uiprpc.SessionAckForInitReply, error) {
	if err := errors.New("TODO"); err != nil {
		return nil, err
	} else {
		return &uiprpc.SessionAckForInitReply{
			Ok: true,
		}, nil
	}
}
