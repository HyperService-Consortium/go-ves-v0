package nsbcli

import (
	"encoding/json"
	transactiontype "github.com/HyperService-Consortium/NSB/application/transaction-type"

	appl "github.com/HyperService-Consortium/NSB/application"
	"github.com/HyperService-Consortium/NSB/grpc/nsbrpc"
	uiptypes "github.com/HyperService-Consortium/go-uip/uiptypes"
)

func (nc *NSBClient) AddBlockCheck(
	user uiptypes.Signer, toAddress []byte,
	chainID uint64, blockID, rootHash []byte, rcType uint8,
) (*ResultInfo, error) {
	// fmt.Println(string(buf.Bytes()))
	fap, err := nc.addBlockCheck(chainID, blockID, rootHash, rcType)
	if err != nil {
		return nil, err
	}
	txHeader, err := nc.CreateContractPacket(user, toAddress, []byte{0}, fap)
	if err != nil {
		return nil, err
	}
	ret, err := nc.sendContractTx(transactiontype.SystemCall, txHeader)
	if err != nil {
		return nil, err
	}
	// fmt.Println(PretiStruct(ret), err)
	return ret, nil
}

func (nc *NSBClient) addBlockCheck(
	chainID uint64, blockID, rootHash []byte, rtType uint8,
) (*nsbrpc.FAPair, error) {
	var args appl.ArgsAddBlockCheck
	args.ChainID = chainID
	args.BlockID = blockID
	args.RootHash = rootHash
	args.RtType = rtType
	b, err := json.Marshal(args)
	if err != nil {
		return nil, err
	}

	var fap = new(nsbrpc.FAPair)
	fap.FuncName = "system.merkleproof@addBlockCheck"
	fap.Args = b
	// fmt.Println(PretiStruct(args), b)
	return fap, err
}
