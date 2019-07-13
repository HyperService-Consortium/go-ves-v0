package opintent

import (
	"encoding/json"
	"errors"
	"fmt"

	types "github.com/Myriad-Dreamin/go-ves/types"

	chain_type "github.com/Myriad-Dreamin/go-ves/go-uiputils/types/chain_type"
	unit_type "github.com/Myriad-Dreamin/go-ves/go-uiputils/types/unit_type"
	account "github.com/Myriad-Dreamin/go-ves/types/account"
)

type hexstring = string

type RawAccountInfo struct {
	ChainId uint64 `json:"domain"`
	Name    string `json:"user_name"`
}

type BasePaymentOpIntent struct {
	Src        *RawAccountInfo `json:"src"`    // key
	Dst        *RawAccountInfo `json:"dst"`    // key
	Amount     hexstring       `json:"amount"` // key
	UnitString string          `json:"unit"`   // optional
}

func initializeError(keyAttr string) error {
	return fmt.Errorf("the attribute %v must be included in the Payment intent", keyAttr)
}

type ChainInfo struct {
	Host      string
	ChainType uint64
}

func (c *ChainInfo) GetHost() string {
	return c.Host
}

func (c *ChainInfo) GetChainType() uint64 {
	return c.ChainType
}

type ChainInfoInterface interface {
	GetHost() string
	GetChainType() uint64
}

type DataType interface {
	GetValue() interface{}
	GetDesc() []byte
}

type help_info string
type chain_id uint64
type meta []byte
type Encoder interface {
}

type Decoder interface {
}

type ProcessorInterface interface {
	CheckAddress([]byte) bool
	// ParseData([]DataType, meta, chain_id) (bool, []byte, help_info)
	// GetEncoder() Encoder
	// GetDecoder() Decoder
}

func SearchChainId(domain uint64) (ChainInfoInterface, error) {
	switch domain {
	case 0:
		return nil, errors.New("nil domain is not allowed")
	case 1: // ethereum chain 1
		return &ChainInfo{
			Host:      "127.0.0.1",
			ChainType: chain_type.Ethereum,
		}, nil
	case 2: // ethereum chain 2
		return &ChainInfo{
			Host:      "127.0.0.1",
			ChainType: chain_type.Ethereum,
		}, nil
	case 3: // tendermint chain 1
		return &ChainInfo{
			Host:      "47.254.66.11",
			ChainType: chain_type.TendermintNSB,
		}, nil
	case 4: // ethereum chain 1
		return &ChainInfo{
			Host:      "47.251.2.73",
			ChainType: chain_type.TendermintNSB,
		}, nil
	default:
		return nil, errors.New("not found")
	}
}

func TempGetRelay(domain uint64) (types.Account, error) {
	switch domain {
	case 0:
		return nil, errors.New("nil domain is not allowed")
	case 1: // ethereum chain 1
		return &account.PureAccount{
			ChainId: 1,
			Address: []byte("abc"),
		}, nil
	case 2: // ethereum chain 2
		return &account.PureAccount{
			ChainId: 2,
			Address: []byte("abc"),
		}, nil
	case 3: // tendermint chain 1
		return &account.PureAccount{
			ChainId: 3,
			Address: []byte("abc"),
		}, nil
	case 4: // ethereum chain 1
		return &account.PureAccount{
			ChainId: 4,
			Address: []byte("abc"),
		}, nil
	default:
		return nil, errors.New("not found")
	}
}

