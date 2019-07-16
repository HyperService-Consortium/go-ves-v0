package nsbcli

import (
	"encoding/hex"
	"fmt"
	"testing"

	tx "github.com/HyperServiceOne/NSB/contract/isc/transaction"
	nmath "github.com/HyperServiceOne/NSB/math"
)

type eacc struct {
	pbk []byte
	sig []byte
}

func (e *eacc) GetPublicKey() []byte {
	return e.pbk
}

func (e *eacc) Sign(b []byte) []byte { return e.sig }

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
	fmt.Println(NewNSBClient(host).CreateISC(&eacc{[]byte{1, 2, 3}, []byte("abcc")}, []uint32{0, 0}, [][]byte{[]byte{1, 2, 3}, []byte{1, 2, 4}}, []*tx.TransactionIntent{&opintent}, []byte{0}))
}
