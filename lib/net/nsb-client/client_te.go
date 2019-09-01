package nsbcli

import (
	"encoding/json"
	"fmt"
	"testing"

	mtest "github.com/HyperService-Consortium/mydrest"
)

var s mtest.TestHelper

const host = "47.251.2.73:26657"

//  PretiStruct helps read struct
func PretiStruct(minterface interface{}) string {
	je, _ := json.MarshalIndent(minterface, "", "\t")
	return string(je)
}

// TestGetAbciInfo is not public, just for test
func TestGetAbciInfo(t *testing.T) {
	b, err := NewNSBClient(host).GetAbciInfo()
	s.AssertNoErr(t, err)
	fmt.Println(PretiStruct(b))
}

// TestGetBlock is not public, just for test
func TestGetBlock(t *testing.T) {
	b, err := NewNSBClient(host).GetBlock(8200)
	s.AssertNoErr(t, err)
	fmt.Println(PretiStruct(b))
}

// TestGetCommitInfo is not public, just for test
func TestGetCommitInfo(t *testing.T) {
	b, err := NewNSBClient(host).GetCommitInfo(8200)
	s.AssertNoErr(t, err)
	fmt.Println(PretiStruct(b))
}

// TestGetBlocks is not public, just for test
func TestGetBlocks(t *testing.T) {
	b, err := NewNSBClient(host).GetBlocks(8199, 8200)
	s.AssertNoErr(t, err)
	fmt.Println(PretiStruct(b))
}

// TestGetBlockResults is not public, just for test
func TestGetBlockResults(t *testing.T) {
	b, err := NewNSBClient(host).GetBlockResults(8200)
	s.AssertNoErr(t, err)
	fmt.Println(PretiStruct(b))
}

// TestGetConsensusParamsInfo is not public, just for test
func TestGetConsensusParamsInfo(t *testing.T) {
	b, err := NewNSBClient(host).GetConsensusParamsInfo(8200)
	s.AssertNoErr(t, err)
	fmt.Println(PretiStruct(b))
}

// TestGetConsensusState is not public, just for test
func TestGetConsensusState(t *testing.T) {
	b, err := NewNSBClient(host).GetConsensusState()
	s.AssertNoErr(t, err)
	fmt.Println(PretiStruct(b))
}

// TestGetGenesis is not public, just for test
func TestGetGenesis(t *testing.T) {
	b, err := NewNSBClient(host).GetGenesis()
	s.AssertNoErr(t, err)
	fmt.Println(PretiStruct(b))
}

// TestGetHealth is not public, just for test
func TestGetHealth(t *testing.T) {
	b, err := NewNSBClient(host).GetHealth()
	s.AssertNoErr(t, err)
	fmt.Println(PretiStruct(b))
}

// TestGetNetInfo is not public, just for test
func TestGetNetInfo(t *testing.T) {
	b, err := NewNSBClient(host).GetNetInfo()
	s.AssertNoErr(t, err)
	fmt.Println(PretiStruct(b))
}

// TestGetNumUnconfirmedTxs is not public, just for test
func TestGetNumUnconfirmedTxs(t *testing.T) {
	b, err := NewNSBClient(host).GetNumUnconfirmedTxs()
	s.AssertNoErr(t, err)
	fmt.Println(PretiStruct(b))
}

// TestGetStatus is not public, just for test
func TestGetStatus(t *testing.T) {
	b, err := NewNSBClient(host).GetStatus()
	s.AssertNoErr(t, err)
	fmt.Println(PretiStruct(b))
}

// TestGetUnconfirmedTxs is not public, just for test
func TestGetUnconfirmedTxs(t *testing.T) {
	b, err := NewNSBClient(host).GetUnconfirmedTxs(1)
	s.AssertNoErr(t, err)
	fmt.Println(PretiStruct(b))
}

// TestGetValidators is not public, just for test
func TestGetValidators(t *testing.T) {
	b, err := NewNSBClient(host).GetUnconfirmedTxs(8200)
	s.AssertNoErr(t, err)
	fmt.Println(PretiStruct(b))
}
