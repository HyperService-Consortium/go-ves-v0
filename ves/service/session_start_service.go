package service

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	logger "github.com/HyperService-Consortium/go-ves/lib/log"
	"github.com/HyperService-Consortium/go-ves/ves/vs"
	"time"

	"golang.org/x/net/context"

	"github.com/HyperService-Consortium/go-uip/uiptypes"
	"github.com/HyperService-Consortium/go-ves/grpc/uiprpc"
	uipbase "github.com/HyperService-Consortium/go-ves/grpc/uiprpc-base"
	"github.com/HyperService-Consortium/go-ves/types"
	"github.com/HyperService-Consortium/go-ves/types/session"
)

type SessionStartService = MultiThreadSerialSessionStartService

func NewSessionStartService(server *vs.VServer, context context.Context, sessionStartRequest *uiprpc.SessionStartRequest) SessionStartService {
	return NewMultiThreadSerialSessionStartService(server, context, sessionStartRequest)
}

type SerialSessionStartService struct {
	*vs.VServer
	context.Context
	*uiprpc.SessionStartRequest
}

func NewSerialSessionStartService(server *vs.VServer, context context.Context, sessionStartRequest *uiprpc.SessionStartRequest) SerialSessionStartService {
	return SerialSessionStartService{VServer: server, Context: context, SessionStartRequest: sessionStartRequest}
}

