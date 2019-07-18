package service

import (
	"errors"
	"fmt"

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
	ses, err := s.FindSessionInfo(s.SessionId)

	// todo: get Session Acked from isc
	// nsbClient.

	if err == nil {
		var success bool
		var help_info string
		success, help_info, err = ses.AckForInit(s.GetUser(), s.GetUserSignature())
		if err != nil {
			// todo, log
			return nil, fmt.Errorf("internal error: %v", err)
		} else if !success {
			return nil, errors.New(help_info)
		} else {

			return &uiprpc.SessionAckForInitReply{
				Ok: true,
			}, nil
		}
	} else {
		return nil, err
	}
}
