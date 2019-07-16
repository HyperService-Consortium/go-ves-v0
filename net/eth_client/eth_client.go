package eth_client

import (
	"encoding/json"

	jsonrpc_client "github.com/Myriad-Dreamin/go-ves/net/rpc-client"
)

type EthClient struct {
	*jsonrpc_client.JsonRpcClient
}

func NewEthClient(host string) *EthClient {
	return &EthClient{
		JsonRpcClient: jsonrpc_client.NewJsonRpcClient(host),
	}
}

type eth_params struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Params  interface{} `json:"params,omitempty"`
	Method  string      `json:"method,omitempty"`
}

func (eth *EthClient) GetEthAccounts() ([]string, error) {
	b, err := eth.PostRequestWithJsonObj(map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "eth_accounts",
		"params":  make([]interface{}, 0, 0),
		"id":      1,
	})
	if err != nil {
		return nil, err
	}

	var x []string
	err = json.Unmarshal(b, &x)

	return x, err
}

func (eth *EthClient) PersonalUnlockAccout(addr string, passphrase string, duration int) (bool, error) {
	b, err := eth.PostRequestWithJsonObj(map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "personal_unlockAccount",
		"params":  []interface{}{addr, passphrase, duration},
		"id":      64,
	})
	if err != nil {
		return false, err
	}

	var x bool
	err = json.Unmarshal(b, &x)
	if err != nil {
		return false, err
	}

	return x, err
}
