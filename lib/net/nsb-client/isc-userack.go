package nsbcli

import (
	"encoding/json"
	transactiontype "github.com/HyperService-Consortium/NSB/application/transaction-type"

	ISC "github.com/HyperService-Consortium/NSB/contract/isc"
	"github.com/HyperService-Consortium/NSB/grpc/nsbrpc"
	uiptypes "github.com/HyperService-Consortium/go-uip/uiptypes"
)

func (nc *NSBClient) UserAck(
	user uiptypes.Signer, contractAddress []byte,
	address, signature []byte,
) (*DeliverTx, error) {
	// fmt.Println(string(buf.Bytes()))
	fap, err := nc.userAck(address, signature)
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
	return &ret.DeliverTx, nil
}

func (nc *NSBClient) userAck(
	address, signature []byte,
) (*nsbrpc.FAPair, error) {

	var args ISC.ArgsUserAck
	args.Address = address
	args.Signature = signature
	b, err := json.Marshal(args)
	if err != nil {
		return nil, err
	}

	var fap = new(nsbrpc.FAPair)
	fap.FuncName = "UserAck"
	fap.Args = b
	// fmt.Println(PretiStruct(args), b)
	return fap, err
}
