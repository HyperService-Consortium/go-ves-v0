package nsbcli

import (
	"encoding/json"
	"fmt"
	"testing"

	appl "github.com/HyperServiceOne/NSB/application"
)

//TestAddAction is not public, just for test
func TestAddAction(t *testing.T) {
	ret, err := NewNSBClient(host).AddAction(&eacc{[]byte{1, 2, 3}, []byte("abcc")}, []byte("abcc"), []byte("abcc"), 0, 1, 1, []byte("12345678"), []byte("12345678123456781234567812345678123456781234567812345678123456781"))
	if err != nil {
		t.Error("USR ACK ERR", err, "\n")
	} else {
		fmt.Println("SSSSSSSSSSSSSSSSS", ret)
	}
}

//TestGetAction is not public, just for test
func TestGetAction(t *testing.T) {
	ret, err := NewNSBClient(host).GetAction(&eacc{[]byte{1, 2, 3}, []byte("abcc")}, []byte("abcc"), []byte("abcc"), 0, 1)
	if err != nil {
		t.Error("USR ACK ERR", err, "\n")
	} else {
		fmt.Println("SSSSSSSSSSSSSSSSS", ret)
	}
}

// TestAddMerkleProof is not public, just for test
func TestAddMerkleProof(t *testing.T) {
	var proof appl.MPTMerkleProof
	proof.RootHash = []byte("roothash")
	jp, err := json.Marshal(proof)
	if err != nil {
		t.Error(err)
		return
	}
	ret, err := NewNSBClient(host).AddMerkleProof(&eacc{[]byte{1, 2, 3}, []byte("abcc")}, []byte("abcc"), 2, []byte("roothash"), jp, []byte("key"), []byte("value"))
	if err != nil {
		t.Error("USR ACK ERR", err, "\n")
	} else {
		fmt.Println("SSSSSSSSSSSSSSSSS", ret)
	}
}

// TestAddBlockCheck is not public, just for test
func TestAddBlockCheck(t *testing.T) {
	const TransactionRoot = 1
	ret, err := NewNSBClient(host).AddBlockCheck(&eacc{[]byte{1, 2, 3}, []byte("abcc")}, []byte("abcc"), 2, []byte("3e0"), []byte("roothash"), TransactionRoot)
	if err != nil {
		t.Error("USR ACK ERR", err, "\n")
	} else {
		fmt.Println("SSSSSSSSSSSSSSSSS", ret)
	}
}

// TestGetMerkleProof is not public, just for test
func TestGetMerkleProof(t *testing.T) {
	const TransactionRoot = 1
	ret, err := NewNSBClient(host).GetMerkleProof(&eacc{[]byte{1, 2, 3}, []byte("abcc")}, []byte("abcc"), 2, []byte("roothash"), []byte("key"), 2, []byte("3e0"), TransactionRoot)
	if err != nil {
		t.Error("USR ACK ERR", err, "\n")
	} else {
		fmt.Println("SSSSSSSSSSSSSSSSS", ret)
	}
}
