package ethclient

import (
	"encoding/json"
	"github.com/tidwall/gjson"

	"github.com/HyperService-Consortium/go-ves/lib/net/eth-client/jsonobj"
	jsonrpc_client "github.com/HyperService-Consortium/go-ves/lib/net/rpc-client"
)

// EthClient provide interface to ethereum rpc service
type EthClient struct {
	*jsonrpc_client.JsonRpcClient
}

// NewEthClient return a pointer of eth-client object
func NewEthClient(host string) *EthClient {
	return &EthClient{
		JsonRpcClient: jsonrpc_client.NewJsonRpcClient(host),
	}
}

//GetEthAccounts return accounts as strings
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

// PersonalUnlockAccout return whether the account is unlocked successfully
func (eth *EthClient) PersonalUnlockAccout(addr string, passphrase string, duration int) (bool, error) {
	b := jsonobj.GetPersonalUnlock(addr, passphrase, duration)
	bb, err := eth.JsonRpcClient.PostWithBody(b)
	jsonobj.ReturnBytes(b)
	if err != nil {
		return false, err
	}

	return gjson.ParseBytes(bb).Bool(), err
}

// SendTransaction return the receipt of sending transaction
func (eth *EthClient) SendTransaction(obj []byte) (string, error) {
	b := jsonobj.GetSendTransaction(obj)
	bb, err := eth.JsonRpcClient.PostWithBody(b)
	jsonobj.ReturnBytes(b)
	if err != nil {
		return "", err
	}

	return gjson.ParseBytes(bb).String(), err
}

// GetStorageAt return the value of position on the address
func (eth *EthClient) GetStorageAt(contractAddress, pos []byte, tag string) (string, error) {
	b := jsonobj.GetStorageAt(contractAddress, pos, tag)
	bb, err := eth.JsonRpcClient.PostWithBody(b)
	jsonobj.ReturnBytes(b)
	if err != nil {
		return "", err
	}

	return gjson.ParseBytes(bb).String(), err
}

func (eth *EthClient) GetTransactionByHash(transactionHash []byte) ([]byte, error) {
	b := jsonobj.GetTransactionByHash(transactionHash)
	bb, err := eth.JsonRpcClient.PostWithBody(b)

	jsonobj.ReturnBytes(b)
	if err != nil {
		return nil, err
	}
	return bb, nil
}

func (eth *EthClient) GetTransactionByStringHash(transactionHash string) ([]byte, error) {
	b := jsonobj.GetTransactionByStringHash(transactionHash)
	bb, err := eth.JsonRpcClient.PostWithBody(b)

	jsonobj.ReturnBytes(b)
	if err != nil {
		return nil, err
	}
	return bb, nil
}

func (eth *EthClient) GetBlockByHash(blockHash []byte, returnFull bool) ([]byte, error) {
	b := jsonobj.GetBlockByHash(blockHash, returnFull)
	bb, err := eth.JsonRpcClient.PostWithBody(b)

	jsonobj.ReturnBytes(b)
	if err != nil {
		return nil, err
	}
	return bb, nil
}

// Do raw rpc invocation
func Do(url string, jsonBody []byte) ([]byte, error) {
	return jsonrpc_client.Do(url, jsonBody)
}
