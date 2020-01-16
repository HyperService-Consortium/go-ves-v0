package nsbcli

import (
	transactiontype "github.com/HyperService-Consortium/NSB/application/transaction-type"

	"github.com/HyperService-Consortium/NSB/grpc/nsbrpc"
	uiptypes "github.com/HyperService-Consortium/go-uip/uiptypes"
)

func (nc *NSBClient) SettleContract(
	user uiptypes.Signer, contractAddress []byte,
) (*DeliverTx, error) {
	// fmt.Println(string(buf.Bytes()))
	fap, err := nc.settleContract()
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

func (nc *NSBClient) settleContract() (*nsbrpc.FAPair, error) {
	var fap = new(nsbrpc.FAPair)
	fap.FuncName = "SettleContract"
	// fmt.Println(PretiStruct(args), b)
	return fap, nil
}
