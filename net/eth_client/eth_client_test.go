package eth_client

import (
	"fmt"
	"testing"
)

const (
	test_host = "127.0.0.1:8545"
)

func TestGetEthAccounts(t *testing.T) {
	x, err := NewEthClient(test_host).GetEthAccounts()
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(x)
}

func TestUnlock(t *testing.T) {
	ok, err := NewEthClient(test_host).PersonalUnlockAccout("0x0ac45f1e6b8d47ac4c73aee62c52794b5898da9f", "123456", 600)
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
