package nsbcli

import (
	"bytes"
	"crypto/rand"
	"encoding/json"

	appl "github.com/HyperService-Consortium/NSB/application"
	cmn "github.com/HyperService-Consortium/NSB/common"
	nsbmath "github.com/HyperService-Consortium/NSB/math"
	uiptypes "github.com/HyperService-Consortium/go-uip/types"
)

func (nc *NSBClient) SettleContract(
	user uiptypes.Signer, contractAddress []byte,
) (*DeliverTx, error) {
	var txHeader cmn.TransactionHeader
	var err error
	var fap appl.FAPair
	fap.FuncName = "SettleContract"
	txHeader.Data, err = json.Marshal(fap)
	if err != nil {
		return nil, err
	}
	txHeader.ContractAddress = contractAddress
	txHeader.From = user.GetPublicKey()

	nonce := make([]byte, 32)
	_, err = rand.Read(nonce)
	if err != nil {
		return nil, err
	}
	txHeader.Nonce = nsbmath.NewUint256FromBytes(nonce)
	txHeader.Value = nsbmath.NewUint256FromBytes([]byte{0})
	// bug: buf.Reset()
	buf := bytes.NewBuffer(make([]byte, mxBytes))

	buf.Write(txHeader.From)
	buf.Write(txHeader.ContractAddress)
	buf.Write(txHeader.Data)
	buf.Write(txHeader.Value.Bytes())
	buf.Write(txHeader.Nonce.Bytes())
	txHeader.Signature = user.Sign(buf.Bytes()).Bytes()
	ret, err := nc.sendContractTx([]byte("sendTransaction"), []byte("isc"), &txHeader)
	// fmt.Println(PretiStruct(ret), err)
	if err != nil {
		return nil, err
	}
	return &ret.DeliverTx, nil
}
