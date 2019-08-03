package service

import (
	"encoding/json"
	"fmt"

	"golang.org/x/net/context"

	bni "github.com/Myriad-Dreamin/go-uip/bni/eth"
	tx "github.com/Myriad-Dreamin/go-uip/op-intent"
	uiprpc "github.com/Myriad-Dreamin/go-ves/grpc/uiprpc"
	uipbase "github.com/Myriad-Dreamin/go-ves/grpc/uiprpc-base"
	types "github.com/Myriad-Dreamin/go-ves/types"
)

type SessionRequireRawTransactService struct {
	types.VESDB
	context.Context
	*uiprpc.SessionRequireRawTransactRequest
}

func (s SessionRequireRawTransactService) Serve() (*uiprpc.SessionRequireRawTransactReply, error) {
	// todo errors.New("TODO")
	s.ActivateSession(s.GetSessionId())
	defer s.InactivateSession(s.GetSessionId())
	ses, err := s.FindSessionInfo(s.SessionId)
	if err != nil {
		return nil, err
	}
	var underTransacting uint32
	underTransacting, err = ses.GetTransactingTransaction()
	if err != nil {
		return nil, err
	}
	var transactionIntent tx.TransactionIntent
	err = s.FindTransaction(ses.GetGUID(), uint64(underTransacting), func(arg1 []byte) error {
		err := json.Unmarshal(arg1, &transactionIntent)
		return err
	})
	if err != nil {
		return nil, err
	}
	var b []byte
	b, err = (&bni.BN{}).Translate(&transactionIntent, s.GetGetter(ses.GetGUID()))
	if err != nil {
		return nil, err
	}
	fmt.Println("tid", underTransacting, "src", transactionIntent.Src, "dst", transactionIntent.Dst)
	return &uiprpc.SessionRequireRawTransactReply{
		RawTransaction: b,
		Tid:            uint64(underTransacting),
		Src: &uipbase.Account{
			Address: transactionIntent.Src,
			ChainId: transactionIntent.ChainID,
		},
		Dst: &uipbase.Account{
			Address: transactionIntent.Dst,
			ChainId: transactionIntent.ChainID,
		},
	}, nil
}
