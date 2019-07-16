package eth_client

import (
	"encoding/json"
	"fmt"
	"reflect"

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

func (eth *EthClient) GetEthAccounts() error {
	b, err := eth.PostRequestWithJsonObj(map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "eth_accounts",
		"params":  make([]interface{}, 0, 0),
		"id":      1,
	})
	if err != nil {
		return err
	}
	fmt.Println(string(b), reflect.TypeOf(b))

	var x []string
	err = json.Unmarshal(b, &x)
	fmt.Println(err, x)

	return err
}

func (eth *EthClient) PersonalUnlockAccout(addr string, passphrase string, duration int) error {
	b, err := eth.PostRequestWithJsonObj(map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "personal_unlockAccount",
		"params":  []interface{}{addr, passphrase, duration},
		"id":      64,
	})
	if err != nil {
		return err
	}
	fmt.Println(string(b), reflect.TypeOf(b))

	return err
}
