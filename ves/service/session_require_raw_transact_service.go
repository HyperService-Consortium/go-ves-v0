package service

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"

	gjson "github.com/tidwall/gjson"

	"golang.org/x/net/context"

	ethbni "github.com/Myriad-Dreamin/go-ves/lib/bni/eth"
	transtype "github.com/Myriad-Dreamin/go-uip/const/trans_type"
	value_type "github.com/Myriad-Dreamin/go-uip/const/value_type"
	tx "github.com/Myriad-Dreamin/go-uip/op-intent"
	uiptypes "github.com/Myriad-Dreamin/go-uip/types"
	uiprpc "github.com/Myriad-Dreamin/go-ves/grpc/uiprpc"
	uipbase "github.com/Myriad-Dreamin/go-ves/grpc/uiprpc-base"
	types "github.com/Myriad-Dreamin/go-ves/types"
)

type SessionRequireRawTransactService struct {
	Resp *uipbase.Account
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

	if transactionIntent.TransType == transtype.ContractInvoke {
		var meta uiptypes.ContractInvokeMeta

		err := json.Unmarshal(transactionIntent.Meta, &meta)
		if err != nil {
			return nil, err
		}

		var intDesc uint16
		for _, param := range meta.Params {
			if intDesc = value_type.FromString(param.Type); intDesc == value_type.Unknown {
				return nil, errors.New("unknown type: " + param.Type)
			}

			result := gjson.ParseBytes(param.Value)
			if !result.Get("constant").Exists() {
				if result.Get("contract").Exists() &&
					result.Get("pos").Exists() &&
					result.Get("field").Exists() {
					ca, err := hex.DecodeString(result.Get("contract").String())
					if err != nil {
						return nil, err
					}
					pos, err := hex.DecodeString(result.Get("contract").String())
					if err != nil {
						return nil, err
					}
					desc := []byte(result.Get("field").String())

					v, err := new(ethbni.BN).GetStorageAt(transactionIntent.ChainID, intDesc, ca, pos, desc)
					if err != nil {
						return nil, err
					}
					vv, err := json.Marshal(v)
					s.SetKV(ses.GetGUID(), desc, vv)
					if err != nil {
						return nil, err
					}
				} else {
					return nil, errors.New("no enough info of source description")
				}
			}
		}

	}

	var b []byte
	b, err = (&ethbni.BN{}).Translate(&transactionIntent, s.GetGetter(ses.GetGUID()))
	if err != nil {
		return nil, err
	}

	if transactionIntent.TransType == transtype.Payment {

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
	} else {

		fmt.Println("tid", underTransacting, "src", transactionIntent.Src, "dst", s.Resp.GetAddress())
		return &uiprpc.SessionRequireRawTransactReply{
			RawTransaction: b,
			Tid:            uint64(underTransacting),
			Src: &uipbase.Account{
				Address: transactionIntent.Src,
				ChainId: transactionIntent.ChainID,
			},
			Dst: s.Resp,
		}, nil
	}

}
