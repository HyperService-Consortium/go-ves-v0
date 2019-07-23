package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"golang.org/x/net/context"

	uiptypes "github.com/Myriad-Dreamin/go-uip/types"
	uiprpc "github.com/Myriad-Dreamin/go-ves/grpc/uiprpc"
	uipbase "github.com/Myriad-Dreamin/go-ves/grpc/uiprpc-base"
	nsbcli "github.com/Myriad-Dreamin/go-ves/net/nsb_client"
	types "github.com/Myriad-Dreamin/go-ves/types"
	session "github.com/Myriad-Dreamin/go-ves/types/session"
)

var nsbClient = nsbcli.NewNSBClient("47.251.2.73:26657")

type SessionStartService = MultiThreadSerialSessionStartService

type SerialSessionStartService struct {
	Signer uiptypes.Signer
	CVes   uiprpc.CenteredVESClient
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
	fmt.Println("accs, txs", owners, txs)
	return nsbClient.CreateISC(s.Signer, make([]uint32, len(owners)), owners, txs, s.Signer.Sign(bytes.Join(anyb.GetTransactions(), []byte{})))
}

func (s *SerialSessionStartService) SessionStart() ([]byte, []uiptypes.Account, error) {
	var ses = new(session.SerialSession)
	success, help_info, err := ses.InitFromOpIntents(s.GetOpintents())
	if err != nil {
		// TODO: log
		return nil, nil, err
	}
	if !success {
		return nil, nil, errors.New(help_info)
	}
	ses.ISCAddress, err = s.RequestNSBForNewSession(ses)
	if err != nil {
		return nil, nil, err
	}
	s.InsertSessionInfo(ses)
	for i := uint32(0); i < ses.TransactionCount; i++ {
		fmt.Println(nsbClient.FreezeInfo(s.Signer, ses.ISCAddress, uint64(i)))
	}
	// s.UpdateTxs
	// s.UpdateAccs
	return ses.ISCAddress, ses.GetAccounts(), nil
}

func (s *SerialSessionStartService) Serve() (*uiprpc.SessionStartReply, error) {
	if b, accs, err := s.SessionStart(); err != nil {
		return nil, err
	} else {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		r, err := s.CVes.InternalRequestComing(ctx, &uiprpc.InternalRequestComingRequest{
			SessionId: b,
			Host:      []byte{127, 0, 0, 1, ((23351) >> 8 & 0xff), 23351 & 0xff},
			Accounts: func() (uaccs []*uipbase.Account) {
				for _, acc := range accs {
					uaccs = append(uaccs, &uipbase.Account{
						Address: acc.GetAddress(),
						ChainId: acc.GetChainId(),
					})
				}
				return
			}(),
		})
		fmt.Println("reply?", r, err)
		if err != nil {
			return nil, err
		}

		return &uiprpc.SessionStartReply{
			Ok:        r.GetOk(),
			SessionId: b,
		}, nil
	}
}

type MultiThreadSerialSessionStartService struct {
	Signer uiptypes.Signer
	CVes   uiprpc.CenteredVESClient
	types.VESDB
	context.Context
	*uiprpc.SessionStartRequest
}

func (s *MultiThreadSerialSessionStartService) RequestNSBForNewSession(anyb types.Session) ([]byte, error) {
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
	fmt.Println("accs, txs", owners, txs)
	return nsbClient.CreateISC(s.Signer, make([]uint32, len(owners)), owners, txs, s.Signer.Sign(bytes.Join(anyb.GetTransactions(), []byte{})))
}

func (s *MultiThreadSerialSessionStartService) SessionStart() ([]byte, []uiptypes.Account, error) {
	var ses = new(session.MultiThreadSerialSession)
	success, help_info, err := ses.InitFromOpIntents(s.GetOpintents())
	if err != nil {
		// TODO: log
		return nil, nil, err
	}
	if !success {
		return nil, nil, errors.New(help_info)
	}
	ses.ISCAddress, err = s.RequestNSBForNewSession(ses)
	if err != nil {
		return nil, nil, err
	}
	err = ses.AfterInitGUID()
	if err != nil {
		return nil, nil, err
	}

	s.InsertSessionInfo(ses)
	for i := uint32(0); i < ses.TransactionCount; i++ {
		fmt.Println(nsbClient.FreezeInfo(s.Signer, ses.ISCAddress, uint64(i)))
	}
	// s.UpdateTxs
	// s.UpdateAccs
	return ses.ISCAddress, ses.GetAccounts(), nil
}

func (s *MultiThreadSerialSessionStartService) Serve() (*uiprpc.SessionStartReply, error) {
	if b, accs, err := s.SessionStart(); err != nil {
		return nil, err
	} else {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		r, err := s.CVes.InternalRequestComing(ctx, &uiprpc.InternalRequestComingRequest{
			SessionId: b,
			Host:      []byte{127, 0, 0, 1, ((23351) >> 8 & 0xff), 23351 & 0xff},
			Accounts: func() (uaccs []*uipbase.Account) {
				for _, acc := range accs {
					uaccs = append(uaccs, &uipbase.Account{
						Address: acc.GetAddress(),
						ChainId: acc.GetChainId(),
					})
				}
				return
			}(),
		})
		fmt.Println("reply?", r, err)
		if err != nil {
			return nil, err
		}

		return &uiprpc.SessionStartReply{
			Ok:        r.GetOk(),
			SessionId: b,
		}, nil
	}
}
