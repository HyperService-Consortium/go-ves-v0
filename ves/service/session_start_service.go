package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	"golang.org/x/net/context"

	uiptypes "github.com/Myriad-Dreamin/go-uip/types"
	uiprpc "github.com/Myriad-Dreamin/go-ves/grpc"
	nsbcli "github.com/Myriad-Dreamin/go-ves/net/nsb_client"
	types "github.com/Myriad-Dreamin/go-ves/types"
	session "github.com/Myriad-Dreamin/go-ves/types/session"
)

var nsbClient = nsbcli.NewNSBClient("47.251.2.73:26657")

type SessionStartService = SerialSessionStartService

type SerialSessionStartService struct {
	Signer uiptypes.Signer
	types.VESDB
	context.Context
	*uiprpc.SessionStartRequest
}

func (s *SerialSessionStartService) RequestNSBForNewSession(anyb types.Session) ([]byte, error) {
	var accs = anyb.GetAccounts()

	var owners = make([][]byte, 0, len(accs)+1)
	owners = append(owners, s.Signer.GetPublicKey())
	for _, owner := range accs {
		owners = append(owners, owner.GetAddress())
	}
	var txs = anyb.GetTransactions()
	var btxs = make([][]byte, 0, len(txs))
	for _, tx := range txs {
		b, err := json.Marshal(tx)
		if err != nil {
			return nil, err
		}
		btxs = append(btxs, b)
	}
	fmt.Println("accs", owners)
	return nsbClient.CreateISC(s.Signer, make([]uint32, len(owners)), owners, nil, s.Signer.Sign(bytes.Join(anyb.GetTransactions(), []byte{})))
}
func (s *SerialSessionStartService) SessionStart() ([]byte, error) {
	var ses = new(session.SerialSession)
	success, help_info, err := ses.InitFromOpIntents(s.GetOpintents())
	if err != nil {
		// TODO: log
		return nil, err
	}
	if !success {
		return nil, errors.New(help_info)
	}
	ses.ISCAddress, err = s.RequestNSBForNewSession(ses)
	if err != nil {
		return nil, err
	}
	s.InsertSessionInfo(ses)

	// s.UpdateTxs
	// s.UpdateAccs
	return ses.ISCAddress, nil
}

func (s *SerialSessionStartService) Serve() (*uiprpc.SessionStartReply, error) {
	if b, err := s.SessionStart(); err != nil {
		return nil, err
	} else {
		return &uiprpc.SessionStartReply{
			Ok:        true,
			SessionId: b,
		}, nil
	}
}
