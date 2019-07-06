package user

import (
	"errors"

	types "github.com/Myriad-Dreamin/go-ves/types"
)

type XORMUserAdapter struct {
	ID        int64  `xorm:"pk autoincr 'id'"`
	Name      string `xorm:"'name'"`
	ChainType uint64 `xorm:"'chain_type'"`
	Address   []byte `xorm:"'address'"`
}

func NewXORMUserAdapter(name string, account types.Account) *XORMUserAdapter {
	return &XORMUserAdapter{
		Name:      name,
		ChainType: account.GetChainType(),
		Address:   account.GetAddress(),
	}
}

func NewUserAdapdators(name string, accounts []types.Account) (accs []*XORMUserAdapter) {
	for _, account := range accounts {
		accs = append(accs, NewXORMUserAdapter(name, account))
	}
	return
}

func XORMUserAdapterToAccounts(accounts []XORMUserAdapter) (accs []types.Account) {
	for _, account := range accounts {
		accs = append(accs, account)
	}
	return
}

func UserFromAdapdator(accounts []XORMUserAdapter) (user *User) {
	if accounts == nil {
		return nil
	}
	return &User{
		Name:     accounts[0].Name,
		Accounts: XORMUserAdapterToAccounts(accounts),
	}
}

func (ua XORMUserAdapter) GetAddress() []byte {
	return ua.Address
}

func (ua XORMUserAdapter) GetChainType() uint64 {
	return ua.ChainType
}

func (ua XORMUserAdapter) GetID() int64 {
	return ua.ID
}

func (ua XORMUserAdapter) GetSlicePtr() interface{} {
	return new([]XORMUserAdapter)
}

func (ua XORMUserAdapter) GetObjectPtr() interface{} {
	return new(XORMUserAdapter)
}

func (ua XORMUserAdapter) ToKVMap() map[string]interface{} {
	return map[string]interface{}{
		"id":         ua.ID,
		"name":       ua.Name,
		"chain_type": ua.ChainType,
		"address":    ua.Address,
	}
}

type UserBase struct {
}

func (ub UserBase) InsertAccount(db types.MultiIndex, name string, account types.Account) error {
	return db.Insert(NewXORMUserAdapter(name, account))
}

func (ub UserBase) FindUser(db types.MultiIndex, name string) (user types.User, err error) {
	condition := XORMUserAdapter{Name: name}
	sli, err := db.Select(&condition)
	if err != nil {
		return
	}
	if sli == nil {
		return nil, errors.New("not found")
	}
	return UserFromAdapdator(sli.([]XORMUserAdapter)), nil
}

func (ub UserBase) FindAccounts(db types.MultiIndex, username string, chainType uint64) (accs []types.Account, err error) {
	condition := XORMUserAdapter{Name: username, ChainType: chainType}
	sli, err := db.Select(&condition)
	if err != nil {
		return
	}
	if sli == nil {
		return nil, errors.New("not found")
	}
	return XORMUserAdapterToAccounts(sli.([]XORMUserAdapter)), nil
}

func (ub UserBase) HasAccount(db types.MultiIndex, name string, account types.Account) (has bool, err error) {
	return db.Get(NewXORMUserAdapter(name, account))
}

func (ub UserBase) InvertFind(db types.MultiIndex, account types.Account) (name string, err error) {

	//condition := User{Accounts}
	return "", nil
}
