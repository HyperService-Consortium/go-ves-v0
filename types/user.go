package types

type User interface {
	// User must has name or other guid

	// user is a kv-object
	KVObject
}

type ChainType = uint32
type Address = []byte
type UserName = string

type Account interface {
	GetChainType() ChainType
	GetAddress() Address
}

type UserBase interface {
	// insert accounts maps from guid to account
	InsertAccount(MultiIndex, UserName, Account) error

	// find accounts which guid is corresponding to user
	FindUser(MultiIndex, UserName) (User, error)

	// find accounts which guid is corresponding to user
	FindAccounts(MultiIndex, UserName, ChainType) ([]Account, error)

	// return true if user has this account
	HasAccount(MultiIndex, UserName, Account) (bool, error)

	// return the user which has this account
	InvertFind(MultiIndex, Account) (UserName, error)
}
