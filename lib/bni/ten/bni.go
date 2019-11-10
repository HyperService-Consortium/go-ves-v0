package bni

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/HyperService-Consortium/go-ves/types"
	"net/url"
	"strings"

	"github.com/HyperService-Consortium/NSB/grpc/nsbrpc"
	"github.com/gogo/protobuf/proto"

	"github.com/HyperService-Consortium/go-ethabi"
	TransType "github.com/HyperService-Consortium/go-uip/const/trans_type"
	valuetype "github.com/HyperService-Consortium/go-uip/const/value_type"
	"github.com/HyperService-Consortium/go-uip/uiptypes"
	ethclient "github.com/HyperService-Consortium/go-ves/lib/net/eth-client"
	nsbclient "github.com/HyperService-Consortium/go-ves/lib/net/nsb-client"
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

func (bn *BN) RouteRaw(destination uiptypes.ChainID, rawTransaction uiptypes.RawTransaction) (
	transactionReceipt uiptypes.TransactionReceipt, err error) {
	ci, err := bn.dns.GetChainInfo(destination)
	if err != nil {
		return nil, err
	}

	_ = ci

	return nil, errors.New("TODO")
	//var txHeader MiddleHeader
	//err = json.Unmarshal(payload, &txHeader)
	//if err != nil {
	//	return nil, err
	//}
	//
	//// bug: buf.Reset()
	//buf := bytes.NewBuffer(make([]byte, 65535))
	//
	//buf.Write(txHeader.Header.Src)
	//buf.Write(txHeader.Header.Dst)
	//buf.Write(txHeader.Header.Data)
	//buf.Write(txHeader.Header.Value)
	//buf.Write(txHeader.Header.Nonce)
	//txHeader.Header.Signature = bn.signer.Sign(buf.Bytes()).Bytes()
	//
	//if err != nil {
	//	return nil, err
	//}
	//
	//b, err := json.Marshal(txHeader.Header)
	//if err != nil {
	//	return nil, err
	//}
	//
	//buf.Reset()
	//buf.Write(txHeader.PreHeader)
	//buf.Write(b)

	//return nsbclient.NewNSBClient((&url.URL{Scheme: "http", Host: ci.GetChainHost(), Path: "/"}).String()).BroadcastTxCommitReturnBytes(buf.Bytes())
}

func (bn *BN) WaitForTransact(chainID uiptypes.ChainID, transactionReceipt uiptypes.TransactionReceipt,
	options ...interface{}) (blockID uiptypes.BlockID, color []byte, err error) {
	// todo

	var info RTxInfo

	err = json.Unmarshal(transactionReceipt, &info)

	return nil, info.transactionReceipt, err
	// cinfo, err := SearchChainId(chainID)
	// if err != nil {
	// 	return nil, err
	// }
	// worker := nsbclient.NewNSBClient(cinfo.GetHost())
	// ddl := time.Now().Add(timeout)
	// for time.Now().Before(ddl) {
	// 	tx, err := worker.GetProof(receipt, `"prove_on_tx_trie"`)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	fmt.Println(string(tx))
	// 	if gjson.GetBytes(tx, "blockNumber").Type != gjson.Null {
	// 		b, _ := hex.DecodeString(gjson.GetBytes(tx, "blockHash").String()[2:])
	// 		return b, nil
	// 	}
	// 	time.Sleep(time.Millisecond * 500)
	//
	// }
	// return nil, errors.New("timeout")
}

func (bn *BN) RouteWithSigner(signer uiptypes.Signer) (uiptypes.Router, error) {
	nbn := *bn
	nbn.signer = signer
	return &nbn, nil
}

