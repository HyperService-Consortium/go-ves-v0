package vesdb

import types "github.com/Myriad-Dreamin/go-ves/types"

type Database struct {
	sindb  types.Index
	muldb  types.MultiIndex
	sesdb  types.SessionBase
	userdb types.UserBase
}

func (db Database) SetIndex(phyDB types.Index) bool {
	db.sindb = phyDB
	return true
}

func (db Database) SetMultiIndex(phyDB types.MultiIndex) bool {
	db.muldb = phyDB
	return true
}

func (db Database) SetSessionBase(logicDB types.SessionBase) bool {
	db.sesdb = logicDB
	return true
}

func (db Database) SetUserBase(logicDB types.UserBase) bool {
	db.userdb = logicDB
	return true
}

func (db Database) InsertSessionInfo(session types.Session) error {
	return db.sesdb.InsertSessionInfo(db.muldb, session)
}

func (db Database) FindSessionInfo(isc_address []byte) (types.Session, error) {
	return db.sesdb.FindSessionInfo(db.muldb, isc_address)
}

func (db Database) UpdateSessionInfo(session types.Session) error {
	return db.sesdb.UpdateSessionInfo(db.muldb, session)
}

func (db Database) DeleteSessionInfo(isc_address []byte) error {
	return db.sesdb.DeleteSessionInfo(db.muldb, isc_address)
}

func (db Database) InsertAccount(user_name string, account types.Account) error {
	return db.userdb.InsertAccount(db.muldb, user_name, account)
}

// func (db Database) DeleteAccount(user_name string, account types.Account) ( error) {

// }

// func (db Database) DeleteAccountByName(user_name string) ( error) {

// }

// func (db Database) DeleteAccountByID(user_id) ( error) {

// }

func (db Database) FindUser(user_name string) (types.User, error) {
	return db.userdb.FindUser(db.muldb, user_name)
}

func (db Database) FindAccounts(user_name string, chain_type uint64) ([]types.Account, error) {
	return db.userdb.FindAccounts(db.muldb, user_name, chain_type)
}

func (db Database) HasAccount(user_name string, account types.Account) (bool, error) {
	return db.userdb.HasAccount(db.muldb, user_name, account)
}

func (db Database) InvertFind(account types.Account) (string, error) {
	return db.userdb.InvertFind(db.muldb, account)
}
