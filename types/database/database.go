package vesdb

import "github.com/Myriad-Dreamin/go-ves/types"

type Database struct {
	pdb    types.MultiIndex
	sesdb  types.SessionBase
	userdb types.UserBase
}

func (db *Database) SetMultiIndex(phyDB types.MultiIndex) bool {
	db.pdb = phyDB
	return true
}

func (db *Database) SetSessionBase(logicDB types.SessionBase) bool {
	db.sesdb = logicDB
	return true
}

func (db *Database) SetUserBase(logicDB types.UserBase) bool {
	db.userdb = logicDB
	return true
}

func (db *Database) InsertSession(session types.Session) error {
	return db.sesdb.InsertSession(db.pdb, session)
}

func (db *Database) FindSession(isc_address []byte) (types.Session, error) {
	return db.sesdb.FindSession(db.pdb, isc_address)
}

func (db *Database) UpdateSession(session types.Session) error {
	return db.sesdb.UpdateSession(db.pdb, session)
}

func (db *Database) DeleteSession(isc_address []byte) error {
	return db.sesdb.DeleteSession(db.pdb, isc_address)
}

func (db *Database) InsertAccount(user_name string, account types.Account) error {
	return db.userdb.InsertAccount(db.pdb, user_name, account)
}

// func (db *Database) DeleteAccount(user_name string, account types.Account) ( error) {

// }

// func (db *Database) DeleteAccountByName(user_name string) ( error) {

// }

// func (db *Database) DeleteAccountByID(user_id) ( error) {

// }

func (db *Database) FindUser(user_name string) (types.User, error) {
	return db.userdb.FindUser(db.pdb, user_name)
}

func (db *Database) FindAccounts(user_name string, chain_type uint32) ([]types.Account, error) {
	return db.userdb.FindAccounts(db.pdb, user_name, chain_type)
}

func (db *Database) HasAccount(user_name string, account types.Account) (bool, error) {
	return db.userdb.HasAccount(db.pdb, user_name, account)
}

func (db *Database) InvertFind(account types.Account) (string, error) {
	return db.userdb.InvertFind(db.pdb, account)
}
