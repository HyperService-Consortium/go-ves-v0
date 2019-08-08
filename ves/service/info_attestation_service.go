package service

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"golang.org/x/net/context"

	ethbni "github.com/Myriad-Dreamin/go-uip/bni/eth"
	tx "github.com/Myriad-Dreamin/go-uip/op-intent"
	uiptypes "github.com/Myriad-Dreamin/go-uip/types"
	uiprpc "github.com/Myriad-Dreamin/go-ves/grpc/uiprpc"
	uipbase "github.com/Myriad-Dreamin/go-ves/grpc/uiprpc-base"
	ethbni "github.com/Myriad-Dreamin/go-ves/lib/bni/eth"
	log "github.com/Myriad-Dreamin/go-ves/lib/log"
	types "github.com/Myriad-Dreamin/go-ves/types"
	nsbi "github.com/Myriad-Dreamin/go-ves/types/nsb-interface"
)

type InformAttestationService struct {
	CVes uiprpc.CenteredVESClient
	Host string
	uiptypes.Signer
	types.VESDB
	context.Context
	*uiprpc.AttestationReceiveRequest
}

func (s *InformAttestationService) Serve() (*uiprpc.AttestationReceiveReply, error) {
	// todo
	s.ActivateSession(s.GetSessionId())
	ses, err := s.FindSessionInfo(s.GetSessionId())
	tid, _ := ses.GetTransactingTransaction()
	fmt.Printf("%T\n", ses)
	fmt.Println("this is about ", tid)
	if err == nil {
		defer func() {
			s.UpdateSessionInfo(ses)
			s.InactivateSession(s.GetSessionId())
		}()
		ses.SetSigner(s.Signer)

		var success bool
		var helpInfo string

		current_tx_id, _ := ses.GetTransactingTransaction()
		success, helpInfo, err = ses.NotifyAttestation(
			nsbi.NSBInterfaceImpl(s.Host, s.Signer),
			&ethbni.BN{},
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

				if len(helpInfo) != 0 {
					log.Infoln("InformAttestationService:", helpInfo)
				}

				// if ret, err := nsbClient.SettleContract(s.Signer, ses.GetGUID()); err != nil {
				// 	return nil, err
				// } else {
				// 	fmt.Printf(
				// 		"closing contract {\n\tinfo: %v,\n\tdata: %v,\n\tlog: %v, \n\ttags: %v\n}\n",
				// 		ret.Info, string(ret.Data), ret.Log, ret.Tags,
				// 	)
				// }

				ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
				defer cancel()

				raccs := ses.GetAccounts()

				var accs []*uipbase.Account

				for _, acc := range raccs {
					accs = append(accs, &uipbase.Account{
						Address: acc.GetAddress(),
						ChainId: acc.GetChainId(),
					})
				}

				_, err = s.CVes.InternalCloseSession(ctx, &uiprpc.InternalCloseSessionRequest{
					SessionId: ses.GetGUID(),
					NsbHost:   []byte{47, 251, 2, 73, uint8(26657 >> 8), uint8(26657 & 0xff)},
					GrpcHost:  []byte{127, 0, 0, 1, ((23351) >> 8 & 0xff), 23351 & 0xff},
					Accounts:  accs,
				})
				// fmt.Println("reply?", r, err)
				if err != nil {
					return nil, err
				}

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
					ChainId: kvs.ChainID,
				})
				log.Printf("sending attestation request to %v %v\n", hex.EncodeToString(kvs.Src), kvs.ChainID)
				// accs = append(accs, &uipbase.Account{
				// 	Address: kvs.Dst,
				// 	ChainId: kvs.ChainID,
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
		s.InactivateSession(s.GetSessionId())
		return nil, err
	}
}
