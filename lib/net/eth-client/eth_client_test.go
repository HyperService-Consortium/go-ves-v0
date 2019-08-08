package ethclient

import (
	"encoding/hex"
	"fmt"
	"testing"
)

const (
	testHost = "127.0.0.1:8545"
)

func TestGetEthAccounts(t *testing.T) {
	x, err := NewEthClient(testHost).GetEthAccounts()
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(x)
}

func TestUnlock(t *testing.T) {
	ok, err := NewEthClient(testHost).PersonalUnlockAccout("0x0ac45f1e6b8d47ac4c73aee62c52794b5898da9f", "123456", 600)

	if ok == false || err != nil {
		if ok == false {
			if err != nil {

				t.Error(err)
			} else {
				t.Errorf("not ok..")
			}
		} else {

			t.Error(err)
		}
		return
	}
}

const objjj = `{"from":"0x0ac45f1e6b8d47ac4c73aee62c52794b5898da9f", "to": "0x981739a13593980763de3353340617ef16da6354", "value": "0x1"}`

func TestSendTransaction(t *testing.T) {
	b, err := NewEthClient(testHost).SendTransaction([]byte(objjj))
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(string(b))
}

func TestGetStorageAt(t *testing.T) {
	var addr = "1234567812345678123456781234567812345678"
	baddr, err := hex.DecodeString(addr)
	if err != nil {
		t.Error(err)
		return
	}
	var pos = []byte{1}
	b, err := NewEthClient(testHost).GetStorageAt(baddr, pos, "latest")
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(string(b))
}

func TestGetTransactionByHash(t *testing.T) {
	txb, err := hex.DecodeString("a41d03fde4e7cf4c58870092c65709db7532956f7d0882156f11f503a6d88d2f")
	if err != nil {
		t.Error(err)
		return
	}
	b, err := NewEthClient(testHost).GetTransactionByHash(txb)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(string(b))
	b, err = NewEthClient(testHost).GetTransactionByStringHash("0xa41d03fde4e7cf4c58870092c65709db7532956f7d0882156f11f503a6d88d2f")
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(string(b))

}

func TestGetBlockByHash(t *testing.T) {
	txb, err := hex.DecodeString("8a8b9aaa48e0fb024abb7105798ad48057cf4fd14100505addabc319ed3d41c6")
	if err != nil {
		t.Error(err)
		return
	}

	b, err := NewEthClient(testHost).GetBlockByHash(txb, true)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(string(b))

	b, err = NewEthClient(testHost).GetBlockByHash(txb, false)
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println(string(b))
}
