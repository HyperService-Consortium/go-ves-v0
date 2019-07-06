package vesdb

import "github.com/Myriad-Dreamin/go-ves/types"

type Database struct {
}

func (db *Database) SetMultiIndex(phyDB types.MultiIndex) (success_or_not bool) {

}

func (db *Database) SetSessionBase(logicDB types.SessionBase) (success_or_not bool) {

}

func (db *Database) SetUserBase(logicDB types.UserBase) (success_or_not bool) {

}

func (db *Database) InsertSession(session types.Session) (err error) {

}

func (db *Database) FindSession(isc_address []byte) (session types.Session, err error) {

}

func (db *Database) UpdateSession(session types.Session) (err error) {

}

func (db *Database) DeleteSession(isc_address []byte) (err error) {

}

func (db *Database) InsertAccount(user_name string, account types.Account) (err error) {

}

// func (db *Database) DeleteAccount(user_name string, account types.Account) (err error) {

// }

// func (db *Database) DeleteAccountByName(user_name string) (err error) {

// }

// func (db *Database) DeleteAccountByID(user_id) (err error) {

// }

func (db *Database) FindUser(user_name string) (User, err error) {

}

func (db *Database) FindAccounts(user_name string, chain_type uint32) (accounts []types.Account, err error) {

}

func (db *Database) HasAccount(user_name string, account types.Account) (has bool, err error) {

}

func (db *Database) InvertFind(account types.Account) (user_name string, err error) {

}
