package service

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/HyperService-Consortium/go-ves/ves/vs"
	"time"

	"golang.org/x/net/context"

	tx "github.com/HyperService-Consortium/go-uip/op-intent"
	"github.com/HyperService-Consortium/go-uip/signaturer"
	"github.com/HyperService-Consortium/go-ves/grpc/uiprpc"
	uipbase "github.com/HyperService-Consortium/go-ves/grpc/uiprpc-base"
)

type SessionAckForInitService struct {
	*vs.VServer
	context.Context
	*uiprpc.SessionAckForInitRequest
}

func NewSessionAckForInitService(server *vs.VServer, context context.Context, sessionAckForInitRequest *uiprpc.SessionAckForInitRequest) SessionAckForInitService {
	return SessionAckForInitService{VServer: server, Context: context, SessionAckForInitRequest: sessionAckForInitRequest}
}

func (s SessionAckForInitService) Serve() (*uiprpc.SessionAckForInitReply, error) {
	s.Logger.Info("session acknowledging... ", "address", hex.EncodeToString(s.GetUser().GetAddress()))

	s.DB.ActivateSession(s.GetSessionId())
	defer s.DB.InactivateSession(s.GetSessionId())
	ses, err := s.DB.FindSessionInfo(s.SessionId)
	// todo: get Session Acked from isc
	// nsbClient.
	if err != nil {
		s.Logger.Error("find session info error", "error", err)
		return nil, err
	}

	var success bool
	var help_info string
	ss := s.GetUserSignature()
	success, help_info, err = ses.AckForInit(s.GetUser(), signaturer.FromRaw(ss.Content, ss.SignatureType))
	if err != nil {
		s.Logger.Error("ack error", "error", err)
		return nil, fmt.Errorf("internal error: %v", err)
	} else if !success {
		return nil, errors.New(help_info)
	}

	if ses.GetAckCount() == uint32(len(ses.GetAccounts())) {

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		txb := ses.GetTransaction(0)
		var kvs tx.TransactionIntent
		err := json.Unmarshal(txb, &kvs)
		if err != nil {
			return nil, err
		}
		var accs []*uipbase.Account
		accs = append(accs, &uipbase.Account{
			Address: kvs.Src,
			ChainId: kvs.ChainID,
		})
		s.Logger.Info("sending attestation request", "chain id", kvs.ChainID, "address", hex.EncodeToString(kvs.Src))

		_, err = s.CVes.InternalAttestationSending(ctx, &uiprpc.InternalRequestComingRequest{
			SessionId: ses.GetGUID(),
			Host:      s.Host,
			Accounts:  accs,
		})
		if err != nil {
			return nil, err
		}
	}
	if err = s.DB.UpdateSessionInfo(ses); err != nil {
		s.Logger.Error("sending attestation request", "err", err)
		return nil, err
	}
	return &uiprpc.SessionAckForInitReply{
		Ok: true,
	}, nil
}
