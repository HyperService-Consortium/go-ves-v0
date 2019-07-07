package service

import (
	"encoding/binary"
	"errors"
	"math/rand"

	"golang.org/x/net/context"

	uiprpc "github.com/Myriad-Dreamin/go-ves/grpc"
	types "github.com/Myriad-Dreamin/go-ves/types"
	session "github.com/Myriad-Dreamin/go-ves/types/session"
)

func RequestNSBForNewSession(anyb []byte) ([]byte, error) {
	var buf = make([]byte, 20)
	binary.PutVarint(buf, rand.Int63())
	return buf, nil
}

type SessionStartService = SerialSessionStartService

type SerialSessionStartService struct {
	types.VESDB
	context.Context
	*uiprpc.SessionStartRequest
}

func (s SerialSessionStartService) SessionStart() error {
	var ses = new(session.SerialSession)
	success, help_info, err := ses.InitFromOpIntents(s.GetOpintents())
	if err != nil {
		// TODO: log
		return err
	}
	if !success {
		return errors.New(help_info)
	}
	ses.ISCAddress, err = RequestNSBForNewSession(ses.Content)
	if err != nil {
		return err
	}
	s.InsertSessionInfo(ses)
	// s.UpdateTxs
	// s.UpdateAccs
	return nil
}

func (s SerialSessionStartService) Serve() (*uiprpc.SessionStartReply, error) {
	if err := s.SessionStart(); err != nil {
		return nil, err
	} else {
		return &uiprpc.SessionStartReply{
			Ok: true,
		}, nil
	}
}
