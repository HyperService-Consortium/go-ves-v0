package session

import (
	"bytes"
	"errors"
	"io"

	verifier "github.com/Myriad-Dreamin/go-ves/crypto/verifier"
	types "github.com/Myriad-Dreamin/go-ves/types"

	bitmap "github.com/Myriad-Dreamin/go-ves/bitmapping"
	const_prefix "github.com/Myriad-Dreamin/go-ves/database/const_prefix"
	serial_helper "github.com/Myriad-Dreamin/go-ves/serial_helper"
)

type SerialSession struct {
	ID               int64           `xorm:"pk unique notnull autoincr 'id'"`
	ISCAddress       []byte          `xorm:"unique 'isc_address'"`
	Accounts         []types.Account `xorm:"-"`
	Transactions     [][]byte        `xorm:"-"`
	TransactionCount uint32          `xorm:"'transaction_count'"`
	UnderTransacting uint32          `xorm:"'under_transacting'"`
	Status           uint8           `xorm:"'status'"`
	Content          []byte          `xorm:"'content'"`
	Acks             []byte          `xorm:"'acks'"`
}

func (ses SerialSession) TableName() string {
	return "ves_session"
}

func (ses SerialSession) ToKVMap() map[string]interface{} {
	return map[string]interface{}{
		"id":                ses.ID,
		"isc_address":       ses.ISCAddress,
		"transaction_count": ses.TransactionCount,
		"under_transacting": ses.UnderTransacting,
		"status":            ses.Status,
		"content":           ses.Content,
		"acks":              ses.Acks,
	}
}

func (ses SerialSession) GetID() int64 {
	return ses.ID
}

func (ses SerialSession) GetGUID() (isc_address []byte) {
	return ses.ISCAddress
}

func (ses SerialSession) GetObjectPtr() interface{} {
	return new(SerialSession)
}

func (ses SerialSession) GetSlicePtr() interface{} {
	return new([]SerialSession)
}

func (ses SerialSession) GetAccounts() []types.Account {

	// move to adapdator
	// if ses.Accounts == nil {
	// 	ses.Accounts = nil
	// }

	return ses.Accounts
}

func (ses SerialSession) GetTransaction(transaction_id uint32) []byte {
	// if ses.Transactions == nil {
	// 	ses.Transactions = make([][]byte, ses.TransactionCount)
	// }
	// if ses.Transactions[transaction_id] == nil {
	// 	ses.Transactions[transaction_id] = nil
	// }

	return ses.Transactions[transaction_id]
}

func (ses SerialSession) GetTransactions() (transactions [][]byte) {
	// if ses.Transactions == nil {
	// 	ses.Transactions = make([][]byte, ses.TransactionCount)
	// }

	return ses.Transactions
}

func (ses SerialSession) GetTransactingTransaction() (transaction_id uint32, err error) {
	// Status
	return ses.UnderTransacting, nil
}

func (ses SerialSession) GetContent() []byte {
	return ses.Content
}

func (ses SerialSession) InitFromOpIntents(types.OpIntents) (bool, string, error) {
	return false, "TODO", nil
}

func Verify(signature types.Signature, contentProviding, publicKey []byte) bool {
	return bytes.Equal(contentProviding, signature.GetContent()) &&
		verifier.Verify(signature, publicKey) == true
}

func (ses SerialSession) AckForInit(
	account types.Account,
	signature types.Signature,
) (success_or_not bool, help_info string, err error) {
	var addr = account.GetAddress()
	for idx, ak := range ses.GetAccounts() {
		if bytes.Equal(ak.GetAddress(), addr) {
			if !bitmap.InLength(ses.Acks, idx) {
				return false, "", errors.New("wrong Acks bytes set..")
			}
			if bitmap.Get(ses.Acks, idx) {
				return false, "have acked", nil
			}
			if !Verify(signature, ses.GetContent(), account.GetAddress()) {
				return false, "verify signature error...", nil
			}
			bitmap.Set(ses.Acks, idx)
			// todo: NSB
			return true, "", nil
		}
	}
	return false, "account not found in this session", nil
}

func (ses SerialSession) ProcessAttestation(
	attestation types.Attestation,
) (success_or_not bool, help_info string, err error) {
	return false, "TODO", nil
}

func (ses SerialSession) SyncFromISC() (err error) {
	return errors.New("TODO")
}

// the database which used by others
type SerialSessionBase struct {
}

func (sb SerialSessionBase) InsertSessionInfo(
	db types.MultiIndex, session types.Session,
) error {
	return db.Insert(session.(SerialSession))
}

func (sb SerialSessionBase) FindSessionInfo(
	db types.MultiIndex, isc_address []byte,
) (session types.Session, err error) {
	var sessions interface{}
	sessions, err = db.Select(&SerialSession{ISCAddress: isc_address})
	if err != nil {
		return
	}
	session = &(sessions.([]SerialSession)[0])
	return
}

func (sb SerialSessionBase) UpdateSessionInfo(
	db types.MultiIndex, session types.Session,
) (err error) {
	return db.Modify(session, session.ToKVMap())
}

func (sb SerialSessionBase) DeleteSessionInfo(
	db types.MultiIndex, isc_address []byte,
) (err error) {
	return db.Delete(&SerialSession{ISCAddress: isc_address})
}

func (sb SerialSessionBase) InsertSessionAccounts(
	db types.Index, isc_address []byte, accounts []types.Account,
) (err error) {
	var k, v []byte
	k, err = serial_helper.DecoratePrefix(const_prefix.AccountsPrefix, isc_address)
	if err != nil {
		return
	}
	v, err = serial_helper.SerializeAccountsInterface(accounts)
	if err != nil {
		return
	}
	db.Set(k, v)
	return
}

func (sb SerialSessionBase) FindSessionAccounts(
	db types.Index, isc_address []byte, getter func(uint64, []byte) error,
) (err error) {
	var k, v []byte
	k, err = serial_helper.DecoratePrefix(const_prefix.AccountsPrefix, isc_address)
	if err != nil {
		return
	}
	v, err = db.Get(k)
	if err != nil {
		return
	}
	var ct uint64
	var n int64
	for {
		n, ct, k, err = serial_helper.UnserializeAccountInterface(v)
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return
		}
		getter(ct, k)
		v = v[n:]
	}
}

// func (sb SerialSessionBase) InsertTransaction(
// 	db types.Index, transaction_id uint64, Transaction []byte
// ) (err error) (
//
// )
