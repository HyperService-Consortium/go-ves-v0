package user

import (
	"errors"

	types "github.com/Myriad-Dreamin/go-ves/types"
)

// an implementation of types.Account is uiprpc.Account from "github.com/Myriad-Dreamin/go-ves/grpc"

type User struct {
	id       int64 `xorm:"pk autoincr"`
	Name     string
	Accounts []types.Account
}

func (this *User) GetId() int64 {
	return this.id
}
func (this *User) GetSlicePtr() interface{} {
	return new([]User)
}

func (this *User) GetObjectPtr() interface{} {
	return &User{}
}

type UserBase struct {
}

func (ub *UserBase) InsertAccount(db types.MultiIndex, name string, account types.Account) error {
	user := User{Name: name}
	user.Accounts = []types.Account{account}
	err := db.Insert(&user)
	return err
}

func (ub *UserBase) FindUser(db types.MultiIndex, name string) (user User, err error) {
	condition := User{Name: name}
	has, err := db.Get(&condition)
	if err != nil {
		return user, err
	}
	if has == false {
		return user, errors.New("Database does not have the object")
	}
	return condition, nil
}

func (ub *UserBase) FindAccounts(db types.MultiIndex, username string, chainType uint64) (accs []types.Account, err error) {
	user, err := ub.FindUser(db, username)
	if err != nil {
		return nil, err
	}
	var ret []types.Account
	for _, acc := range user.Accounts {
		if acc.GetChainType() == chainType {
			ret = append(ret, acc)
		}
	}
	return ret, nil
}

func (ub *UserBase) HasAccount(db types.MultiIndex, name string, account types.Account) (has bool, err error) {
	user, err := ub.FindUser(db, name)
	if err != nil {
		return false, err
	}
	for _, acc := range user.Accounts {
		if string(acc.GetAddress()) == string(account.GetAddress()) && acc.GetChainType() == account.GetChainType() {
			return true, nil
		}
	}
	return false, nil
}

func (ub *UserBase) InvertFind(db types.MultiIndex, account types.Account) (name string, err error) {

	//condition := User{Accounts}
	return "", nil
}
