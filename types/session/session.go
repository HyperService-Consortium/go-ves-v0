package session

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"unsafe"

	uiptypes "github.com/Myriad-Dreamin/go-uip/types"
	account "github.com/Myriad-Dreamin/go-uip/types/account"
	types "github.com/Myriad-Dreamin/go-ves/types"

	bitmap "github.com/Myriad-Dreamin/go-ves/bitmapping"
	const_prefix "github.com/Myriad-Dreamin/go-ves/database/const_prefix"
	serial_helper "github.com/Myriad-Dreamin/go-ves/serial_helper"

	opintents "github.com/Myriad-Dreamin/go-uip/op-intent"
)

const (
	host = "47.251.2.73:26657"
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
	AckCount         uint32             `json:"-" xorm:"'ack_count'"`
	Signer           uiptypes.Signer    `json:"-" xorm:"-"`
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
		"ack_count":         ses.AckCount,
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

func (ses *SerialSession) GetAckCount() uint32 {
	return ses.AckCount
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

// must be used in single-thread env
type comparator struct {
	s      map[string]bool
	hacker [8]byte
}

func makeComparator() *comparator {
	return &comparator{s: make(map[string]bool)}
}

func (c *comparator) Insert(a uint64, b []byte) bool {
	h := md5.New()
	*(*uint64)(unsafe.Pointer(&c.hacker)) = a
	h.Write(c.hacker[:])
	h.Write(b)
	var nb = h.Sum(nil)
	if _, ok := c.s[string(nb)]; ok {
		return false
	} else {
		c.s[string(nb)] = true
		return true
	}
}

func (ses *SerialSession) SetSigner(signer uiptypes.Signer) {
	ses.Signer = signer
}

func (ses *SerialSession) InitFromOpIntents(opIntents uiptypes.OpIntents) (bool, string, error) {
	intents, err := opintents.NewOpIntentInitializer().InitOpIntent(opIntents)
	if err != nil {
		return false, err.Error(), nil
	}
	fmt.Println(intents, err)
	ses.Transactions = make([][]byte, 0, len(intents))

	ses.Accounts = nil
	c := makeComparator()
	for _, intent := range intents {
		fmt.Println("insert", len(ses.Transactions))
		ses.Transactions = append(ses.Transactions, intent.Bytes())
		fmt.Println(string(intent.Bytes()), hex.EncodeToString(intent.Src), hex.EncodeToString(intent.Dst))

		if c.Insert(intent.ChainId, intent.Src) {
			ses.Accounts = append(ses.Accounts, &account.PureAccount{ChainId: intent.ChainId, Address: intent.Src})
		}
		if c.Insert(intent.ChainId, intent.Dst) {
			ses.Accounts = append(ses.Accounts, &account.PureAccount{ChainId: intent.ChainId, Address: intent.Dst})
		}
	}
	ses.TransactionCount = uint32(len(intents))
	ses.UnderTransacting = 0
	ses.Status = 0
	ses.Acks = make([]byte, (len(ses.Accounts)+7)>>3)
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
	return true
	// return bytes.Equal(contentProviding, signature.GetContent()) &&
	// 	verifier.Verify(signature, publicKey) == true
}

func (ses *SerialSession) AckForInit(
	account uiptypes.Account,
	signature uiptypes.Signature,
) (success_or_not bool, help_info string, err error) {
	var addr = account.GetAddress()
	fmt.Println(ses.Acks, len(ses.Acks))
	for idx, ak := range ses.GetAccounts() {
		fmt.Println("comparing", hex.EncodeToString(ak.GetAddress()), hex.EncodeToString(addr))
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
			ses.AckCount++
			// if ses.AckCount == uint32(len(ses.Accounts)) {
			//
			// }
			// todo: NSB
			return true, "", nil
		}
	}
	return false, "account not found in this session", nil
}

type IterateAttestation struct {
	uiptypes.Attestation
	sigs []uiptypes.Signature
}

func (atte *IterateAttestation) GetSignatures() []uiptypes.Signature {
	return atte.sigs
}

func iter(atte uiptypes.Attestation, signer uiptypes.Signer) uiptypes.Attestation {
	return &IterateAttestation{atte, append(atte.GetSignatures(), &temSignature{
		sigtype: todo,
		content: signer.Sign(atte.GetSignatures()[len(atte.GetSignatures())-1].GetContent()),
	})}
}

const todo = 12333

type temSignature struct {
	content []byte
	sigtype uint32
}

func (t *temSignature) GetContent() []byte {
	return t.content
}

func (t *temSignature) GetSignatureType() uint32 {
	return t.sigtype
}

func isRawTransaction(tag uint8) bool {
	return (tag & 0x1) == 0x1
}

func (ses *SerialSession) ProcessAttestation(
	nsb types.NSBInterface, bn types.BNInterface, atte uiptypes.Attestation,
) (success_or_not bool, help_info string, err error) {
	// todo
	type Type = uint64

	const (
		Unknown Type = 0 + iota
		Initing
		Inited
		Instantiating
		Instantiated
		Open
		Opened
		Closed
	)

	tid, sigs := atte.GetTid(), atte.GetSignatures()

	if tid != uint64(ses.UnderTransacting) {
		return false, "this transaction is not undertransacting", nil
	}

	switch uint64(len(sigs)) + Instantiating - 1 {
	// case Unknown:
	// 	return nil, errors.New("transaction is of the status unknown")
	// case Initing:
	// 	return nil, errors.New("transaction is of the status initing")
	// case Inited:
	// 	return nil, errors.New("transaction is of the status inited")
	case Instantiating:
		nsb.InsuranceClaim(ses.GetGUID(), iter(atte, ses.Signer))
		return true, "", nil
	case Instantiated:
		chainID, tag, payload, err := serial_helper.UnserializeAttestationContent(atte.GetContent())

		if err != nil {
			return false, err.Error(), nil
		}

		// type = s.GetAtte().GetContent()
		// content = type.Content
		// s.BroadcastTxCommit(content)
		if isRawTransaction(tag) {
			cb, err := bn.RouteRaw(chainID, payload)
			fmt.Println("cbing ", string(cb))
			if err != nil {
				return false, err.Error(), nil
			}
		} else {
			cb, err := bn.Route(chainID, payload)
			fmt.Println("cbing ", string(cb))
			if err != nil {
				return false, err.Error(), nil
			}
		}

		nsb.InsuranceClaim(ses.GetGUID(), iter(atte, ses.Signer))

		return true, "", nil
	case Open:
		nsb.InsuranceClaim(ses.GetGUID(), iter(atte, ses.Signer))
		return true, "", nil
	case Opened:
		nsb.InsuranceClaim(ses.GetGUID(), iter(atte, ses.Signer))
		ses.UnderTransacting++
		return true, "", nil
	case Closed:
		ses.UnderTransacting++
		return true, "", nil
	default:
		return false, "", errors.New("unknown aid types")
	}
}

func (ses *SerialSession) SyncFromISC() (err error) {
	return errors.New("TODO")
}

// the database which used by others
type SerialSessionBase struct {
	sesref map[uint32]*SerialSession
}

func NewSerialSessionBase() *SerialSessionBase {
	return &SerialSessionBase{
		sesref: make(map[uint32]*SerialSession),
	}
}

func (sb *SerialSessionBase) InsertSessionInfo(
	db types.MultiIndex, idb types.Index, session types.Session,
) error {
	if sesref[]
	err := sb.InsertSessionAccounts(idb, session.GetGUID(), session.GetAccounts())
	if err != nil {
		return err
	}
	for idx, tx := range session.GetTransactions() {
		err = sb.InsertTransaction(idb, session.GetGUID(), uint64(idx), tx)
		if err != nil {
			return err
		}
	}
	return db.Insert(session.(*SerialSession))
}

func (sb *SerialSessionBase) FindSessionInfo(
	db types.MultiIndex, idb types.Index, isc_address []byte,
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
	sb.FindSessionAccounts(idb, isc_address, func(arg1 uint64, arg2 []byte) error {
		f[0].Accounts = append(f[0].Accounts, &account.PureAccount{ChainId: arg1, Address: arg2})
		fmt.Println("finded", hex.EncodeToString(arg2))
		return nil
	})
	for idx := uint64(f[0].TransactionCount); idx != 0; idx-- {
		sb.FindTransaction(idb, isc_address, idx, func(arg []byte) error {
			f[0].Transactions = append(f[0].Transactions, arg)
			fmt.Println("finded tx", len(arg))
			return nil
		})
	}
	session = &f[0]
	fmt.Println("getid", session.GetID())
	return
}

func (sb *SerialSessionBase) UpdateSessionInfo(
	db types.MultiIndex, idb types.Index, session types.Session,
) (err error) {
	fmt.Println("updateid", session.GetID())
	return db.Modify(session, session.ToKVMap())
}

func (sb *SerialSessionBase) DeleteSessionInfo(
	db types.MultiIndex, idb types.Index, isc_address []byte,
) (err error) {
	return db.Delete(&SerialSession{ISCAddress: isc_address})
}

func (sb *SerialSessionBase) InsertSessionAccounts(
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

func (sb *SerialSessionBase) FindSessionAccounts(
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

// func (sb *SerialSessionBase) InsertTransaction(
// 	db types.Index, transaction_id uint64, Transaction []byte
// ) (err error) (
//
// )
/***********************************Tx***********************/
func (sb *SerialSessionBase) InsertTransaction(
	db types.Index, isc_address []byte, transaction_id uint64, transaction []byte,
) (err error) {
	var k []byte
	k, err = serial_helper.DecoratePrefix([]byte{
		uint8((transaction_id >> 56) & 0xff), uint8((transaction_id >> 48) & 0xff),
		uint8((transaction_id >> 40) & 0xff), uint8((transaction_id >> 32) & 0xff),
		uint8((transaction_id >> 24) & 0xff), uint8((transaction_id >> 16) & 0xff),
		uint8((transaction_id >> 8) & 0xff), uint8((transaction_id >> 0) & 0xff),
	}, isc_address)
	if err != nil {
		return
	}
	k, err = serial_helper.DecoratePrefix(const_prefix.TransactionPrefix, k)
	//TransactionPrefix = []byte("ts")
	if err != nil {
		return
	}
	err = db.Set(k, transaction)
	return
}

func (sb *SerialSessionBase) FindTransaction(
	db types.Index, isc_address []byte, transaction_id uint64, getter func([]byte) error,
) (err error) {
	var k, v []byte
	k, err = serial_helper.DecoratePrefix([]byte{
		uint8((transaction_id >> 56) & 0xff), uint8((transaction_id >> 48) & 0xff),
		uint8((transaction_id >> 40) & 0xff), uint8((transaction_id >> 32) & 0xff),
		uint8((transaction_id >> 24) & 0xff), uint8((transaction_id >> 16) & 0xff),
		uint8((transaction_id >> 8) & 0xff), uint8((transaction_id >> 0) & 0xff),
	}, isc_address)
	if err != nil {
		return
	}
	k, err = serial_helper.DecoratePrefix(const_prefix.TransactionPrefix, k)
	if err != nil {
		return
	}
	v, err = db.Get(k)
	if err != nil {
		return
	}
	err = getter(v)
	return
}
