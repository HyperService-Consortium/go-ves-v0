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

func TestGetConsensusState(t *testing.T) {
	b, err := NewNSBClient(host).GetConsensusState()
	s.AssertNoErr(t, err)
	fmt.Println(PretiJson(b))
}

func TestGetGenesis(t *testing.T) {
	b, err := NewNSBClient(host).GetGenesis()
	s.AssertNoErr(t, err)
	fmt.Println(PretiJson(b))
}
func TestGetHealth(t *testing.T) {
	b, err := NewNSBClient(host).GetHealth()
	s.AssertNoErr(t, err)
	fmt.Println(PretiJson(b))
}
func TestGetNetInfo(t *testing.T) {
	b, err := NewNSBClient(host).GetNetInfo()
	s.AssertNoErr(t, err)
	fmt.Println(PretiJson(b))
}
func TestGetNumUnconfirmedTxs(t *testing.T) {
	b, err := NewNSBClient(host).GetNumUnconfirmedTxs()
	s.AssertNoErr(t, err)
	fmt.Println(PretiJson(b))
}
func TestGetStatus(t *testing.T) {
	b, err := NewNSBClient(host).GetStatus()
	s.AssertNoErr(t, err)
	fmt.Println(PretiJson(b))
}
func TestGetUnconfirmedTxs(t *testing.T) {
	b, err := NewNSBClient(host).GetUnconfirmedTxs(1)
	s.AssertNoErr(t, err)
	fmt.Println(PretiJson(b))
}
func TestGetValidators(t *testing.T) {
	b, err := NewNSBClient(host).GetUnconfirmedTxs(8200)
	s.AssertNoErr(t, err)
	fmt.Println(PretiJson(b))
}
