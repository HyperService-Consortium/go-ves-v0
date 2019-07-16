package eth_client

import "testing"

func TestGetEthAccounts(t *testing.T) {
	err := NewEthClient("127.0.0.1:8545").GetEthAccounts()
	if err != nil {
		t.Error(err)
		return
	}
}

func TestUnlock(t *testing.T) {
	err := NewEthClient("127.0.0.1:8545").PersonalUnlockAccout("0x0ac45f1e6b8d47ac4c73aee62c52794b5898da9f", "123456", 600)
	if err != nil {
		t.Error(err)
		return
	}
}
