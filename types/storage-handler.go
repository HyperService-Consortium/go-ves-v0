package types

import "github.com/HyperService-Consortium/go-uip/uiptypes"

type StorageHandler interface {
	GetTransactionProof(index Index, chainID uiptypes.ChainID, blockID uiptypes.BlockID, color []byte) (uiptypes.MerkleProof, error)
	GetStorageAt(index Index, chainID uiptypes.ChainID, typeID uiptypes.TypeID, contractAddress uiptypes.ContractAddress, pos []byte, description []byte) (uiptypes.Variable, error)
	SetStorageOf(index Index, chainID uiptypes.ChainID, typeID uiptypes.TypeID, contractAddress uiptypes.ContractAddress, pos []byte, description []byte, variable uiptypes.Variable) error
}

//type SessionKVBase interface {
//	SetKV(Index, isc_address, provedKey, provedValue) error
//	GetKV(Index, isc_address, provedKey) (provedValue, error)
//	GetSetter(Index, isc_address) KVSetter
//	GetGetter(Index, isc_address) KVGetter
//}
//
//type provedKey = []byte
//type provedValue = []byte
//
//type KVSetter interface {
//	Set(provedKey, provedValue) error
//}
//
//type KVGetter interface {
//	Get(provedKey) (provedValue, error)
//}
