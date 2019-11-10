package service

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"

	gjson "github.com/tidwall/gjson"

	"golang.org/x/net/context"

	transtype "github.com/HyperService-Consortium/go-uip/const/trans_type"
	value_type "github.com/HyperService-Consortium/go-uip/const/value_type"
	tx "github.com/HyperService-Consortium/go-uip/op-intent"
	uiptypes "github.com/HyperService-Consortium/go-uip/uiptypes"
	uiprpc "github.com/HyperService-Consortium/go-ves/grpc/uiprpc"
	uipbase "github.com/HyperService-Consortium/go-ves/grpc/uiprpc-base"
	ethbni "github.com/HyperService-Consortium/go-ves/lib/bni/eth"
	tenbni "github.com/HyperService-Consortium/go-ves/lib/bni/ten"
	types "github.com/HyperService-Consortium/go-ves/types"
)

type SessionRequireRawTransactService struct {
	Resp *uipbase.Account
	types.VESDB
	context.Context
	*uiprpc.SessionRequireRawTransactRequest
}

var bnis = map[uint64]uiptypes.BlockChainInterface{
	1: new(ethbni.BN),
	2: new(ethbni.BN),
	3: new(tenbni.BN),
	4: new(tenbni.BN),
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

	bn := bnis[transactionIntent.ChainID]

	if transactionIntent.TransType == transtype.ContractInvoke {
		var meta uiptypes.ContractInvokeMeta

		err := json.Unmarshal(transactionIntent.Meta, &meta)
		if err != nil {
			return nil, err
		}

		var intDesc uiptypes.TypeID
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
					pos, err := hex.DecodeString(result.Get("pos").String())
					if err != nil {
						return nil, err
					}
					desc := []byte(result.Get("field").String())

					v, err := bn.GetStorageAt(transactionIntent.ChainID, intDesc, ca, pos, desc)
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

	var b uiptypes.RawTransaction
	b, err = bn.Translate(&transactionIntent, s.GetGetter(ses.GetGUID()))
	if err != nil {
		return nil, err
	}

	if transactionIntent.TransType == transtype.Payment {

		fmt.Println("tid", underTransacting, "src", transactionIntent.Src, "dst", transactionIntent.Dst)
		return &uiprpc.SessionRequireRawTransactReply{
			RawTransaction: b.Bytes(),
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
			RawTransaction: b.Bytes(),
			Tid:            uint64(underTransacting),
			Src: &uipbase.Account{
				Address: transactionIntent.Src,
				ChainId: transactionIntent.ChainID,
			},
			Dst: s.Resp,
		}, nil
	}

}
