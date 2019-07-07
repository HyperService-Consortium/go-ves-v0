package types

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

type VESDB interface {
	SetIndex(Index) success_or_not

	SetMultiIndex(MultiIndex) success_or_not

	SetSessionBase(SessionBase) success_or_not

	SetUserBase(UserBase) success_or_not

	// insert accounts maps from guid to account
	InsertSessionInfo(Session) error

	// find accounts which guid is corresponding to user
	FindSessionInfo(isc_address) (Session, error)

	UpdateSessionInfo(Session) error

	DeleteSessionInfo(isc_address) error

	// insert accounts maps from guid to account
	InsertAccount(user_name, Account) error

	// DeleteAccount(user_name, Account) error

	// DeleteAccountByName(user_name) error

	// DeleteAccountByID(user_id) error

	// find accounts which guid is corresponding to user
	FindUser(user_name) (User, error)

	// find accounts which guid is corresponding to user
	FindAccounts(user_name, chain_id) ([]Account, error)

	// return true if user has this account
	HasAccount(user_name, Account) (bool, error)

	// return the user which has this account
	InvertFind(Account) (user_name, error)
}
