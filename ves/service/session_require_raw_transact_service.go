package service

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/HyperService-Consortium/go-ves/config"
	payment_option "github.com/HyperService-Consortium/go-ves/lib/bni/payment-option"
	logger "github.com/HyperService-Consortium/go-ves/lib/log"
	"github.com/tidwall/gjson"

	"golang.org/x/net/context"

	transtype "github.com/HyperService-Consortium/go-uip/const/trans_type"
	"github.com/HyperService-Consortium/go-uip/const/value_type"
	tx "github.com/HyperService-Consortium/go-uip/op-intent"
	"github.com/HyperService-Consortium/go-uip/uiptypes"
	"github.com/HyperService-Consortium/go-ves/grpc/uiprpc"
	uipbase "github.com/HyperService-Consortium/go-ves/grpc/uiprpc-base"
	ethbni "github.com/HyperService-Consortium/go-ves/lib/bni/eth"
	tenbni "github.com/HyperService-Consortium/go-ves/lib/bni/ten"
	"github.com/HyperService-Consortium/go-ves/types"
)

type SessionRequireRawTransactService struct {
	Resp *uipbase.Account
	types.VESDB
	context.Context
	*uiprpc.SessionRequireRawTransactRequest
}

func (s SessionRequireRawTransactService) GetTransactionProof(chainID uiptypes.ChainID, blockID uiptypes.BlockID, color []byte) (uiptypes.MerkleProof, error) {
	// todo
	panic("implement me")
}

var bnis = map[uint64]uiptypes.BlockChainInterface{
	1: ethbni.NewBN(config.ChainDNS),
	2: ethbni.NewBN(config.ChainDNS),
	3: tenbni.NewBN(config.ChainDNS),
	4: tenbni.NewBN(config.ChainDNS),
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
	//fmt.Println(ses.(*session.MultiThreadSerialSession).Transactions)
	//fmt.Println("underTransacting", underTransacting, ses.GetGUID(), ses)

	var transactionIntent tx.TransactionIntent
	err = s.FindTransaction(ses.GetGUID(), uint64(underTransacting), func(arg1 []byte) error {
		err := json.Unmarshal(arg1, &transactionIntent)
		return err
	})
	if err != nil {
		return nil, err
	}
	//fmt.Println(".......")

	bn := bnis[transactionIntent.ChainID]

	if transactionIntent.TransType == transtype.ContractInvoke {
		fmt.Println("aaaaaa")
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
	} else {
		n, ok, err := payment_option.NeedInconsistentValueOption(gjson.ParseBytes(transactionIntent.Meta))
		if err != nil {
			logger.Warn("omit need inc-val option")
			return nil, err
		}
		if ok {
			v, err := bnis[2].GetStorageAt(n.ChainID, n.TypeID, n.ContractAddress, n.Pos, n.Description)
			if err != nil {
				logger.Error("get failed")
				return nil, err
			}
			logger.Info("getting state from blockchain", "address", n.ContractAddress, "value-type:", v.GetValue(), v.GetType(), "at pos", n.Pos)
			err = s.VESDB.SetStorageOf(n.ChainID, n.TypeID, n.ContractAddress, n.Pos, n.Description, v)
			if err != nil {
				logger.Error("set failed")
				return nil, err
			}
		}
	}

	var b uiptypes.RawTransaction
	b, err = bn.Translate(&transactionIntent, s)
	if err != nil {
		logger.Error("translate error", "err is ", err)
		return nil, err
	}

	if transactionIntent.TransType == transtype.Payment {

		fmt.Println("tid", underTransacting, "src", transactionIntent.Src, "dst", transactionIntent.Dst)
		x, err := b.Serialize()
		if err != nil {
			return nil, err
		}

		return &uiprpc.SessionRequireRawTransactReply{
			RawTransaction: x,
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
		x, err := b.Serialize()
		if err != nil {
			return nil, err
		}


		fmt.Println("tid", underTransacting, "src", transactionIntent.Src, "dst", s.Resp.GetAddress())
		return &uiprpc.SessionRequireRawTransactReply{
			RawTransaction: x,
			Tid:            uint64(underTransacting),
			Src: &uipbase.Account{
				Address: transactionIntent.Src,
				ChainId: transactionIntent.ChainID,
			},
			Dst: s.Resp,
		}, nil
	}

}
