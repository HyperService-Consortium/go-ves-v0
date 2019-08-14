package nsbcli

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"io"

	appl "github.com/HyperServiceOne/NSB/application"
	cmn "github.com/HyperServiceOne/NSB/common"
	ISC "github.com/HyperServiceOne/NSB/contract/isc"
	nsbmath "github.com/HyperServiceOne/NSB/math"
	uiptypes "github.com/Myriad-Dreamin/go-uip/types"
)

func (nc *NSBClient) UserAck(
	user uiptypes.Signer, contractAddress []byte,
	address, signature []byte,
) (*DeliverTx, error) {
	var txHeader cmn.TransactionHeader
	var buf = bytes.NewBuffer(make([]byte, mxBytes))
	buf.Reset()
	// fmt.Println(string(buf.Bytes()))
	err := nc.userAck(buf, address, signature)
	if err != nil {
		return nil, err
	}

	var fap appl.FAPair
	fap.FuncName = "UserAck"
	fap.Args = buf.Bytes()
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
	buf = bytes.NewBuffer(make([]byte, mxBytes))

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

func (nc *NSBClient) userAck(
	w io.Writer,
	address, signature []byte,
) error {
	var args ISC.ArgsUserAck
	args.Address = address
	args.Signature = signature
	b, err := json.Marshal(args)
	if err != nil {
		return err
	}

	// fmt.Println(PretiStruct(args), b)
	_, err = w.Write(b)
	return err
}
