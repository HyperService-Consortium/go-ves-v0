package bni

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/HyperService-Consortium/go-ethabi"
	"github.com/HyperService-Consortium/go-rlp"
	"github.com/HyperService-Consortium/go-uip/const/value_type"
	merkleproof "github.com/HyperService-Consortium/go-uip/merkle-proof"
	"github.com/HyperService-Consortium/go-uip/uiptypes"
	ethclient "github.com/HyperService-Consortium/go-ves/lib/net/eth-client"
	"github.com/tidwall/gjson"
	"math/big"
	"net/url"
	"strconv"
)

//type Storage interface {
//
//	// blockID + color decide an only transaction on chain with id chainID
//	GetTransactionProof(chainID ChainID, blockID BlockID, color []byte) (MerkleProof, error)
//	GetStorageAt(chainID ChainID, typeID TypeID, contractAddress ContractAddress, pos []byte, description []byte) (Variable, error)
//}

func (bn *BN) GetStorageAt(chainID uiptypes.ChainID, typeID uiptypes.TypeID, contractAddress uiptypes.ContractAddress, pos []byte, description []byte) (uiptypes.Variable, error) {
	// todo
	ci, err := bn.dns.GetChainInfo(chainID)
	if err != nil {
		return nil, err
	}

	switch typeID {
	case value_type.Bool:
		s, err := ethclient.NewEthClient((&url.URL{Scheme: "http", Host: ci.GetChainHost(), Path: "/"}).String()).GetStorageAt(contractAddress, pos, "latest")
		if err != nil {
			return nil, err
		}

		b, err := hex.DecodeString(s[2:])
		if err != nil {
			return nil, err
		}
		bs, err := ethabi.NewDecoder().Decodes([]string{"bool"}, b)
		return variable{
			Type:  typeID,
			Value: bs[0],
		}, nil
	case value_type.Uint256:
		s, err := ethclient.NewEthClient((&url.URL{Scheme: "http", Host: ci.GetChainHost(), Path: "/"}).String()).GetStorageAt(contractAddress, pos, "latest")
		if err != nil {
			return nil, err
		}

		b, err := hex.DecodeString(s[2:])
		if err != nil {
			return nil, err
		}
		bs, err := ethabi.NewDecoder().Decodes([]string{"uint256"}, b)
		return variable{
			Type:  typeID,
			Value: bs[0],
		}, nil
	}

	return nil, nil
}

func (bn *BN) GetTransactionByStringHash(host string, index string) (*Transaction, error) {
	b, err := ethclient.NewEthClient(host).GetTransactionByStringHash(index)
	if err != nil {
		return nil, err
	}

	// b = bytes.Replace(b, []byte("0x"), nil, -1)
	ret := gjson.ParseBytes(b)

	if !ret.Exists() {
		return nil, errors.New("not exists")
	}

	var qwq Transaction
	var data = new(Txdata)
	qwq.data = data
	if nonce := ret.Get("nonce").String(); len(nonce) > 2 {
		data.AccountNonce, err = strconv.ParseUint(nonce[2:], 16, 64)
		if err != nil {
			return nil, err
		}
	}
	var ok bool
	if amount := ret.Get("value").String(); len(amount) > 2 {

		data.Amount, ok = new(big.Int).SetString(amount[2:], 16)
		if !ok {
			return nil, errors.New("cant conv amount")
		}
	}
	if gas := ret.Get("gas").String(); len(gas) > 2 {

		data.GasLimit, err = strconv.ParseUint(gas[2:], 16, 64)
		if err != nil {
			return nil, err
		}
	}
	if hexdata := ret.Get("input").String(); len(hexdata) > 2 {

		data.Payload, err = hex.DecodeString(hexdata[2:])
		if err != nil {
			return nil, err
		}
	}
	if price := ret.Get("gasPrice").String(); len(price) > 2 {

		data.Price, ok = new(big.Int).SetString(price[2:], 16)
		if !ok {
			return nil, errors.New("cant conv price")
		}
	}
	if r := ret.Get("r").String(); len(r) > 2 {

		data.R, ok = new(big.Int).SetString(r[2:], 16)
		if !ok {
			return nil, errors.New("cant conv R")
		}
	}
	if s := ret.Get("s").String(); len(s) > 2 {

		data.S, ok = new(big.Int).SetString(s[2:], 16)
		if !ok {
			return nil, errors.New("cant conv S")
		}
	}
	if v := ret.Get("v").String(); len(v) > 2 {

		data.V, ok = new(big.Int).SetString(v[2:], 16)
		if !ok {
			return nil, errors.New("cant conv V")
		}
	}
	if toAddress := ret.Get("to").String(); len(toAddress) > 2 {
		data.Recipient, err = hex.DecodeString(toAddress[2:])
		if err != nil {
			return nil, err
		}
	}

	// fmt.Println(hex.EncodeToString(qwq.Hash()), ret.Get("hash"))

	return &qwq, nil
}

