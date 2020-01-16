package bni

import (
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/HyperService-Consortium/go-uip/uiptypes"
	ethclient "github.com/HyperService-Consortium/go-ves/lib/net/eth-client"
	"github.com/tidwall/gjson"
	"strconv"
	"time"
)

func (bn *BN) MustWithSigner() bool {
	return true
}

func (bn *BN) RouteWithSigner(signer uiptypes.Signer) (uiptypes.Router, error) {
	nbn :=  *bn
	nbn.signer = signer
	return &nbn, nil
}

func (bn *BN) RouteRaw(destination uiptypes.ChainID, rawTransaction uiptypes.RawTransaction) (
	transactionReceipt uiptypes.TransactionReceipt, err error) {
	if !rawTransaction.Signed() {
		return nil, ErrNotSigned
	}

	return bn.createTransactionReceipt(
		bn.sendTransaction(destination, rawTransaction))
}

func (bn *BN) sendTransaction(
	destination uiptypes.ChainID, rawTransaction uiptypes.RawTransaction) (
	[]byte, error) {
	ci, err := bn.dns.GetChainInfo(destination)
	if err != nil {
		return nil, err
	}
	// todo receipt
	b, err := rawTransaction.Serialize()
	if err != nil {
		return nil, err
	}
	return ethclient.HTTPDo(ci.GetChainHost(), b)
}

func (bn *BN) createTransactionReceipt(b []byte, err error) (
	uiptypes.TransactionReceipt, error) {
	if err != nil {
		return nil, err
	}
	var x string
	err = json.Unmarshal(b, &x)
	if err != nil {
		return nil, err
	}

	if len(x) == 0 {
		return nil, errors.New("deploy failed?")
	}

	b, err = hex.DecodeString(x[2:])
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (bn *BN) WaitForTransact(chainID uiptypes.ChainID, transactionReceipt uiptypes.TransactionReceipt,
	rOptions ...interface{}) (blockID uiptypes.BlockID, color []byte, err error) {
	options := parseOptions(rOptions)
	chainInfo, err := bn.dns.GetChainInfo(chainID)
	if err != nil {
		return nil, nil, err
	}
	worker := ethclient.NewEthClient(chainInfo.GetChainHost())
	ddl := time.Now().Add(options.timeout)
	for time.Now().Before(ddl) {
		tx, err := worker.GetTransactionByHash(transactionReceipt)
		if err != nil {
			return nil, nil, err
		}
		fmt.Println(string(tx))
		if gjson.GetBytes(tx, "blockNumber").Type != gjson.Null {
			b, _ := hex.DecodeString(gjson.GetBytes(tx, "blockHash").String()[2:])
			idx, _ := strconv.ParseUint(gjson.GetBytes(tx, "transactionIndex").String()[2:], 16, 64)
			var a = make([]byte, 8)
			binary.BigEndian.PutUint64(a, idx)
			return b, a, nil
		}
		time.Sleep(time.Millisecond * 500)
	}
	return nil, nil, errors.New("timeout")
}

func (bn *BN) Route(intent *uiptypes.TransactionIntent, storage uiptypes.Storage) ([]byte, error) {
	// todo
	rawTransaction, err := bn.Translate(intent, storage)
	if err != nil {
		return nil, err
	}
	if !rawTransaction.Signed() {
		rawTransaction, err = rawTransaction.Sign(bn.signer)
		if err != nil {
			return nil, err
		}
	}
	return bn.RouteRaw(intent.ChainID, rawTransaction)
}