func (s SerialSessionStartService) RequestNSBForNewSession(anyb types.Session) ([]byte, error) {
	var accs = anyb.GetAccounts()

	var owners = make([][]byte, 0, len(accs)+1)
	// todo
	// owners = append(owners, s.Signer.GetPublicKey())
	for _, owner := range accs {
		owners = append(owners, owner.GetAddress())
		s.Logger.Info("waiting", "address", hex.EncodeToString(owner.GetAddress()))
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
	// s.Logger.Info("accs, txs", owners, txs)
	return s.NsbClient.CreateISC(s.Signer, make([]uint32, len(owners)), owners, txs, s.Signer.Sign(bytes.Join(anyb.GetTransactions(), []byte{})).Bytes())
}

func (s SerialSessionStartService) SessionStart() ([]byte, []uiptypes.Account, error) {
	var ses = new(session.SerialSession)
	ses.Signer = s.Signer
	success, help_info, err := ses.InitFromOpIntents(s.GetOpintents())
	if err != nil {
		// TODO: log
		s.Logger.Error("error", "error", err)
		return nil, nil, err
	}
	if !success {
		s.Logger.Error("error", "error", err)
		return nil, nil, errors.New(help_info)
	}
	ses.ISCAddress, err = s.RequestNSBForNewSession(ses)
	if err != nil {
		s.Logger.Error("error", "error", err)
		return nil, nil, err
	}
	err = s.DB.InsertSessionInfo(ses)
	if err != nil {
		s.Logger.Error("error", "error", err)
		return nil, nil, err
	}
	for i := uint32(0); i < ses.TransactionCount; i++ {
		_, err = s.NsbClient.FreezeInfo(s.Signer, ses.ISCAddress, uint64(i))
		if err != nil {
			s.Logger.Error("error", "error", err)
			return nil, nil, err
		}
	}
	// s.UpdateTxs
	// s.UpdateAccs
	return ses.ISCAddress, ses.GetAccounts(), nil
}

func (s SerialSessionStartService) Serve() (*uiprpc.SessionStartReply, error) {
	if b, accs, err := s.SessionStart(); err != nil {
		s.Logger.Error("error", "error", err)
		return nil, err
	} else {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		r, err := s.CVes.InternalRequestComing(ctx, &uiprpc.InternalRequestComingRequest{
			SessionId: b,
			Host:      s.Host,
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

		if err != nil {
			s.Logger.Error("error", "error", err)
			return nil, err
		}

		return &uiprpc.SessionStartReply{
			Ok:        r.GetOk(),
			SessionId: b,
		}, nil
	}
}

type MultiThreadSerialSessionStartService struct {
	*vs.VServer
	context.Context
	*uiprpc.SessionStartRequest
}

func NewMultiThreadSerialSessionStartService(server *vs.VServer, context context.Context, sessionStartRequest *uiprpc.SessionStartRequest) MultiThreadSerialSessionStartService {
	return MultiThreadSerialSessionStartService{VServer: server, Context: context, SessionStartRequest: sessionStartRequest}
}

func (s MultiThreadSerialSessionStartService) RequestNSBForNewSession(anyb types.Session) ([]byte, error) {
	var accs = anyb.GetAccounts()

	var owners = make([][]byte, 0, len(accs)+1)
	// todo
	// owners = append(owners, s.Signer.GetPublicKey())
	for _, owner := range accs {
		owners = append(owners, owner.GetAddress())
		s.Logger.Info("waiting", hex.EncodeToString(owner.GetAddress()))
	}
	var txs = anyb.GetTransactions()
	var btxs = make([][]byte, 0, len(txs))
	for _, tx := range txs {
		b, err := json.Marshal(tx)
		if err != nil {
			s.Logger.Error("error", "error", err)
			return nil, err
		}
		btxs = append(btxs, b)
	}

	return s.NsbClient.CreateISC(s.Signer, make([]uint32, len(owners)), owners, txs, s.Signer.Sign(bytes.Join(anyb.GetTransactions(), []byte{})).Bytes())
}

func (s MultiThreadSerialSessionStartService) SessionStart() ([]byte, []uiptypes.Account, error) {

	var ses = new(session.MultiThreadSerialSession)
	ses.Signer = s.Signer
	success, help_info, err := ses.InitFromOpIntents(s.GetOpintents())
	if err != nil {
		s.Logger.Error("error", "error", err)
		return nil, nil, err
	}
	if !success {
		return nil, nil, errors.New(help_info)
	}
	ses.ISCAddress, err = s.RequestNSBForNewSession(ses)
	if ses.ISCAddress == nil {
		err = fmt.Errorf("request isc failed: %v", err)
		s.Logger.Error("error", "error", err)
		return nil, nil, err
	}
	if err != nil {
		err = fmt.Errorf("request isc failed on request: %v", err)
		s.Logger.Error("error", "error", err)
		return nil, nil, err
	}
	err = ses.AfterInitGUID()
	logger.Println("after init guid...", ses.ISCAddress, hex.EncodeToString(ses.ISCAddress))
	if err != nil {
		s.Logger.Error("error", "error", err)
		return nil, nil, err
	}

	err = s.DB.InsertSessionInfo(ses)
	if err != nil {
		s.Logger.Error("error", "error", err)
		return nil, nil, err
	}
	for i := uint32(0); i < ses.TransactionCount; i++ {
		//s.Logger.Info()
		_, err := s.NsbClient.FreezeInfo(s.Signer, ses.ISCAddress, uint64(i))
		if err != nil {
			s.Logger.Error("error", "error", err)
			return nil, nil, err
		}
	}

	// s.UpdateTxs
	// s.UpdateAccs
	return ses.ISCAddress, ses.GetAccounts(), nil
}

func (s MultiThreadSerialSessionStartService) Serve() (*uiprpc.SessionStartReply, error) {
	fmt.Println(`###########################################################################################################################################
#                                                              SessionStart                                                               #
###########################################################################################################################################`)

	if b, accs, err := s.SessionStart(); err != nil {
		s.Logger.Error("error", "error", err)
		return nil, err
	} else {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		fmt.Println(`###########################################################################################################################################
#                                                              SendingRequest                                                             #
###########################################################################################################################################`)
		r, err := s.CVes.InternalRequestComing(ctx, &uiprpc.InternalRequestComingRequest{
			SessionId: b,
			Host:      s.Host,
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

		if err != nil {
			s.Logger.Error("error", "error", err)
			return nil, err
		}

		return &uiprpc.SessionStartReply{
			Ok:        r.GetOk(),
			SessionId: b,
		}, nil
	}
}
