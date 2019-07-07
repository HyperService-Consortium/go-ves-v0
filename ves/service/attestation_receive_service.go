package service

import (
	"errors"

	"golang.org/x/net/context"

	uiprpc "github.com/Myriad-Dreamin/go-ves/grpc"
	types "github.com/Myriad-Dreamin/go-ves/types"
)

type AttestationReceiveService struct {
	types.VESDB
	context.Context
	*uiprpc.AttestationReceiveRequest
}

func (s AttestationReceiveService) Serve() (*uiprpc.AttestationReceiveReply, error) {
	if err := errors.New("TODO"); err != nil {
		return nil, err
	} else {
		return &uiprpc.AttestationReceiveReply{
			Ok: true,
		}, nil
	}
}
