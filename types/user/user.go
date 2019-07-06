package user

import types "github.com/Myriad-Dreamin/go-ves/types"

// an implementation of types.Account is uiprpc.Account from "github.com/Myriad-Dreamin/go-ves/grpc"

type User struct {
	Name     string
	Accounts []*AccountAdapdator
}

type AccountAdapdator struct {
	ChainType uint64
	Address   []byte
}

func NewAccountAdapdator(account types.Account) *AccountAdapdator {
	return &AccountAdapdator{
		ChainType: account.GetChainType(),
		Address:   account.GetAddress(),
	}
}

func NewAccountsAdapdator(accounts []types.Account) (accs []*AccountAdapdator) {
	for _, account := range accounts {
		accs = append(accs, NewAccountAdapdator(account))
	}
	return
}

type UserBase struct {
}

func (ub *UserBase) InsertAccount(db types.MultiIndex, name string, account types.Account) error {
	a := NewAccountAdapdator(account)
}

func (ub *UserBase) FindUser(db types.MultiIndex, name string) (user User, err error) {

}

func (ub *UserBase) FindAccounts(db types.MultiIndex, username string, chainType uint32) (accs []Account, err error) {

}

func (ub *UserBase) HasAccount(db types.MultiIndex, name string, account types.Account) (has bool, err error) {

}

func (ub *UserBase) InvertFind(db types.MultiIndex, account types.Account) (name string, err error) {

}
