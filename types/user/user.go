package user

import (
	types "github.com/Myriad-Dreamin/go-ves/types"
)

// an implementation of types.Account is uiprpc.Account from "github.com/Myriad-Dreamin/go-ves/grpc"

type User struct {
	ID       int64           `xorm:"pk autoincr 'id'"`
	Name     string          `xorm:"'name'"`
	Accounts []types.Account `xorm:"-"`
}

func ConvertAccounts(accounts []types.Account) (ret []map[string]interface{}) {
	for _, account := range accounts {
		ret = append(ret, map[string]interface{}{
			"chain_id": account.GetChainId(),
			"address":  account.GetAddress(),
		})
	}
	return
}

func (u User) ToKVMap() map[string]interface{} {
	return map[string]interface{}{
		"name":     u.Name,
		"accounts": ConvertAccounts(u.Accounts),
	}
}

func (u User) GetName() string {
	return u.Name
}

func (u User) GetAccounts() []types.Account {
	return u.Accounts
}

func (this User) GetID() int64 {
	return this.ID
}

func (this User) GetSlicePtr() interface{} {
	return new([]User)
}

func (this User) GetObjectPtr() interface{} {
	return &User{}
}
