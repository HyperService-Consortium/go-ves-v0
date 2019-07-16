package types

import uiptypes "github.com/Myriad-Dreamin/go-uip/types"

type user_name = string

type User interface {
	// user must has name or other guid

	// user is a kv-object
	KVObject

	GetName() user_name
	GetAccounts() []uiptypes.Account
}

type has = bool

// the database which used by others
type UserBase interface {
	// insert accounts maps from guid to account
	InsertAccount(MultiIndex, user_name, uiptypes.Account) error

	// DeleteAccount(MultiIndex, user_name, Account) error

	// DeleteAccountByName(MultiIndex, user_name) error

	// DeleteAccountByID(MultiIndex, user_id) error

	// find accounts which guid is corresponding to user
	FindUser(MultiIndex, user_name) (User, error)

	// find accounts which guid is corresponding to user
	FindAccounts(MultiIndex, user_name, uiptypes.ChainId) ([]uiptypes.Account, error)

	// return true if user has this account
	HasAccount(MultiIndex, user_name, uiptypes.Account) (has, error)

	// return the user which has this account
	InvertFind(MultiIndex, uiptypes.Account) (user_name, error)
}
