package b

import (
	"encoding/json"
	"fmt"
	"testing"
)

type EthAccount struct {
	Address    string `json:"address"`
	ChainID    uint64 `json:"chain_id"`
	PassPhrase string `json:"pass_phrase"`
}

type EthAccountAlias struct {
	EthAccount
	Alias string `json:"alias"`
}

var jsonString = `
[
{
    "address": "0x666",
    "chain_id": 3,
    "pass_phrase": "123456",
    "alias": "myriad dreamin"
},
{
    "address": "0x666",
    "chain_id": 3,
    "pass_phrase": "123456",
    "alias": "myriad dreamin"
}
]
`

func TestJson(t *testing.T) {
	var ks = make([]*EthAccountAlias, 0)
	err := json.Unmarshal([]byte(jsonString), &ks)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(ks[0], ks[1])
}
