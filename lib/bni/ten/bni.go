package bni

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	transactiontype "github.com/HyperService-Consortium/NSB/application/transaction-type"
	"github.com/HyperService-Consortium/NSB/math"
	"github.com/HyperService-Consortium/go-ves/types"
	"github.com/gogo/protobuf/proto"
	"net/url"
	"strings"

	"github.com/HyperService-Consortium/NSB/grpc/nsbrpc"
	TransType "github.com/HyperService-Consortium/go-uip/const/trans_type"
	"github.com/HyperService-Consortium/go-uip/uiptypes"
	nsbcli "github.com/HyperService-Consortium/go-ves/lib/net/nsb-client"
)

func decoratePrefix(hexs string) string {
	if strings.HasPrefix(hexs, "0x") {
		return hexs
	} else {
		return "0x" + hexs
	}
}

type BN struct {
	dns types.ChainDNSInterface
	signer uiptypes.Signer
}
var (
	errorDecodeSrcAddress = errors.New("the src address should be in length of 32")
	errorDecodeDstAddress = errors.New("the dst address should be in length of 32 or 0")
)

func (bn *BN) Deserialize(raw []byte) (uiptypes.RawTransaction, error) {

	var txHeader nsbrpc.TransactionHeader
	err := proto.Unmarshal(raw, &txHeader)
	if err != nil {
		return nil, err
	}
	if len(txHeader.Src) != 32 {
		return nil, errorDecodeSrcAddress
	}
	if len(txHeader.Dst) != 32 && len(txHeader.Dst) != 0 {
		return nil, errorDecodeDstAddress
	}

	return &rawTransaction{
		Type:   transactiontype.Type(raw[0]),
		Header: &txHeader,
	}, nil
}

var (
	ErrNotSigned = errors.New("not signed")
)

