package service

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"golang.org/x/net/context"

	uiptypes "github.com/HyperService-Consortium/go-uip/uiptypes"
	uiprpc "github.com/HyperService-Consortium/go-ves/grpc/uiprpc"
	uipbase "github.com/HyperService-Consortium/go-ves/grpc/uiprpc-base"
	nsbcli "github.com/HyperService-Consortium/go-ves/lib/net/nsb-client"
	types "github.com/HyperService-Consortium/go-ves/types"
	session "github.com/HyperService-Consortium/go-ves/types/session"
)

type SessionStartService = MultiThreadSerialSessionStartService

type SerialSessionStartService struct {
	Signer    uiptypes.Signer
	NsbClient *nsbcli.NSBClient
	CVes      uiprpc.CenteredVESClient
	types.VESDB
	context.Context
	*uiprpc.SessionStartRequest
}

func (s *SerialSessionStartService) RequestNSBForNewSession(anyb types.Session) ([]byte, error) {
	var accs = anyb.GetAccounts()

	var owners = make([][]byte, 0, len(accs)+1)
	// todo
	// owners = append(owners, s.Signer.GetPublicKey())
	for _, owner := range accs {
		owners = append(owners, owner.GetAddress())
		fmt.Println("waiting", hex.EncodeToString(owner.GetAddress()))
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
	// fmt.Println("accs, txs", owners, txs)
	return s.NsbClient.CreateISC(s.Signer, make([]uint32, len(owners)), owners, txs, s.Signer.Sign(bytes.Join(anyb.GetTransactions(), []byte{})).Bytes())
}

func (s *SerialSessionStartService) SessionStart() ([]byte, []uiptypes.Account, error) {
	var ses = new(session.SerialSession)
	ses.Signer = s.Signer
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
		fmt.Println(s.NsbClient.FreezeInfo(s.Signer, ses.ISCAddress, uint64(i)))
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
		// fmt.Println("reply?", r, err)
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
	Signer    uiptypes.Signer
	NsbClient *nsbcli.NSBClient
	CVes      uiprpc.CenteredVESClient
	types.VESDB
	context.Context
	*uiprpc.SessionStartRequest
}

func (s *MultiThreadSerialSessionStartService) RequestNSBForNewSession(anyb types.Session) ([]byte, error) {
	var accs = anyb.GetAccounts()

	var owners = make([][]byte, 0, len(accs)+1)
	// todo
	// owners = append(owners, s.Signer.GetPublicKey())
	for _, owner := range accs {
		owners = append(owners, owner.GetAddress())
		fmt.Println("waiting", hex.EncodeToString(owner.GetAddress()))
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
	// fmt.Println("accs, txs", owners, txs)
	return s.NsbClient.CreateISC(s.Signer, make([]uint32, len(owners)), owners, txs, s.Signer.Sign(bytes.Join(anyb.GetTransactions(), []byte{})).Bytes())
}

func (s *MultiThreadSerialSessionStartService) SessionStart() ([]byte, []uiptypes.Account, error) {
	var ses = new(session.MultiThreadSerialSession)
	ses.Signer = s.Signer
	success, help_info, err := ses.InitFromOpIntents(s.GetOpintents())
	if err != nil {
		// TODO: log
		return nil, nil, err
	}
	if !success {
		return nil, nil, errors.New(help_info)
	}
	ses.ISCAddress, err = s.RequestNSBForNewSession(ses)
	if ses.ISCAddress == nil {
		return nil, nil, fmt.Errorf("request isc failed: %v", err)
	}
	if err != nil {
		return nil, nil, fmt.Errorf("request isc failed on request: %v", err)
	}
	err = ses.AfterInitGUID()
	fmt.Println("after init guid...", ses.ISCAddress, hex.EncodeToString(ses.ISCAddress))
	if err != nil {
		return nil, nil, err
	}

	err = s.InsertSessionInfo(ses)
	if err != nil {
		return nil, nil, err
	}
	for i := uint32(0); i < ses.TransactionCount; i++ {
		fmt.Println(s.NsbClient.FreezeInfo(s.Signer, ses.ISCAddress, uint64(i)))
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
		// fmt.Println("reply?", r, err)
		if err != nil {
			return nil, err
		}

		return &uiprpc.SessionStartReply{
			Ok:        r.GetOk(),
			SessionId: b,
		}, nil
	}
}
