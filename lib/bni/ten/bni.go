package bni

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/url"

	appl "github.com/HyperServiceOne/NSB/application"
	cmn "github.com/HyperServiceOne/NSB/common"
	nsbmath "github.com/HyperServiceOne/NSB/math"
	"github.com/Myriad-Dreamin/go-ethabi"
	TransType "github.com/Myriad-Dreamin/go-uip/const/trans_type"
	valuetype "github.com/Myriad-Dreamin/go-uip/const/value_type"
	opintent "github.com/Myriad-Dreamin/go-uip/op-intent"
	"github.com/Myriad-Dreamin/go-uip/types"
	ethclient "github.com/Myriad-Dreamin/go-ves/lib/net/eth-client"
	nsbclient "github.com/Myriad-Dreamin/go-ves/lib/net/nsb-client"
)

type BN struct {
	signer types.Signer
}

type MiddleHeader struct {
	Header    cmn.TransactionHeader `json:"h"`
	PreHeader []byte                `json:"p"`
}

func (bn *BN) MustWithSigner() bool {
	return true
}

func (bn *BN) RouteWithSigner(signer types.Signer) types.Router {
	nbn := new(BN)
	nbn.signer = signer
	return nbn
}

func (bn *BN) RouteRaw(destination uint64, payload []byte) ([]byte, error) {
	ci, err := SearchChainId(destination)
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

	buf.Write(txHeader.Header.From)
	buf.Write(txHeader.Header.ContractAddress)
	buf.Write(txHeader.Header.Data)
	buf.Write(txHeader.Header.Value.Bytes())
	buf.Write(txHeader.Header.Nonce.Bytes())
	txHeader.Header.Signature = bn.signer.Sign(buf.Bytes()).Bytes()

	if err != nil {
		return nil, err
	}

	b, err := json.Marshal(txHeader.Header)
	if err != nil {
		return nil, err
	}

	buf.Reset()
	buf.Write(txHeader.PreHeader)
	buf.Write(b)

	return nsbclient.NewNSBClient((&url.URL{Scheme: "http", Host: ci.GetHost(), Path: "/"}).String()).BroadcastTxCommitReturnBytes(buf.Bytes())
}

func (bn *BN) Route(intent *types.TransactionIntent, kvGetter types.KVGetter) ([]byte, error) {
	// todo
	onChainTransaction, err := bn.Translate(intent, kvGetter)
	if err != nil {
		return nil, err
	}
	return bn.RouteRaw(intent.ChainID, onChainTransaction)
}

func (bn *BN) Translate(
	intent *opintent.TransactionIntent,
	kvGetter types.KVGetter,
) ([]byte, error) {
	switch intent.TransType {
	case TransType.Payment:
		var txHeader MiddleHeader

		// Nonce
		nonce := make([]byte, 32)
		_, err := rand.Read(nonce)
		if err != nil {
			return nil, err
		}
		txHeader.Header.Nonce = nsbmath.NewUint256FromBytes(nonce)

		// Data
		var args appl.ArgsTransfer
		value, err := hex.DecodeString(intent.Amt)

		if err != nil {
			return nil, err
		}
		args.Value = nsbmath.NewUint256FromBytes(value)

		b, err := json.Marshal(args)
		if err != nil {
			return nil, err
		}

		var fap appl.FAPair
		fap.FuncName = "transfer"
		fap.Args = b
		txHeader.Header.Data, err = json.Marshal(fap)
		if err != nil {
			return nil, err
		}

		// Rest
		txHeader.Header.ContractAddress = intent.Dst
		txHeader.Header.From = intent.Src
		txHeader.Header.Value = args.Value
		txHeader.PreHeader = []byte("systemCall\x19system.token\x18")

		return json.Marshal(txHeader)
	case TransType.ContractInvoke:
		// var meta types.ContractInvokeMeta
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

func (bn *BN) CheckAddress(addr []byte) bool {
	return len(addr) == 32
}

func (bn *BN) GetStorageAt(chainID uint64, typeID uint16, contractAddress []byte, pos []byte, desc []byte) (interface{}, error) {

	ci, err := SearchChainId(chainID)
	if err != nil {
		return nil, err
	}

	switch typeID {
	case valuetype.Bool:
		s, err := ethclient.NewEthClient((&url.URL{Scheme: "http", Host: ci.GetHost(), Path: "/"}).String()).GetStorageAt(contractAddress, pos, "latest")
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
		s, err := ethclient.NewEthClient(ci.GetHost()).GetStorageAt(contractAddress, pos, "latest")
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
