package nsbcli

import (
	"encoding/json"
	transactiontype "github.com/HyperService-Consortium/NSB/application/transaction-type"

	appl "github.com/HyperService-Consortium/NSB/application"
	"github.com/HyperService-Consortium/NSB/grpc/nsbrpc"
	uiptypes "github.com/HyperService-Consortium/go-uip/uiptypes"
)

func (nc *NSBClient) AddMerkleProof(
	user uiptypes.Signer, toAddress []byte,
	merkletype uint16, rootHash, proof, key, value []byte,
) (*ResultInfo, error) {
	// fmt.Println(string(buf.Bytes()))
	fap, err := nc.addMerkleProof(merkletype, rootHash, proof, key, value)
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

func (nc *NSBClient) addMerkleProof(
	merkletype uint16, rootHash []byte, proof []byte, key []byte, value []byte,
) (*nsbrpc.FAPair, error) {
	var args appl.ArgsValidateMerkleProof

	args.Type = merkletype
	args.RootHash = rootHash

	args.Proof = proof
	args.Key = key
	args.Value = value
	b, err := json.Marshal(args)
	if err != nil {
		return nil, err
	}

	var fap = new(nsbrpc.FAPair)
	fap.FuncName = "system.merkleproof@validateMerkleProof"
	fap.Args = b
	// fmt.Println(PretiStruct(args), b)
	return fap, err
}
