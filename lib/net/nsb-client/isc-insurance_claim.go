package nsbcli

import (
	"encoding/binary"
	transactiontype "github.com/HyperService-Consortium/NSB/application/transaction-type"

	"github.com/HyperService-Consortium/NSB/grpc/nsbrpc"
	uiptypes "github.com/HyperService-Consortium/go-uip/uiptypes"
)

func (nc *NSBClient) InsuranceClaim(
	user uiptypes.Signer, contractAddress []byte,
	tid, aid uint64,
) (*DeliverTx, error) {
	// fmt.Println(string(buf.Bytes()))
	fap, err := nc.insuranceClaim(tid, aid)
	if err != nil {
		return nil, err
	}
	txHeader, err := nc.CreateContractPacket(user, contractAddress, []byte{0}, fap)
	if err != nil {
		return nil, err
	}
	ret, err := nc.sendContractTx(transactiontype.SendTransaction, txHeader)
	if err != nil {
		return nil, err
	}
	// fmt.Println(PretiStruct(ret), err)
	return &ret.DeliverTx, nil
}

func (nc *NSBClient) insuranceClaim(
	tid, aid uint64,
) (*nsbrpc.FAPair, error) {
	var fap = new(nsbrpc.FAPair)
	fap.FuncName = "InsuranceClaim"
	fap.Args = make([]byte, 16)
	binary.BigEndian.PutUint64(fap.Args[0:8], tid)
	binary.BigEndian.PutUint64(fap.Args[8:], aid)
	// fmt.Println(PretiStruct(args), b)
	return fap, nil
}
