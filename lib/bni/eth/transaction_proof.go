package bni

import (
	"bytes"
	"math/big"
	"sync/atomic"

	trie "github.com/HyperService-Consortium/go-mpt"
	"github.com/HyperService-Consortium/go-rlp"
	leveldb "github.com/syndtr/goleveldb/leveldb"
	leveldbstorage "github.com/syndtr/goleveldb/leveldb/storage"
	"golang.org/x/crypto/sha3"
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