func (bn *BN) RouteRaw(destination uiptypes.ChainID, rawTransaction uiptypes.RawTransaction) (
	transactionReceipt uiptypes.TransactionReceipt, err error) {
	if !rawTransaction.Signed() {
		return nil, ErrNotSigned
	}
	ci, err := bn.dns.GetChainInfo(destination)
	if err != nil {
		return nil, err
	}
	// todo receipt
	b, err := rawTransaction.Serialize()
	if err != nil {
		return nil, err
	}
	b, err = nsbcli.NewNSBClient((&url.URL{Scheme: "http", Host: ci.GetChainHost(), Path: "/"}).String()).BroadcastTxCommitReturnBytes(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (bn *BN) WaitForTransact(_ uiptypes.ChainID, transactionReceipt uiptypes.TransactionReceipt,
	options ...interface{}) (blockID []byte, color []byte, err error) {
	var res nsbcli.ResultInfo
	err = json.Unmarshal(transactionReceipt, &res)
	if err != nil {
		return nil, nil, err
	}

	return []byte(res.Height), []byte(res.Hash), err
}

func (bn *BN) RouteWithSigner(signer uiptypes.Signer) (uiptypes.Router, error) {
	nbn := *bn
	nbn.signer = signer
	return &nbn, nil
}

type rawTransaction struct {
	Type transactiontype.Type
	Header *nsbrpc.TransactionHeader
}

func newRawTransaction(_type transactiontype.Type, header *nsbrpc.TransactionHeader) *rawTransaction {
	return &rawTransaction{Type: _type, Header: header}
}


func (r *rawTransaction) Serialize() ([]byte, error) {
	return nsbcli.GlobalClient.Serialize(r.Type, r.Header)
}

func (r *rawTransaction) Signed() bool {
	return len(r.Header.Signature) != 0
}

func (r *rawTransaction) Sign(user uiptypes.Signer) (uiptypes.RawTransaction, error) {
	if !bytes.Equal(r.Header.Src, user.GetPublicKey()) {
		return nil, fmt.Errorf("sign error user is %v, want is %v", hex.EncodeToString(user.GetPublicKey()), hex.EncodeToString(r.Header.Src))
	}
	r.Header = nsbcli.GlobalClient.Sign(user, r.Header)
	return r, nil
}

func (bn *BN) Translate(intent *uiptypes.TransactionIntent, storage uiptypes.Storage) (uiptypes.RawTransaction, error) {
	switch intent.TransType {
	case TransType.Payment:
		header, err := nsbcli.GlobalClient.CreateTransferPacket(intent.Src, intent.Dst, math.NewUint256FromHexString(intent.Amt))
		if err != nil {
			return nil, err
		}
		return newRawTransaction(transactiontype.Type(intent.TransType), header), nil
	case TransType.ContractInvoke:
		// var meta uiptypes.ContractInvokeMeta
		//
		// err := json.Unmarshal(intent.Meta, &meta)
		// if err != nil {
		// 	return nil, err
		// }
		// //_ = meta
		//
		// data, err := ContractInvocationDataABI(&meta, storage)
		// if err != nil {
		// 	return nil, err
		// }
		//
		// hexdata := hex.EncodeToString(data)
		// // meta.FuncName
		//
		// return json.Marshal(map[string]interface{}{
		// 	"jsonrpc": "2.0",
		// 	"method":  "eth_sendTransaction",
		// 	"params": []interface{}{
		// 		map[string]interface{}{
		// 			"from":  decoratePrefix(hex.EncodeToString(intent.Src)),
		// 			"to":    decoratePrefix(hex.EncodeToString(intent.Dst)),
		// 			"value": decoratePrefix(intent.Amt),
		// 			"data":  decoratePrefix(hexdata),
		// 		},
		// 	},
		// 	"id": 1,
		// })
		return nil, errors.New("todo")
	default:
		return nil, errors.New("cant translate")
	}
}

func (bn *BN) GetStorageAt(chainID uiptypes.ChainID, typeID uiptypes.TypeID, contractAddress uiptypes.ContractAddress, pos []byte, description []byte) (uiptypes.Variable, error) {
	return nil, errors.New("todo")
	//ci, err := bn.dns.GetChainInfo(chainID)
	//if err != nil {
	//	return nil, err
	//}
	//
	//switch typeID {
	//case valuetype.Bool:
	//	s, err := ethclient.NewEthClient((&url.URL{Scheme: "http", Host: ci.GetChainHost(), Path: "/"}).String()).GetStorageAt(contractAddress, pos, "latest")
	//	if err != nil {
	//		return nil, err
	//	}
	//
	//	b, err := hex.DecodeString(s[2:])
	//	if err != nil {
	//		return nil, err
	//	}
	//	bs, err := ethabi.NewDecoder().Decodes([]string{"bool"}, b)
	//	return bs[0], nil
	//case valuetype.Uint256:
	//	s, err := ethclient.NewEthClient(ci.GetChainHost()).GetStorageAt(contractAddress, pos, "latest")
	//	if err != nil {
	//		return nil, err
	//	}
	//
	//	b, err := hex.DecodeString(s[2:])
	//	if err != nil {
	//		return nil, err
	//	}
	//	bs, err := ethabi.NewDecoder().Decodes([]string{"uint256"}, b)
	//	return bs[0], nil
	//}

	//return nil, nil
}

func NewBN(dns types.ChainDNSInterface) *BN {
	return &BN{dns: dns}
}

func (bn *BN) MustWithSigner() bool {
	return true
}


//func (bn *BN) Route(intent *uiptypes.TransactionIntent, kvGetter uiptypes.KVGetter) ([]byte, error) {
//	// todo
//	onChainTransaction, err := bn.Translate(intent, kvGetter)
//	if err != nil {
//		return nil, err
//	}
//	return bn.RouteRaw(intent.ChainID, onChainTransaction)
//}


func (bn *BN) CheckAddress(addr []byte) bool {
	return len(addr) == 32
}

