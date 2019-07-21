package user

import (
	"errors"

	uiptypes "github.com/Myriad-Dreamin/go-uip/types"
	types "github.com/Myriad-Dreamin/go-ves/types"
)

type XORMUserAdapter struct {
	ID      int64  `xorm:"pk autoincr 'id'"`
	Name    string `xorm:"'name'"`
	ChainId uint64 `xorm:"'chain_id'"`
	Address []byte `xorm:"'address'"`
}

func NewXORMUserAdapter(name string, account uiptypes.Account) *XORMUserAdapter {
	return &XORMUserAdapter{
		Name:    name,
		ChainId: account.GetChainId(),
		Address: account.GetAddress(),
	}
}

func NewXORMUserAdapterWithOnlyAccount(account uiptypes.Account) *XORMUserAdapter {
	return &XORMUserAdapter{
		ChainId: account.GetChainId(),
		Address: account.GetAddress(),
	}
}

func NewUserAdapdators(name string, accounts []uiptypes.Account) (accs []*XORMUserAdapter) {
	for _, account := range accounts {
		accs = append(accs, NewXORMUserAdapter(name, account))
	}
	return
}

func XORMUserAdapterToAccounts(accounts []XORMUserAdapter) (accs []uiptypes.Account) {
	for _, account := range accounts {
		accs = append(accs, account)
	}
	return
}

func UserFromAdapdator(name string, accounts []XORMUserAdapter) (user *User) {
	if accounts == nil {
		return &User{
			Name: name,
		}
	}
	return &User{
		Name:     name,
		Accounts: XORMUserAdapterToAccounts(accounts),
	}
}

func (ua XORMUserAdapter) TableName() string {
	return "users"
}

func (ua XORMUserAdapter) GetAddress() []byte {
	return ua.Address
}

func (ua XORMUserAdapter) GetChainId() uint64 {
	return ua.ChainId
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
		"id":       ua.ID,
		"name":     ua.Name,
		"chain_id": ua.ChainId,
		"address":  ua.Address,
	}
}

type XORMUserBase struct {
}

func (ub XORMUserBase) InsertAccount(db types.MultiIndex, name string, account uiptypes.Account) error {
	return db.Insert(NewXORMUserAdapter(name, account))
}

func (ub XORMUserBase) FindUser(db types.MultiIndex, name string) (user types.User, err error) {
	condition := XORMUserAdapter{Name: name}
	sli, err := db.Select(&condition)
	if err != nil {
		return
	}
	if sli == nil {
		return nil, errors.New("not found")
	}
	return UserFromAdapdator(name, sli.([]XORMUserAdapter)), nil
}

func (ub XORMUserBase) FindAccounts(db types.MultiIndex, username string, chainType uint64) (accs []uiptypes.Account, err error) {
	condition := XORMUserAdapter{Name: username, ChainId: chainType}
	sli, err := db.Select(&condition)
	if err != nil {
		return
	}
	if sli == nil {
		return nil, errors.New("not found")
	}
	return XORMUserAdapterToAccounts(sli.([]XORMUserAdapter)), nil
}

func (ub XORMUserBase) HasAccount(db types.MultiIndex, name string, account uiptypes.Account) (has bool, err error) {
	return db.Get(NewXORMUserAdapter(name, account))
}

func (ub XORMUserBase) InvertFind(db types.MultiIndex, account uiptypes.Account) (name string, err error) {
	condition := NewXORMUserAdapterWithOnlyAccount(account)
	sli, err := db.Select(condition)
	if err != nil {
		return
	}
	if sli == nil {
		return "", errors.New("not found")
	}
	return (sli.([]XORMUserAdapter))[0].Name, nil
}
