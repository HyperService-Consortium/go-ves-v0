package nsbcli

import (
	"bytes"
	transactiontype "github.com/HyperService-Consortium/NSB/application/transaction-type"
	"github.com/HyperService-Consortium/NSB/grpc/nsbrpc"
	"github.com/HyperService-Consortium/go-uip/uiptypes"
	"github.com/gogo/protobuf/proto"
	"math/rand"
	"time"
)

func (*NSBClient) Sign(s uiptypes.Signer, txHeader *nsbrpc.TransactionHeader) *nsbrpc.TransactionHeader {
	// bug: buf.Reset()
	buf := bytes.NewBuffer(make([]byte, mxBytes))
	buf.Write(txHeader.Src)
	buf.Write(txHeader.Dst)
	buf.Write(txHeader.Data)
	buf.Write(txHeader.Value)
	buf.Write(txHeader.Nonce)

	txHeader.Signature = s.Sign(buf.Bytes()).Bytes()
	return txHeader
}

func (nc *NSBClient) CreateContractPacket(
	s uiptypes.Signer, toAddress, value []byte, pair *nsbrpc.FAPair,
) (*nsbrpc.TransactionHeader, error) {
	data, err := proto.Marshal(pair)
	if err != nil {
		return nil, err
	}
	return nc.CreateNormalPacket(s, toAddress, data, value)
}

func (nc *NSBClient) CreateUnsignedContractPacket(
	srcAddress, dstAddress, value []byte, pair *nsbrpc.FAPair,
) (*nsbrpc.TransactionHeader, error) {
	data, err := proto.Marshal(pair)
	if err != nil {
		return nil, err
	}
	return nc.CreateUnsignedNormalPacket(srcAddress, dstAddress, data, value)
}

func (nc *NSBClient) CreateUnsignedNormalPacket(
	srcAddress, dstAddress, data, value []byte,
) (*nsbrpc.TransactionHeader, error) {
	txHeader := new(nsbrpc.TransactionHeader)
	var err error

	txHeader.Data = data
	txHeader.Dst = dstAddress
	txHeader.Src = srcAddress

	nonce := make([]byte, 32)
	_, err = rand.Read(nonce)
	if err != nil {
		return nil, err
	}
	txHeader.Nonce = nonce
	txHeader.Value = value
	return txHeader, nil
}

func (nc *NSBClient) CreateNormalPacket(
	s uiptypes.Signer, toAddress, data, value []byte,
) (*nsbrpc.TransactionHeader, error) {
	txHeader, err := nc.CreateUnsignedNormalPacket(s.GetPublicKey(), toAddress, data, value)
	if err != nil {
		return nil, err
	}
	nc.Sign(s, txHeader)
	return txHeader, nil
}

func (nc *NSBClient) sendTx(t transactiontype.Type, txHeader *nsbrpc.TransactionHeader, err error) (*ResultInfo, error) {
	if err != nil {
		return nil, err
	}
	ret, err := nc.sendContractTx(t, txHeader)
	if err != nil {
		return nil, err
	}
	// fmt.Println(PretiStruct(ret), err)
	return ret, nil
}

func (nc *NSBClient) systemCall(txHeader *nsbrpc.TransactionHeader, err error) (*ResultInfo, error) {
	return nc.sendTx(transactiontype.SystemCall, txHeader, err)
}

func (nc *NSBClient) createContract(txHeader *nsbrpc.TransactionHeader, err error) (*ResultInfo, error) {
	return nc.sendTx(transactiontype.CreateContract, txHeader, err)
}

func (nc *NSBClient) sendTransaction(txHeader *nsbrpc.TransactionHeader, err error) (*ResultInfo, error) {
	return nc.sendTx(transactiontype.SendTransaction, txHeader, err)
}

func (nc *NSBClient) sign(user uiptypes.Signer, txHeader *nsbrpc.TransactionHeader, err error) (*nsbrpc.TransactionHeader, error) {
	if err != nil {
		return nil, err
	}
	nc.Sign(user, txHeader)
	return txHeader, nil
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