func (bn *BN) Translate(intent *uiptypes.TransactionIntent, kvGetter uiptypes.KVGetter) (uiptypes.RawTransaction, error) {
	switch intent.TransType {
	case TransType.Payment:
		return nil, errors.New("todo")
		//var txHeader MiddleHeader
		//
		//// Nonce
		//nonce := make([]byte, 32)
		//_, err := rand.Read(nonce)
		//if err != nil {
		//	return nil, err
		//}
		//txHeader.Header.Nonce = nonce
		//
		//// Data
		//var args appl.ArgsTransfer
		//value, err := hex.DecodeString(intent.Amt)
		//
		//if err != nil {
		//	return nil, err
		//}
		//args.Value = nsbmath.NewUint256FromBytes(value)
		//
		//b, err := json.Marshal(args)
		//if err != nil {
		//	return nil, err
		//}
		//
		//var fap appl.FAPair
		//fap.FuncName = "transfer"
		//fap.Args = b
		//txHeader.Header.Data, err = json.Marshal(fap)
		//if err != nil {
		//	return nil, err
		//}
		//
		//// Rest
		//txHeader.Header.Src = intent.Dst
		//txHeader.Header.Dst = intent.Src
		//txHeader.Header.Value = args.Value.Bytes()
		//txHeader.PreHeader = []byte("systemCall\x19system.token\x18")
		//
		//return json.Marshal(txHeader)
	case TransType.ContractInvoke:
		// var meta uiptypes.ContractInvokeMeta
		//
		// err := json.Unmarshal(intent.Meta, &meta)
		// if err != nil {
		// 	return nil, err
		// }
		// //_ = meta
		//
		// data, err := ContractInvocationDataABI(&meta, kvGetter)
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

func (bn *BN) GetStorageAt(chainID uiptypes.ChainID, typeID uiptypes.TypeID, contractAddress uiptypes.Contract, pos uiptypes.Pos, description uiptypes.Desc) (interface{}, error) {

	ci, err := bn.dns.GetChainInfo(chainID)
	if err != nil {
		return nil, err
	}

	switch typeID {
	case valuetype.Bool:
		s, err := ethclient.NewEthClient((&url.URL{Scheme: "http", Host: ci.GetChainHost(), Path: "/"}).String()).GetStorageAt(contractAddress, pos, "latest")
		if err != nil {
			return nil, err
		}

		b, err := hex.DecodeString(s[2:])
		if err != nil {
			return nil, err
		}
		bs, err := ethabi.NewDecoder().Decodes([]string{"bool"}, b)
		return bs[0], nil
	case valuetype.Uint256:
		s, err := ethclient.NewEthClient(ci.GetChainHost()).GetStorageAt(contractAddress, pos, "latest")
		if err != nil {
			return nil, err
		}

		b, err := hex.DecodeString(s[2:])
		if err != nil {
			return nil, err
		}
		bs, err := ethabi.NewDecoder().Decodes([]string{"uint256"}, b)
		return bs[0], nil
	}

	return nil, nil
}

func NewBN(dns types.ChainDNSInterface) *BN {
	return &BN{dns: dns}
}

type MiddleHeader struct {
	Header    nsbrpc.TransactionHeader `json:"h"`
	PreHeader []byte                   `json:"p"`
}

func (bn *BN) MustWithSigner() bool {
	return true
}

type RTxInfo struct {
	ret                []byte
	transactionReceipt []byte
}

func (bn *BN) RouteRawTransaction(destination uint64, payload []byte) ([]byte, error) {
	ci, err := bn.dns.GetChainInfo(destination)
	if err != nil {
		return nil, err
	}

	var txHeader MiddleHeader
	err = json.Unmarshal(payload, &txHeader)
	if err != nil {
		return nil, err
	}

	// bug: buf.Reset()
	buf := bytes.NewBuffer(make([]byte, 65535))

	buf.Write(txHeader.Header.Src)
	buf.Write(txHeader.Header.Dst)
	buf.Write(txHeader.Header.Data)
	buf.Write(txHeader.Header.Value)
	buf.Write(txHeader.Header.Nonce)
	txHeader.Header.Signature = bn.signer.Sign(buf.Bytes()).Bytes()

	if err != nil {
		return nil, err
	}

	b, err := proto.Marshal(&txHeader.Header)
	if err != nil {
		return nil, err
	}

	buf.Reset()
	buf.Write(txHeader.PreHeader)
	buf.Write(b)

	var ret RTxInfo

	ret.ret, err = nsbclient.NewNSBClient((&url.URL{Scheme: "http", Host: ci.GetChainHost(), Path: "/"}).String()).BroadcastTxCommitReturnBytes(buf.Bytes())
	if err != nil {
		return nil, err
	}

	ret.transactionReceipt = b
	b, err = json.Marshal(ret)
	return b, nil
}

func (bn *BN) Route(intent *uiptypes.TransactionIntent, kvGetter uiptypes.KVGetter) ([]byte, error) {
	// todo
	onChainTransaction, err := bn.Translate(intent, kvGetter)
	if err != nil {
		return nil, err
	}
	return bn.RouteRaw(intent.ChainID, onChainTransaction)
}


func (bn *BN) CheckAddress(addr []byte) bool {
	return len(addr) == 32
}

