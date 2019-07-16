package session

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"

	uiptypes "github.com/Myriad-Dreamin/go-uip/types"
	verifier "github.com/Myriad-Dreamin/go-ves/crypto/verifier"
	types "github.com/Myriad-Dreamin/go-ves/types"

	bitmap "github.com/Myriad-Dreamin/go-ves/bitmapping"
	const_prefix "github.com/Myriad-Dreamin/go-ves/database/const_prefix"
	serial_helper "github.com/Myriad-Dreamin/go-ves/serial_helper"

	opintents "github.com/Myriad-Dreamin/go-uip/op-intent"
)

type SerialSession struct {
	ID               int64              `json:"-" xorm:"pk unique notnull autoincr 'id'"`
	ISCAddress       []byte             `json:"-" xorm:"notnull 'isc_address'"`
	Accounts         []uiptypes.Account `json:"-" xorm:"-"`
	Transactions     [][]byte           `json:"transactions" xorm:"-"`
	TransactionCount uint32             `json:"-" xorm:"'transaction_count'"`
	UnderTransacting uint32             `json:"-" xorm:"'under_transacting'"`
	Status           uint8              `json:"-" xorm:"'status'"`
	Content          []byte             `json:"-" xorm:"'content'"`
	Acks             []byte             `json:"-" xorm:"'acks'"`
}

func randomSession() *SerialSession {
	var buf = make([]byte, 20)
	binary.PutVarint(buf, rand.Int63())
	return &SerialSession{
		ISCAddress: buf,
	}
}

func (ses *SerialSession) TableName() string {
	return "ves_session"
}

func (ses *SerialSession) ToKVMap() map[string]interface{} {
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

func (ses *SerialSession) GetID() int64 {
	return ses.ID
}

func (ses *SerialSession) GetGUID() (isc_address []byte) {
	return ses.ISCAddress
}

func (ses *SerialSession) GetObjectPtr() interface{} {
	return new(SerialSession)
}

func (ses *SerialSession) GetSlicePtr() interface{} {
	return new([]SerialSession)
}

func (ses *SerialSession) GetAccounts() []uiptypes.Account {

	// move to adapdator
	// if ses.Accounts == nil {
	// 	ses.Accounts = nil
	// }

	return ses.Accounts
}

func (ses *SerialSession) GetTransaction(transaction_id uint32) []byte {
	// if ses.Transactions == nil {
	// 	ses.Transactions = make([][]byte, ses.TransactionCount)
	// }
	// if ses.Transactions[transaction_id] == nil {
	// 	ses.Transactions[transaction_id] = nil
	// }

	return ses.Transactions[transaction_id]
}

func (ses *SerialSession) GetTransactions() (transactions [][]byte) {
	// if ses.Transactions == nil {
	// 	ses.Transactions = make([][]byte, ses.TransactionCount)
	// }

	return ses.Transactions
}

func (ses *SerialSession) GetTransactingTransaction() (transaction_id uint32, err error) {
	// Status
	return ses.UnderTransacting, nil
}

func (ses *SerialSession) GetContent() []byte {
	return ses.Content
}

func (ses *SerialSession) InitFromOpIntents(opIntents uiptypes.OpIntents) (bool, string, error) {
	intents, err := opintents.NewOpIntentInitializer().InitOpIntent(opIntents)
	if err != nil {
		return false, err.Error(), nil
	}
	fmt.Println(intents, err)
	ses.Transactions = make([][]byte, 0, len(intents))
	for _, intent := range intents {
		fmt.Println("insert", len(ses.Transactions))
		ses.Transactions = append(ses.Transactions, intent.Bytes())
		fmt.Println(string(intent.Bytes()), hex.EncodeToString(intent.Src), hex.EncodeToString(intent.Dst))
	}
	ses.TransactionCount = uint32(len(intents))
	ses.UnderTransacting = 0
	ses.Status = 0
	ses.Content, err = json.Marshal(ses)
	if err != nil {
		return false, "", err
	}

	// TransactionCount uint32          `xorm:"'transaction_count'"`
	// UnderTransacting uint32          `xorm:"'under_transacting'"`
	// Status           uint8           `xorm:"'status'"`
	// Content          []byte          `xorm:"'content'"`
	// Acks             []byte          `xorm:"'acks'"`

	return true, "", nil
}

func Verify(signature uiptypes.Signature, contentProviding, publicKey []byte) bool {
	return bytes.Equal(contentProviding, signature.GetContent()) &&
		verifier.Verify(signature, publicKey) == true
}

func (ses *SerialSession) AckForInit(
	account uiptypes.Account,
	signature uiptypes.Signature,
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

func (ses *SerialSession) ProcessAttestation(
	attestation uiptypes.Attestation,
) (success_or_not bool, help_info string, err error) {
	return false, "TODO", nil
}

func (ses *SerialSession) SyncFromISC() (err error) {
	return errors.New("TODO")
}

// the database which used by others
type SerialSessionBase struct {
}

func (sb SerialSessionBase) InsertSessionInfo(
	db types.MultiIndex, session types.Session,
) error {
	return db.Insert(session.(*SerialSession))
}

func (sb SerialSessionBase) FindSessionInfo(
	db types.MultiIndex, isc_address []byte,
) (session types.Session, err error) {
	var sessions interface{}
	sessions, err = db.Select(&SerialSession{ISCAddress: isc_address})
	if err != nil {
		return
	}
	f := sessions.([]SerialSession)
	if f == nil {
		return nil, errors.New("not found")
	}
	session = &f[0]
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
	db types.Index, isc_address []byte, accounts []uiptypes.Account,
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
