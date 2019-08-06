package nsbcli

import (
	"fmt"
	"testing"
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

//TestAddMerkleProof is not public, just for test
// func TestAddMerkleProof(t *testing.T) {
// 	ret, err := NewNSBClient(host).AddMerkleProof(&eacc{[]byte{1, 2, 3}, []byte("abcc")}, []byte("abcc"), []byte("abcc"), 0, 1, 1, []byte("12345678"), []byte("12345678123456781234567812345678123456781234567812345678123456781"))
// 	if err != nil {
// 		t.Error("USR ACK ERR", err, "\n")
// 	} else {
// 		fmt.Println("SSSSSSSSSSSSSSSSS", ret)
// 	}
// }
//
//TestGetAction is not public, just for test
// func TestGetAction(t *testing.T) {
// 	ret, err := NewNSBClient(host).GetAction(&eacc{[]byte{1, 2, 3}, []byte("abcc")}, []byte("abcc"), []byte("abcc"), 0, 1)
// 	if err != nil {
// 		t.Error("USR ACK ERR", err, "\n")
// 	} else {
// 		fmt.Println("SSSSSSSSSSSSSSSSS", ret)
// 	}
// }
