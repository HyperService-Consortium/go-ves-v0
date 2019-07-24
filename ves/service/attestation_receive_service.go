package service

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"golang.org/x/net/context"

	tx "github.com/Myriad-Dreamin/go-uip/op-intent"
	uiptypes "github.com/Myriad-Dreamin/go-uip/types"
	uiprpc "github.com/Myriad-Dreamin/go-ves/grpc/uiprpc"
	uipbase "github.com/Myriad-Dreamin/go-ves/grpc/uiprpc-base"
	log "github.com/Myriad-Dreamin/go-ves/log"
	types "github.com/Myriad-Dreamin/go-ves/types"
	bni "github.com/Myriad-Dreamin/go-ves/types/bn-interface"
	nsbi "github.com/Myriad-Dreamin/go-ves/types/nsb-interface"
)

type AttestationReceiveService struct {
	CVes uiprpc.CenteredVESClient
	Host string
	uiptypes.Signer
	types.VESDB
	context.Context
	*uiprpc.AttestationReceiveRequest
}

type AtteAdapdator struct {
	*uipbase.Attestation
}

func (atte *AtteAdapdator) GetSignatures() []uiptypes.Signature {
	var ss = atte.Attestation.GetSignatures()
	ret := make([]uiptypes.Signature, len(ss))
	for _, s := range ss {
		ret = append(ret, uiptypes.Signature(s))
	}
	return ret
}

func (s *AttestationReceiveService) Serve() (*uiprpc.AttestationReceiveReply, error) {
	// todo
	ses, err := s.FindSessionInfo(s.GetSessionId())
	defer s.UpdateSessionInfo(ses)

	if err == nil {
		ses.SetSigner(s.Signer)

		var success bool
		var helpInfo string

		current_tx_id, _ := ses.GetTransactingTransaction()
		success, helpInfo, err = ses.ProcessAttestation(
			nsbi.NSBInterfaceImpl(s.Host, s.Signer),
			&bni.BN{},
			&AtteAdapdator{s.GetAtte()},
		)

		fixed_tx_id, _ := ses.GetTransactingTransaction()

		if err != nil {
			// todo, log
			return nil, fmt.Errorf("internal error: %v", err)
		} else if !success {
			return nil, errors.New(helpInfo)
		} else {
			if fixed_tx_id == uint32(len(ses.GetTransactions())) {
				// close
				return &uiprpc.AttestationReceiveReply{
					Ok: true,
				}, nil
			}
			if fixed_tx_id != current_tx_id {

				ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
				defer cancel()
				txb := ses.GetTransaction(fixed_tx_id)
				var kvs tx.TransactionIntent
				err := json.Unmarshal(txb, &kvs)
				if err != nil {
					return nil, err
				}
				var accs []*uipbase.Account
				accs = append(accs, &uipbase.Account{
					Address: kvs.Src,
					ChainId: kvs.ChainId,
				})
				log.Printf("sending attestation request to %v %v\n", hex.EncodeToString(kvs.Src), kvs.ChainId)
				// accs = append(accs, &uipbase.Account{
				// 	Address: kvs.Dst,
				// 	ChainId: kvs.ChainId,
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
			return &uiprpc.AttestationReceiveReply{
				Ok: true,
			}, nil
		}

	} else {
		return nil, err
	}
}
