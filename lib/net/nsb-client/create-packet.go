package nsbcli

import (
	"bytes"
	"github.com/HyperService-Consortium/NSB/grpc/nsbrpc"
	uiptypes "github.com/HyperService-Consortium/go-uip/uiptypes"
	"github.com/golang/protobuf/proto"
	"math/rand"
	"time"
)

func (nc *NSBClient) CreateContractPacket(
	s uiptypes.Signer, toAddress, value []byte, pair *nsbrpc.FAPair,
) (*nsbrpc.TransactionHeader, error) {
	data, err := proto.Marshal(pair)
	if err != nil {
		return nil, err
	}
	return nc.CreateNormalPacket(s, toAddress, data, value)
}


func (nc *NSBClient) CreateNormalPacket(
	s uiptypes.Signer, toAddress, data, value []byte,
) (*nsbrpc.TransactionHeader, error) {
	txHeader := new(nsbrpc.TransactionHeader)
	var err error

	txHeader.Data = data
	txHeader.Dst = toAddress
	txHeader.Src = s.GetPublicKey()

	nonce := make([]byte, 32)
	_, err = rand.Read(nonce)
	if err != nil {
		return nil, err
	}
	txHeader.Nonce = nonce
	txHeader.Value = value
	// bug: buf.Reset()
	buf := bytes.NewBuffer(make([]byte, mxBytes))

	buf.Write(txHeader.Src)
	buf.Write(txHeader.Dst)
	buf.Write(txHeader.Data)
	buf.Write(txHeader.Value)
	buf.Write(txHeader.Nonce)
	txHeader.Signature = s.Sign(buf.Bytes()).Bytes()
	return txHeader, nil
}


func init() {
	rand.Seed(time.Now().UnixNano())
}
