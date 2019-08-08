package service

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"golang.org/x/net/context"

	tx "github.com/Myriad-Dreamin/go-uip/op-intent"
	signaturer "github.com/Myriad-Dreamin/go-uip/signaturer"
	uiptypes "github.com/Myriad-Dreamin/go-uip/types"
	uiprpc "github.com/Myriad-Dreamin/go-ves/grpc/uiprpc"
	uipbase "github.com/Myriad-Dreamin/go-ves/grpc/uiprpc-base"
	log "github.com/Myriad-Dreamin/go-ves/lib/log"
	types "github.com/Myriad-Dreamin/go-ves/types"
)

type SessionAckForInitService struct {
	CVes uiprpc.CenteredVESClient
	uiptypes.Signer
	types.VESDB
	context.Context
	*uiprpc.SessionAckForInitRequest
}

func (s SessionAckForInitService) Serve() (*uiprpc.SessionAckForInitReply, error) {
	s.ActivateSession(s.GetSessionId())
	defer s.InactivateSession(s.GetSessionId())
	ses, err := s.FindSessionInfo(s.SessionId)
	// todo: get Session Acked from isc
	// nsbClient.
	log.Println("session acking... ", hex.EncodeToString(s.GetUser().GetAddress()))
	if err == nil {
		var success bool
		var help_info string
		success, help_info, err = ses.AckForInit(s.GetUser(), signaturer.FromBaseSignature(s.GetUserSignature()))
		if err != nil {
			// todo, log
			return nil, fmt.Errorf("internal error: %v", err)
		} else if !success {
			return nil, errors.New(help_info)
		} else {
			// fmt.Println(ses.GetAckCount(), uint32(len(ses.GetAccounts())))
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
				log.Printf("sending attestation request to %v %v\n", hex.EncodeToString(kvs.Src), kvs.ChainID)
				// accs = append(accs, &uipbase.Account{
				// 	Address: kvs.Dst,
				// 	ChainID: kvs.ChainID,
				// })
				_, err = s.CVes.InternalAttestationSending(ctx, &uiprpc.InternalRequestComingRequest{
					SessionId: ses.GetGUID(),
					Host:      []byte{127, 0, 0, 1, ((23351) >> 8 & 0xff), 23351 & 0xff},
					Accounts:  accs,
				})
				// fmt.Println("reply?", r, err)
				if err != nil {
					return nil, err
				}

				/*
					atte := &uipbase.Attestation{
						Content: ses.GetTransaction(0),
						Signatures: []*uipbase.Signature{
							&uipbase.Signature{
								Content:       s.Sign(ses.GetTransaction(0)),
								SignatureType: 123456,
							},
						},
					}
				*/
			}
			err = s.VESDB.UpdateSessionInfo(ses)
			if err != nil {
				fmt.Println("uperr", err)
				return nil, err
			}
			return &uiprpc.SessionAckForInitReply{
				Ok: true,
			}, nil
		}
	} else {
		return nil, err
	}
}
