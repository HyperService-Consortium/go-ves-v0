package types

import (
	"github.com/HyperService-Consortium/go-uip/uiptypes"
)

type KVPair struct {
	Key   string
	Value interface{}
}

type KVMap = map[string]interface{}

type KVObject interface {
	GetObjectPtr() interface{}
	GetSlicePtr() interface{}
	GetID() int64
	ToKVMap() KVMap
}

type Index interface {
	Get([]byte) ([]byte, error)
	Set([]byte, []byte) error
	Delete([]byte) error
	Batch([][]byte, [][]byte) error
}

type MultiIndex interface {
	// RegisterObject(KVObject) error

	Insert(KVObject) error

	Get(KVObject) (bool, error)

	Select(KVObject) (interface{}, error)

	SelectAll(KVObject) (interface{}, error)

	// 要求只Delete到一个
	Delete(KVObject) error

	// 可以Delete多个
	MultiDelete(KVObject) error

	Modify(KVObject, KVMap) error

	MultiModify(KVObject, KVMap) error
}

type KVPMultiIndex interface {
	Insert(...KVPair) error

	Select([]interface{}, ...KVPair) error

	// 要求只Delete到一个
	Delete(...KVPair) error
	// 可以Delete多个
	MultiDelete(...KVPair) error

	// 要求只Update到一个
	Modify([]KVPair, ...KVPair) error
	// 可以Update到多个
	MultiModify([]KVPair, ...KVPair) error
}


type ORMMultiIndex interface {
	MultiIndex
	// 要求只Update到一个
	// Modify(ORMObject, ORMObject) error
	// 可以Update到多个
	// MultiModify(ORMObject, ORMObject) error
}

type chain_id = uint64

type ChainInfo interface {
	GetChainType() uiptypes.ChainType
	GetChainHost() string
}

type StorageHandlerInterface interface {
	GetStorageAt(chainID uiptypes.ChainID, typeID uiptypes.TypeID, contractAddress uiptypes.ContractAddress, pos []byte, description []byte) (uiptypes.Variable, error)
	SetStorageOf(chainID uiptypes.ChainID, typeID uiptypes.TypeID, contractAddress uiptypes.ContractAddress, pos []byte, description []byte, variable uiptypes.Variable) error
}

type VESDB interface {
	SetIndex(Index) success_or_not
	SetMultiIndex(MultiIndex) success_or_not
	SetSessionBase(SessionBase) success_or_not
	SetUserBase(UserBase) success_or_not
	SetSessionKVBase(SessionKVBase) success_or_not
	SetStorageHandler(StorageHandler) success_or_not
	SetChainDNS(ChainDNS) success_or_not

	// insert accounts maps from guid to account
	InsertSessionInfo(Session) error

	// find accounts which guid is corresponding to user
	FindSessionInfo(isc_address) (Session, error)

	UpdateSessionInfo(Session) error

	DeleteSessionInfo(isc_address) error

	FindTransaction(isc_address, transaction_id, getter) error

	ActivateSession(isc_address)

	InactivateSession(isc_address)

	// insert accounts maps from guid to account
	InsertAccount(user_name, uiptypes.Account) error

	// DeleteAccount(user_name, Account) error

	// DeleteAccountByName(user_name) error

	// DeleteAccountByID(user_id) error

	// find accounts which guid is corresponding to user
	FindUser(user_name) (User, error)
	// find accounts which guid is corresponding to user
	FindAccounts(user_name, uint64) ([]uiptypes.Account, error)
	// return true if user has this account
	HasAccount(user_name, uiptypes.Account) (bool, error)
	// return the user which has this account
	InvertFind(uiptypes.Account) (user_name, error)

	SetKV(isc_address, provedKey, provedValue) error
	GetKV(isc_address, provedKey) (provedValue, error)

	GetSetter(isc_address) KVSetter
	GetGetter(isc_address) KVGetter

	StorageHandlerInterface
	ChainDNSInterface
}