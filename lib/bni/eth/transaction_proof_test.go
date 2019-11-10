package bni

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/HyperService-Consortium/go-rlp"
)

func TestGetTransaction(t *testing.T) {
	const txString = "a41d03fde4e7cf4c58870092c65709db7532956f7d0882156f11f503a6d88d2f"
	txbytes, err := hex.DecodeString(txString)
	if err != nil {
		t.Error(err)
		return
	}

	// fmt.Println(new(BN).GetTransaction("127.0.0.1:8545", txbytes))
	new(BN).GetTransaction("127.0.0.1:8545", txbytes)
}

func TestGetTransactionProofByIndex(t *testing.T) {
	const blockString = "8a8b9aaa48e0fb024abb7105798ad48057cf4fd14100505addabc319ed3d41c6"
	blockbytes, err := hex.DecodeString(blockString)
	if err != nil {
		t.Error(err)
		return
	}
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(3))

	fmt.Println(new(BN).GetTransactionProof(1, blockbytes, buf))
}

func TestGetTransactionProof(t *testing.T) {
	const txString = "cc507383e3c78f95c2551eb74ac2a216bae80eaf977ce68264351d7378a83b0e"
	txbytes, err := hex.DecodeString(txString)
	if err != nil {
		t.Error(err)
		return
	}

	blockBytes, index, err := new(BN).WaitForTransact(1, txbytes)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(blockBytes)

	qwq, err := new(BN).GetTransactionProof(1, blockBytes, index)
	if err != nil {
		t.Error(err)
		return
	}
	var tx Txdata
	err = rlp.DecodeBytes(qwq.GetValue(), &tx)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(tx)
}
