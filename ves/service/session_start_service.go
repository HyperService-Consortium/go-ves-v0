package service

import (
	"bytes"
	"encoding/json"
	"errors"

	"golang.org/x/net/context"

	uiprpc "github.com/Myriad-Dreamin/go-ves/grpc"
	nsbcli "github.com/Myriad-Dreamin/go-ves/net/nsb_client"
	types "github.com/Myriad-Dreamin/go-ves/types"
	session "github.com/Myriad-Dreamin/go-ves/types/session"
)

var nsbClient = nsbcli.NewNSBClient("47.251.2.73:26657")

type SessionStartService = SerialSessionStartService

type SerialSessionStartService struct {
	signer types.TenSigner
	types.VESDB
	context.Context
	*uiprpc.SessionStartRequest
}

func (s *SerialSessionStartService) RequestNSBForNewSession(anyb types.Session) ([]byte, error) {
	var accs = anyb.GetAccounts()
	var owners = make([][]byte, 0, len(accs)+1)
	owners = append(owners, s.signer.GetPublicKey())
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
	return nsbClient.CreateISC(s.signer, []uint32{0, 0, 0}, owners, nil, s.signer.Sign(bytes.Join(anyb.GetTransactions(), []byte{})))
}
func (s *SerialSessionStartService) SessionStart() error {
	var ses = new(session.SerialSession)
	success, help_info, err := ses.InitFromOpIntents(s.GetOpintents())
	if err != nil {
		// TODO: log
		return err
	}
	if !success {
		return errors.New(help_info)
	}
	ses.ISCAddress, err = s.RequestNSBForNewSession(ses)
	if err != nil {
		return err
	}
	s.InsertSessionInfo(ses)

	// s.UpdateTxs
	// s.UpdateAccs
	return nil
}

func (s *SerialSessionStartService) Serve() (*uiprpc.SessionStartReply, error) {
	if err := s.SessionStart(); err != nil {
		return nil, err
	} else {
		return &uiprpc.SessionStartReply{
			Ok: true,
		}, nil
	}
}
