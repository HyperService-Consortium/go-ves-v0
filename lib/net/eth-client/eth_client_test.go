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
