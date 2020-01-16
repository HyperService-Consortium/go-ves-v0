package bni

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	base_raw_transaction "github.com/HyperService-Consortium/go-uip/base-raw-transaction"
	"github.com/HyperService-Consortium/go-uip/const/trans_type"
	"github.com/HyperService-Consortium/go-uip/uiptypes"
	payment_option "github.com/HyperService-Consortium/go-ves/lib/bni/payment-option"
	"github.com/tidwall/gjson"
)

//type Translator interface {
//	Translate(intent *TransactionIntent, storage Storage) (rawTransaction RawTransaction, err error)
//
//	// reflect.DeepEqual(Deserialize(rawTransaction.Byte()), rawTransaction) == true
//	Deserialize(raw []byte) (rawTransaction RawTransaction, err error)
//}

func (bn *BN) Translate(intent *uiptypes.TransactionIntent, storage uiptypes.Storage) (uiptypes.RawTransaction, error) {
	switch intent.TransType {
	case trans_type.Payment:
		meta := gjson.ParseBytes(intent.Meta)
		value, err := payment_option.ParseInconsistentValueOption(meta, storage, intent.Amt)
		if err != nil {
			return nil, err
		}

		//fmt.Println(value, ".........")

		b, err := json.Marshal(map[string]interface{}{
			"jsonrpc": "2.0",
			"method":  "eth_sendTransaction",
			"params": []interface{}{
				map[string]interface{}{
					"from":  decoratePrefix(hex.EncodeToString(intent.Src)),
					"to":    decoratePrefix(hex.EncodeToString(intent.Dst)),
					"value": decorateValuePrefix(value),
				},
			},
			"id": 1,
		})
		//fmt.Println("...", string(b))
		return base_raw_transaction.Transaction(b), err
	case trans_type.ContractInvoke:
		var meta uiptypes.ContractInvokeMeta
		err := json.Unmarshal(intent.Meta, &meta)
		if err != nil {
			return nil, err
		}
		//_ = meta
		// todo, test
		data, err := ContractInvocationDataABI(intent.ChainID, &meta, storage)
		if err != nil {
			return nil, err
		}

		hexdata := hex.EncodeToString(data)
		// meta.FuncName

		b, err := json.Marshal(map[string]interface{}{
			"jsonrpc": "2.0",
			"method":  "eth_sendTransaction",
			"params": []interface{}{
				map[string]interface{}{
					"from":  decoratePrefix(hex.EncodeToString(intent.Src)),
					"to":    decoratePrefix(hex.EncodeToString(intent.Dst)),
					// todo
					//"value": decoratePrefix(intent.Amt),
					"data":  decorateValuePrefix(hexdata),
				},
			},
			"id": 1,
		})
		return base_raw_transaction.Transaction(b), err
	default:
		return nil, errors.New("cant translate")
	}
}

func (bn *BN) Deserialize(raw []byte) (rawTransaction uiptypes.RawTransaction, err error) {
	// todo
	return base_raw_transaction.Transaction(raw), nil
}
