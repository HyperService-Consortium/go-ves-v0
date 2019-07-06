package types

type chain_type = uint64
type address = []byte
type user_name = string

// an implementation of types.Account is uiprpc.Account from "github.com/Myriad-Dreamin/go-ves/grpc"
type Account interface {
	GetChainType() chain_type
	GetAddress() address
}

type Signature interface {
	GetSignatureType() uint32
	GetContent() []byte
}

type User interface {
	// user must has name or other guid

	// user is a kv-object
	KVObject

	GetName() user_name
	GetAccounts() []Account
}

type has = bool

// the database which used by others
type UserBase interface {
	// insert accounts maps from guid to account
	InsertAccount(MultiIndex, user_name, Account) error

	// DeleteAccount(MultiIndex, user_name, Account) error

	// DeleteAccountByName(MultiIndex, user_name) error

	// DeleteAccountByID(MultiIndex, user_id) error

	// find accounts which guid is corresponding to user
	FindUser(MultiIndex, user_name) (User, error)

	// find accounts which guid is corresponding to user
	FindAccounts(MultiIndex, user_name, chain_type) ([]Account, error)

	// return true if user has this account
	HasAccount(MultiIndex, user_name, Account) (has, error)

	// return the user which has this account
	InvertFind(MultiIndex, Account) (user_name, error)
}
