package nsbcli

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"encoding/json"
	"io"

	appl "github.com/HyperService-Consortium/NSB/application"
	cmn "github.com/HyperService-Consortium/NSB/common"
	nsbmath "github.com/HyperService-Consortium/NSB/math"
	uiptypes "github.com/HyperService-Consortium/go-uip/types"
)

func (nc *NSBClient) InsuranceClaim(
	user uiptypes.Signer, contractAddress []byte,
	tid, aid uint64,
) (*DeliverTx, error) {
	var txHeader cmn.TransactionHeader
	var buf = bytes.NewBuffer(make([]byte, mxBytes))
	buf.Reset()
	// fmt.Println(string(buf.Bytes()))
	err := nc.insuranceClaim(buf, tid, aid)
	if err != nil {
		return nil, err
	}

	var fap appl.FAPair
	fap.FuncName = "InsuranceClaim"
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

func (nc *NSBClient) insuranceClaim(
	w io.Writer,
	tid, aid uint64,
) error {
	// var args ISC.ArgsInsuranceClaim
	// args.Tid = tid
	// args.Aid = aid
	// b, err := json.Marshal(args)
	err := binary.Write(w, binary.BigEndian, tid)
	if err != nil {
		return err
	}
	err = binary.Write(w, binary.BigEndian, aid)

	// .Println(PretiStruct(args), b)
	// _, err = w.Write(b)
	return err
}
