package nsbcli

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"testing"

	tx "github.com/HyperServiceOne/NSB/contract/isc/transaction"
	nmath "github.com/HyperServiceOne/NSB/math"
	uiptypes "github.com/Myriad-Dreamin/go-uip/types"
)

type sigg struct {
	sig []byte
}

func (s *sigg) GetSignatureType() uint32 {
	return 0
}
func (s *sigg) GetContent() []byte {
	return s.sig
}

func (s *sigg) Bytes() []byte {
	return s.sig
}
func (s *sigg) String() string {
	return hex.EncodeToString(s.sig)
}
func (s *sigg) FromBytes([]byte) bool {
	return true
}
func (s *sigg) FromString(string) bool {
	return true
}

func (s *sigg) Equal(uiptypes.HexType) bool {
	return true
}
func (s *sigg) IsValid() bool {
	return true
}

type eacc struct {
	pbk []byte
	sig []byte
}

func (e *eacc) GetPublicKey() []byte {
	return e.pbk
}

func (e *eacc) Sign(b []byte) uiptypes.Signature { return &sigg{e.sig} }

type obj map[string]interface{}

func TestCreateISC(t *testing.T) {
	var opintent = tx.TransactionIntent{
		Fr:  []byte{1, 2, 3},
		To:  []byte{1, 2, 3, 4},
		Seq: nmath.NewUint256FromHexString("00"),
		Amt: nmath.NewUint256FromHexString("02e0"),
	}
	q, _ := hex.DecodeString("01020301020373635f6f776e657273223a5b2241514944222c2241514945225d2c2272657175697265645f66756e6473223a5b302c305d2c227665735f7369676e6174757265223a2241413d3d222c227472616e73616374696f6e5f696e74656e7473223a5b7b2266726f6d223a2241514944222c22746f223a224151494442413d3d222c22736571223a22222c22616d74223a224175413d222c226d657461223a6e756c6c")
	fmt.Println(string(q))
	fmt.Println(NewNSBClient(host).CreateISC(&eacc{[]byte{1, 2, 3}, []byte("abcc")}, []uint32{0, 0}, [][]byte{[]byte{1, 2, 3}, []byte{1, 2, 4}}, [][]byte{opintent.Bytes()}, []byte{0}))
}
func TestUpdateTxInfo(t *testing.T) {
	var iscAddress, err = base64.StdEncoding.DecodeString("6HaW8O5JuraT9Ew5eW6etChPc2uyahVGhm9KjF+BBqc=")
	if err != nil {
		t.Error(err)
		return
	}
	ret, err := NewNSBClient(host).UpdateTxInfo(&eacc{[]byte{1, 2, 3}, []byte("abcc")}, iscAddress, 0, new(tx.TransactionIntent))
	if err != nil {
		t.Error("UPD ERR", err, "\n")
	} else {
		fmt.Println("SSSSSSSSSSSSSSSSS", ret)
	}
}
func TestFreezeInfo(t *testing.T) {
	ret, err := NewNSBClient(host).FreezeInfo(&eacc{[]byte{1, 2, 3}, []byte("abcc")}, []byte("6HaW8O5JuraT9Ew5eW6etChPc2uyahVGhm9KjF+BBqc="), 1)
	if err != nil {
		t.Error("FREEZE ERR", err, "\n")
	} else {
		fmt.Println("SSSSSSSSSSSSSSSSS", ret)
	}
}

func TestUserAck(t *testing.T) {
	ret, err := NewNSBClient(host).UserAck(&eacc{[]byte{1, 2, 3}, []byte("abcc")}, []byte("6HaW8O5JuraT9Ew5eW6etChPc2uyahVGhm9KjF+BBqc="), []byte("6HaW8O5JuraT9Ew5eW6etChPc2uyahVGhm9KjF+BBqc="), []byte("SB"))
	if err != nil {
		t.Error("USR ACK ERR", err, "\n")
	} else {
		fmt.Println("SSSSSSSSSSSSSSSSS", ret)
	}
}

func TestInsuranceClaim(t *testing.T) {
	_, err := NewNSBClient(host).InsuranceClaim(&eacc{[]byte{1, 2, 3}, []byte("abcc")}, []byte("6HaW8O5JuraT9Ew5eW6etChPc2uyahVGhm9KjF+BBqc="), 0, 1)
	if err != nil {
		t.Error("INSURANCE CLAIM ERR", err, "\n")
	} else {
		fmt.Println("SSSSSSSSSSSSSSSSS")
	}
}
func TestSettleContract(t *testing.T) {
	ret, err := NewNSBClient(host).SettleContract(&eacc{[]byte{1, 2, 3}, []byte("abcc")}, []byte("6HaW8O5JuraT9Ew5eW6etChPc2uyahVGhm9KjF+BBqc="))
	if err != nil {
		t.Error("SETTLE CONTRACT ERR", err, "\n")
	} else {
		fmt.Println("SSSSSSSSSSSSSSSSS", ret)
	}
}
