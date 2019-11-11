package service

import (
	"golang.org/x/net/context"

	uiprpc "github.com/HyperService-Consortium/go-ves/grpc/uiprpc"
	types "github.com/HyperService-Consortium/go-ves/types"
)

type SessionRequireTransactService struct {
	types.VESDB
	context.Context
	*uiprpc.SessionRequireTransactRequest
}

func (s SessionRequireTransactService) Serve() (*uiprpc.SessionRequireTransactReply, error) {
	// todo errors.New("TODO")
	s.ActivateSession(s.GetSessionId())
	defer s.InactivateSession(s.GetSessionId())
	return &uiprpc.SessionRequireTransactReply{
		// Tx: true,
	}, nil
}
