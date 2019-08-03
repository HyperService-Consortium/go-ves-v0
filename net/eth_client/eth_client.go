package eth_client

import (
	"encoding/json"
	"fmt"

	jsonrpc_client "github.com/Myriad-Dreamin/go-ves/net/rpc-client"
	"github.com/tidwall/gjson"

	jsonobj "github.com/Myriad-Dreamin/go-ves/net/eth_client/jsonobj"
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

	b, err := eth.JsonRpcClient.PostWithBody(jsonobj.GetAccount())
	if err != nil {
		return nil, err
	}

	var x []string
	err = json.Unmarshal(b, &x)
	if err != nil {
		return nil, err
	}

	return x, err
}

func (eth *EthClient) PersonalUnlockAccout(addr string, passphrase string, duration int) (bool, error) {
	b := jsonobj.GetPersonalUnlock(addr, passphrase, duration)
	bb, err := eth.JsonRpcClient.PostWithBody(b)
	jsonobj.ReturnBytes(b)
	if err != nil {
		return false, err
	}

	return gjson.ParseBytes(bb).Bool(), err
}

func (eth *EthClient) SendTransaction(obj []byte) (string, error) {
	b := jsonobj.GetSendTransaction(obj)
	bb, err := eth.JsonRpcClient.PostWithBody(b)
	jsonobj.ReturnBytes(b)
	if err != nil {
		return "", err
	}

	return gjson.ParseBytes(bb).String(), err
}

func (eth *EthClient) GetStorageAt(contractAddress, pos []byte, tag string) (string, error) {
	b := jsonobj.GetStorageAt(contractAddress, pos, tag)
	fmt.Println(string(b))
	bb, err := eth.JsonRpcClient.PostWithBody(b)
	jsonobj.ReturnBytes(b)
	if err != nil {
		return "", err
	}

	return gjson.ParseBytes(bb).String(), err
}

func Do(url string, jsonBody []byte) ([]byte, error) {
	return jsonrpc_client.Do(url, jsonBody)
}