func (bn *BN) GetTransaction(host string, index []byte) (*Transaction, error) {
	b, err := ethclient.NewEthClient(host).GetTransactionByHash(index)
	if err != nil {
		return nil, err
	}

	// b = bytes.Replace(b, []byte("0x"), nil, -1)
	ret := gjson.ParseBytes(b)

	if !ret.Exists() {
		return nil, errors.New("not exists")
	}

	var qwq Transaction
	var data = new(Txdata)
	qwq.data = data
	if nonce := ret.Get("nonce").String(); len(nonce) > 2 {
		data.AccountNonce, err = strconv.ParseUint(nonce[2:], 16, 64)
		if err != nil {
			return nil, err
		}
	}
	var ok bool
	if amount := ret.Get("value").String(); len(amount) > 2 {

		data.Amount, ok = new(big.Int).SetString(amount[2:], 16)
		if !ok {
			return nil, errors.New("cant conv amount")
		}
	}
	if gas := ret.Get("gas").String(); len(gas) > 2 {

		data.GasLimit, err = strconv.ParseUint(gas[2:], 16, 64)
		if err != nil {
			return nil, err
		}
	}
	if hexdata := ret.Get("input").String(); len(hexdata) > 2 {

		data.Payload, err = hex.DecodeString(hexdata[2:])
		if err != nil {
			return nil, err
		}
	}
	if price := ret.Get("gasPrice").String(); len(price) > 2 {

		data.Price, ok = new(big.Int).SetString(price[2:], 16)
		if !ok {
			return nil, errors.New("cant conv price")
		}
	}
	if r := ret.Get("r").String(); len(r) > 2 {

		data.R, ok = new(big.Int).SetString(r[2:], 16)
		if !ok {
			return nil, errors.New("cant conv R")
		}
	}
	if s := ret.Get("s").String(); len(s) > 2 {

		data.S, ok = new(big.Int).SetString(s[2:], 16)
		if !ok {
			return nil, errors.New("cant conv S")
		}
	}
	if v := ret.Get("v").String(); len(v) > 2 {

		data.V, ok = new(big.Int).SetString(v[2:], 16)
		if !ok {
			return nil, errors.New("cant conv V")
		}
	}
	if toAddress := ret.Get("to").String(); len(toAddress) > 2 {
		data.Recipient, err = hex.DecodeString(toAddress[2:])
		if err != nil {
			return nil, err
		}
	}

	// fmt.Println(hex.EncodeToString(qwq.Hash()), ret.Get("hash"))

	return &qwq, nil
}

func (bn *BN) GetTransactionProof(chainID uint64, blockID []byte, additional []byte) (uiptypes.MerkleProof, error) {
	cinfo, err := bn.dns.GetChainInfo(chainID)
	if err != nil {
		return nil, err
	}

	b, err := ethclient.NewEthClient(cinfo.GetChainHost()).GetBlockByHash(blockID, false)
	if err != nil {
		return nil, err
	}

	// b = bytes.Replace(b, []byte("0x"), nil, -1)
	ret := gjson.ParseBytes(b)

	if !ret.Exists() {
		return nil, errors.New("block does not not exist")
	}

	rawTxs := ret.Get("transactions").Array()

	// fmt.Println(ret.Get("transactionsRoot"), rawTxs)

	index := binary.BigEndian.Uint64(additional)

	var txs Transactions
	var tx *Transaction
	for _, rawTx := range rawTxs {
		tx, err = bn.GetTransactionByStringHash(cinfo.GetChainHost(), rawTx.String())
		if err != nil {
			return nil, err
		}
		txs = append(txs, tx)
	}

	txTrie, err := NewTxTrie(txs)
	if err != nil {
		return nil, err
	}
	hash, err := txTrie.Commit(nil)
	if err != nil {
		return nil, err
	}
	if hash.Hex() != ret.Get("transactionsRoot").String() {
		return nil, fmt.Errorf("debugging: hash.Hex()[%v] != transactionsRoot[%v]", hash.Hex(), ret.Get("transactionsRoot").String())
	}

	keybuf := new(bytes.Buffer)
	keybuf.Reset()
	rlp.Encode(keybuf, uint(index))

	proof, err := txTrie.TryProve(keybuf.Bytes())
	if err != nil {
		return nil, err
	}

	return merkleproof.NewMPTUsingKeccak256(proof, keybuf.Bytes(), txTrie.Get(keybuf.Bytes())), nil
}

func (bn *BN) GetTransactionProofByHash(chainID uint64, blockID []byte, additional []byte) (uiptypes.MerkleProof, error) {
	cinfo, err := bn.dns.GetChainInfo(chainID)
	if err != nil {
		return nil, err
	}

	b, err := ethclient.NewEthClient(cinfo.GetChainHost()).GetBlockByHash(blockID, false)
	if err != nil {
		return nil, err
	}

	// b = bytes.Replace(b, []byte("0x"), nil, -1)
	ret := gjson.ParseBytes(b)

	if !ret.Exists() {
		return nil, errors.New("block does not not exist")
	}

	rawTxs := ret.Get("transactions").Array()

	// fmt.Println(ret.Get("transactionsRoot"), rawTxs)

	var txs Transactions
	var tx *Transaction
	var index uint64
	for idx, rawTx := range rawTxs {
		tx, err = bn.GetTransactionByStringHash(cinfo.GetChainHost(), rawTx.String())
		if err != nil {
			return nil, err
		}

		if bytes.Equal(additional, tx.Hash()) {
			index = uint64(idx)
		}
		txs = append(txs, tx)
	}

	txTrie, err := NewTxTrie(txs)
	if err != nil {
		return nil, err
	}
	hash, err := txTrie.Commit(nil)
	if err != nil {
		return nil, err
	}
	if hash.Hex() != ret.Get("transactionsRoot").String() {
		return nil, fmt.Errorf("debugging: hash.Hex()[%v] != transactionsRoot[%v]", hash.Hex(), ret.Get("transactionsRoot").String())
	}

	keybuf := new(bytes.Buffer)
	keybuf.Reset()
	err = rlp.Encode(keybuf, uint(index))
	if err != nil {
		return nil, err
	}

	proof, err := txTrie.TryProve(keybuf.Bytes())
	if err != nil {
		return nil, err
	}

	return merkleproof.NewMPTUsingKeccak256(proof, keybuf.Bytes(), txTrie.Get(keybuf.Bytes())), nil
}