func TempSearchAccount(name string, chainId uint64) (types.Account, error) {
	switch name {
	case "":
		return nil, errors.New("nil name is not allowed")
	case "a1": // ethereum chain 1
		switch chainId {
		case 0:
			return nil, errors.New("nil domain is not allowed")
		case 1: // ethereum chain 1
			return &account.PureAccount{
				ChainId: 1,
				Address: []byte("abc"),
			}, nil
		case 2: // ethereum chain 2
			return &account.PureAccount{
				ChainId: 1,
				Address: []byte("abc"),
			}, nil
		case 3: // tendermint chain 1
			return &account.PureAccount{
				ChainId: 1,
				Address: []byte("abc"),
			}, nil
		case 4: // ethereum chain 1
			return &account.PureAccount{
				ChainId: 1,
				Address: []byte("abc"),
			}, nil
		default:
			return nil, errors.New("not found")
		}
	case "a2": // ethereum chain 1
		switch chainId {
		case 0:
			return nil, errors.New("nil domain is not allowed")
		case 1: // ethereum chain 1
			return &account.PureAccount{
				ChainId: 1,
				Address: []byte("abc"),
			}, nil
		case 2: // ethereum chain 2
			return &account.PureAccount{
				ChainId: 1,
				Address: []byte("abc"),
			}, nil
		case 3: // tendermint chain 1
			return &account.PureAccount{
				ChainId: 1,
				Address: []byte("abc"),
			}, nil
		case 4: // ethereum chain 1
			return &account.PureAccount{
				ChainId: 1,
				Address: []byte("abc"),
			}, nil
		default:
			return nil, errors.New("not found")
		}
	default:
		return nil, errors.New("not registered account")
	}
}

func (ier *OpIntentInitializer) InitPaymentOpIntent(
	name string,
	subIntent json.RawMessage,
) (txs []*TransactionIntent, err error) {
	var paymentIntent BasePaymentOpIntent
	err = json.Unmarshal(subIntent, &paymentIntent)
	var tx *TransactionIntent
	if err != nil {
		return
	}
	if paymentIntent.Src == nil {
		return nil, initializeError("src")
	}
	if paymentIntent.Dst == nil {
		return nil, initializeError("dst")
	}
	if len(paymentIntent.Amount) == 0 {
		return nil, initializeError("amount")
	}
	if t, ok := unit_type.Mapping[paymentIntent.UnitString]; !ok {
		return nil, errors.New("unknown unit type")
	} else {
		var srcInfo, dstInfo types.Account
		srcInfo, err = TempSearchAccount(paymentIntent.Src.Name, paymentIntent.Src.ChainId)
		if err != nil {
			return nil, err
		}
		dstInfo, err = TempSearchAccount(paymentIntent.Dst.Name, paymentIntent.Dst.ChainId)
		if err != nil {
			return nil, err
		}
		if tx, err = ier.genPayment(srcInfo, nil, paymentIntent.Amount, t); err != nil {
			return
		}
		txs = append(txs, tx)
		if tx, err = ier.genPayment(nil, dstInfo, paymentIntent.Amount, t); err != nil {
			return
		}
		txs = append(txs, tx)
		// cinfo, err = SearchChainId(paymentIntent.Src.ChainId)
		// if err != nil {
		// 	return nil, err
		// }
		// var processor ProcessorInterface
		// switch cinfo.GetChainType() {
		// case chain_type.Ethereum:
		// 	processor = eth_processor.GetProcessor()
		// default:
		// 	return nil, errors.New("unsupport chain_type")
		// }
		// if !processor.CheckAddress(paymentIntent.Src.)
		return
	}
}

// type PaymentMeta struct {
// 	OpType string `json:"op_type"`
// }

// var pm = []byte(`{"op_type": "transact"}`)

func (ier *OpIntentInitializer) genPayment(src types.Account, dst types.Account, amt string, utid unit_type.Type) (tx *TransactionIntent, err error) {
	if src == nil {
		if src, err = TempGetRelay(dst.GetChainId()); err != nil {
			return nil, err
		}
	} else {
		if dst, err = TempGetRelay(src.GetChainId()); err != nil {
			return nil, err
		}
	}

	tx = &TransactionIntent{
		Src:       src.GetAddress(),
		Dst:       dst.GetAddress(),
		TransType: 0,
		Amt:       amt,
		ChainId:   dst.GetChainId(),
	}
	return
}