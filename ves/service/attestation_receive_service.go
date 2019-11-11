package service

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/HyperService-Consortium/go-ves/config"
	"github.com/HyperService-Consortium/go-ves/ves/vs"
	"time"

	"golang.org/x/net/context"

	tx "github.com/HyperService-Consortium/go-uip/op-intent"
	"github.com/HyperService-Consortium/go-uip/uiptypes"
	"github.com/HyperService-Consortium/go-ves/grpc/uiprpc"
	uipbase "github.com/HyperService-Consortium/go-ves/grpc/uiprpc-base"
	ethbni "github.com/HyperService-Consortium/go-ves/lib/bni/eth"
	// bni "github.com/HyperService-Consortium/go-ves/types/bn-interface"

	"github.com/HyperService-Consortium/go-uip/signaturer"
	nsbi "github.com/HyperService-Consortium/go-ves/types/nsb-interface"
)

type AttestationReceiveService struct {
	*vs.VServer
	context.Context
	*uiprpc.AttestationReceiveRequest
}

func NewAttestationReceiveService(server *vs.VServer, context context.Context, attestationReceiveRequest *uiprpc.AttestationReceiveRequest) AttestationReceiveService {
	return AttestationReceiveService{VServer: server, Context: context, AttestationReceiveRequest: attestationReceiveRequest}
}

type AtteAdapdator struct {
	*uipbase.Attestation
}

func (atte *AtteAdapdator) GetSignatures() []uiptypes.Signature {
	var ss = atte.Attestation.GetSignatures()
	ret := make([]uiptypes.Signature, len(ss))
	for _, s := range ss {
		ret = append(ret, signaturer.FromRaw(s.Content, s.SignatureType))
	}
	return ret
}

func (s AttestationReceiveService) Serve() (*uiprpc.AttestationReceiveReply, error) {
	// todo
	s.DB.ActivateSession(s.GetSessionId())
	ses, err := s.DB.FindSessionInfo(s.GetSessionId())
	if err != nil {
		s.DB.InactivateSession(s.GetSessionId())
		s.Logger.Error("find session info", "sid", hex.EncodeToString(s.GetSessionId()), "err", err)
		return nil, err
	}

	defer func() {
		err = s.DB.UpdateSessionInfo(ses)
		if err != nil {
			s.Logger.Error("update failed", "sid", hex.EncodeToString(ses.GetGUID()), "err", err)
		}
		s.DB.InactivateSession(s.GetSessionId())
	}()

	ses.SetSigner(s.Signer)

	var success bool
	var helpInfo string

	currentTXID, _ := ses.GetTransactingTransaction()
	success, helpInfo, err = ses.ProcessAttestation(
		nsbi.NSBInterfaceFromClient(s.NsbClient, s.Signer),
		ethbni.NewBN(config.ChainDNS),
		&AtteAdapdator{s.GetAtte()},
	)
	if err != nil {
		s.Logger.Error("process transaction internal error", "sid", hex.EncodeToString(ses.GetGUID()),
			"tid", s.GetAtte().Tid, "aid", s.GetAtte().Aid, "err", err)
		return nil, fmt.Errorf("internal error: %v", err)
	} else if !success {
		s.Logger.Error("process transaction error", "sid", hex.EncodeToString(ses.GetGUID()),
			"tid", s.GetAtte().Tid, "aid", s.GetAtte().Aid, "err", err)
		return nil, errors.New(helpInfo)
	}

	fixedTXID, err := ses.GetTransactingTransaction()

	if err != nil {
		s.Logger.Error("get transaction error", "sid", hex.EncodeToString(ses.GetGUID()), "getting id", fixedTXID, "err", err)
		return nil, fmt.Errorf("internal error: %v", err)
	}

	if fixedTXID == uint32(len(ses.GetTransactions())) {
		// close
		return &uiprpc.AttestationReceiveReply{
			Ok: true,
		}, nil
	}

	if fixedTXID != currentTXID {

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		txb := ses.GetTransaction(fixedTXID)
		var kvs tx.TransactionIntent
		err := json.Unmarshal(txb, &kvs)
		if err != nil {
			s.Logger.Error("unmarshal error", "sid", hex.EncodeToString(ses.GetGUID()), "err", err)
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
			s.Logger.Error("send message error", "sid", hex.EncodeToString(ses.GetGUID()), "err", err)
			return nil, err
		}
	}
	return &uiprpc.AttestationReceiveReply{
		Ok: true,
	}, nil
}
