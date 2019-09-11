package bni

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"sync/atomic"
	"time"

	trie "github.com/HyperService-Consortium/go-mpt"
	"github.com/HyperService-Consortium/go-rlp"
	merkleproof "github.com/HyperService-Consortium/go-uip/merkle-proof"
	uiptypes "github.com/HyperService-Consortium/go-uip/types"
	ethclient "github.com/HyperService-Consortium/go-ves/lib/net/eth-client"
	leveldb "github.com/syndtr/goleveldb/leveldb"
	leveldbstorage "github.com/syndtr/goleveldb/leveldb/storage"
	"golang.org/x/crypto/sha3"

	chaininfo "github.com/HyperService-Consortium/go-uip/temporary-chain-info"

	gjson "github.com/tidwall/gjson"
)

type Transaction struct {
	data *Txdata
	// caches
	hash atomic.Value
	size atomic.Value
	from atomic.Value
}

type Txdata struct {
	AccountNonce uint64   `json:"nonce"    gencodec:"required"`
	Price        *big.Int `json:"gasPrice" gencodec:"required"`
	GasLimit     uint64   `json:"gas"      gencodec:"required"`
	Recipient    []byte   `json:"to"       rlp:"nil"` // nil means contract creation
	Amount       *big.Int `json:"value"    gencodec:"required"`
	Payload      []byte   `json:"input"    gencodec:"required"`

	// Signature values
	V *big.Int `json:"v" gencodec:"required"`
	R *big.Int `json:"r" gencodec:"required"`
	S *big.Int `json:"s" gencodec:"required"`

	// This is only used when marshaling to JSON.
	Hash []byte `json:"hash" rlp:"-"`
}

type DerivableList interface {
	Len() int
	GetRlp(i int) []byte
}

// Transactions is a Transaction slice type for basic sorting.
type Transactions []*Transaction

// Len returns the length of s.
func (s Transactions) Len() int { return len(s) }

// GetRlp implements Rlpable and returns the i'th element of s in rlp.
func (s Transactions) GetRlp(i int) []byte {
	enc, _ := rlp.EncodeToBytes(s[i].data)

	return enc
}

func rlpHash(x interface{}) []byte {
	hw := sha3.NewLegacyKeccak256()
	// WARNING: ignoring errors
	rlp.Encode(hw, x)
	return hw.Sum(nil)
}

func (tx *Transaction) Hash() []byte {
	if hash := tx.hash.Load(); hash != nil {
		return hash.([]byte)
	}
	v := rlpHash(tx.data)
	tx.hash.Store(v)
	return v
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

var emptyHash = trie.HexToHash("0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421")
var __op, _ = leveldb.Open(leveldbstorage.NewMemStorage(), nil)
var __v, _ = trie.NewNodeBasefromDB(__op)

func NewVoidTrie() (*trie.Trie, error) {
	return trie.NewTrie(emptyHash, __v)
}

func NewTxTrie(list DerivableList) (*trie.Trie, error) {
	keybuf := new(bytes.Buffer)
	txTrie, err := NewVoidTrie()

	if err != nil {
		return nil, err
	}
	for i := 0; i < list.Len(); i++ {
		keybuf.Reset()
		rlp.Encode(keybuf, uint(i))
		txTrie.Update(keybuf.Bytes(), list.GetRlp(i))
	}
	return txTrie, nil
}

func (bn *BN) GetTransactionProof(chainID uint64, blockID []byte, additional []byte) (uiptypes.MerkleProof, error) {
	cinfo, err := chaininfo.SearchChainId(chainID)
	if err != nil {
		return nil, err
	}

	b, err := ethclient.NewEthClient(cinfo.GetHost()).GetBlockByHash(blockID, false)
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
		tx, err = bn.GetTransactionByStringHash(cinfo.GetHost(), rawTx.String())
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

func (bn *BN) WaitForTransact(chainID uint64, receipt []byte, opt *uiptypes.WaitOption) ([]byte, []byte, error) {
	cinfo, err := chaininfo.SearchChainId(chainID)
	if err != nil {
		return nil, nil, err
	}
	worker := ethclient.NewEthClient(cinfo.GetHost())
	ddl := time.Now().Add(opt.Timeout)
	for time.Now().Before(ddl) {
		tx, err := worker.GetTransactionByHash(receipt)
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

func (bn *BN) GetTransactionProofByHash(chainID uint64, blockID []byte, additional []byte) (uiptypes.MerkleProof, error) {
	cinfo, err := chaininfo.SearchChainId(chainID)
	if err != nil {
		return nil, err
	}

	b, err := ethclient.NewEthClient(cinfo.GetHost()).GetBlockByHash(blockID, false)
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
		tx, err = bn.GetTransactionByStringHash(cinfo.GetHost(), rawTx.String())
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
	rlp.Encode(keybuf, uint(index))

	proof, err := txTrie.TryProve(keybuf.Bytes())
	if err != nil {
		return nil, err
	}

	return merkleproof.NewMPTUsingKeccak256(proof, keybuf.Bytes(), txTrie.Get(keybuf.Bytes())), nil
}
