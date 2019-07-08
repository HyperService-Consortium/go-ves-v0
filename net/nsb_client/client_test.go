package nsbcli

import (
	"encoding/json"
	"fmt"
	"testing"

	mtest "github.com/Myriad-Dreamin/mydrest"
)

var s mtest.TestHelper

const host = "47.251.2.73:26657"

func PretiJson(minterface interface{}) string {
	je, _ := json.MarshalIndent(minterface, "", "\t")
	return string(je)
}

func TestGetAbciInfo(t *testing.T) {
	b, err := NewNSBClient(host).GetAbciInfo()
	s.AssertNoErr(t, err)
	fmt.Println(PretiJson(b))
}

func TestGetBlock(t *testing.T) {
	b, err := NewNSBClient(host).GetBlock(8200)
	s.AssertNoErr(t, err)
	fmt.Println(PretiJson(b))
}

func TestGetCommitInfo(t *testing.T) {
	b, err := NewNSBClient(host).GetCommitInfo(8200)
	s.AssertNoErr(t, err)
	fmt.Println(PretiJson(b))
}

func TestGetBlocks(t *testing.T) {
	b, err := NewNSBClient(host).GetBlocks(8199, 8200)
	s.AssertNoErr(t, err)
	fmt.Println(PretiJson(b))
}

func TestGetBlockResults(t *testing.T) {
	b, err := NewNSBClient(host).GetBlockResults(8200)
	s.AssertNoErr(t, err)
	fmt.Println(PretiJson(b))
}

func TestGetConsensusParamsInfo(t *testing.T) {
	b, err := NewNSBClient(host).GetConsensusParamsInfo(8200)
	s.AssertNoErr(t, err)
	fmt.Println(PretiJson(b))
}
