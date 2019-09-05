package session

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"sync"

	TransType "github.com/HyperService-Consortium/go-uip/const/trans_type"
	TxState "github.com/HyperService-Consortium/go-uip/const/transaction_state_type"
	opintents "github.com/HyperService-Consortium/go-uip/op-intent"
	uiptypes "github.com/HyperService-Consortium/go-uip/types"
	account "github.com/HyperService-Consortium/go-uip/types/account"
	bitmap "github.com/HyperService-Consortium/go-ves/lib/bitmapping/redis-bitmap"
	const_prefix "github.com/HyperService-Consortium/go-ves/lib/database/const_prefix"
	redispool "github.com/HyperService-Consortium/go-ves/lib/database/redis"
	log "github.com/HyperService-Consortium/go-ves/lib/log"
	serial_helper "github.com/HyperService-Consortium/go-ves/lib/serial_helper"
	types "github.com/HyperService-Consortium/go-ves/types"
)

type MultiThreadSerialSession struct {
	ID               int64  `json:"-" xorm:"pk unique notnull autoincr 'id'"`
	ISCAddress       []byte `json:"-" xorm:"notnull 'isc_address'"`
	TransactionCount uint32 `json:"-" xorm:"'transaction_count'"`
	UnderTransacting uint32 `json:"-" xorm:"'under_transacting'"`
	Status           uint8  `json:"-" xorm:"'status'"`
	Content          []byte `json:"-" xorm:"'content'"`

	// index
	Accounts     []uiptypes.Account `json:"-" xorm:"-"`
	Transactions [][]byte           `json:"transactions" xorm:"-"`

	// Acks     []byte `json:"-" xorm:"'-'"`
	// AckCount uint32 `json:"-" xorm:"'-'"`

	// handler
	Signer uiptypes.Signer `json:"-" xorm:"-"`
}

func randomMultiThreadSerialSession() *MultiThreadSerialSession {
	var buf = make([]byte, 20)
	binary.PutVarint(buf, rand.Int63())
	return &MultiThreadSerialSession{
		ISCAddress: buf,
	}
}

func (ses *MultiThreadSerialSession) TableName() string {
	return "mves_session"
}

func (ses *MultiThreadSerialSession) ToKVMap() map[string]interface{} {
	return map[string]interface{}{
		"id":                ses.ID,
		"isc_address":       ses.ISCAddress,
		"transaction_count": ses.TransactionCount,
		"under_transacting": ses.UnderTransacting,
		"status":            ses.Status,
		"content":           ses.Content,
		// "acks":              ses.Acks,
		// "ack_count":         ses.AckCount,
	}
}

func (ses *MultiThreadSerialSession) GetID() int64 {
	return ses.ID
}

func (ses *MultiThreadSerialSession) GetGUID() (isc_address []byte) {
	return ses.ISCAddress
}

func (ses *MultiThreadSerialSession) GetObjectPtr() interface{} {
	return new(MultiThreadSerialSession)
}

func (ses *MultiThreadSerialSession) GetSlicePtr() interface{} {
	return new([]MultiThreadSerialSession)
}

func (ses *MultiThreadSerialSession) GetAccounts() []uiptypes.Account {

	// move to adapdator
	// if ses.Accounts == nil {
	// 	ses.Accounts = nil
	// }

	return ses.Accounts
}

func (ses *MultiThreadSerialSession) GetAckCount() uint32 {
	count, err := bitmap.GetBitMap(ses.GetGUID(), redispool.RedisCacheClient.Pool.Get()).Count()
	if err != nil {
		log.Debugln("debugging of get ackcount", err)
		return 0
	}
	return uint32(count)
}

func (ses *MultiThreadSerialSession) GetTransaction(transaction_id uint32) []byte {
	// if ses.Transactions == nil {
	// 	ses.Transactions = make([][]byte, ses.TransactionCount)
	// }
	// if ses.Transactions[transaction_id] == nil {
	// 	ses.Transactions[transaction_id] = nil
	// }

	return ses.Transactions[transaction_id]
}

func (ses *MultiThreadSerialSession) GetTransactions() (transactions [][]byte) {
	// if ses.Transactions == nil {
	// 	ses.Transactions = make([][]byte, ses.TransactionCount)
	// }

	return ses.Transactions
}

func (ses *MultiThreadSerialSession) GetTransactingTransaction() (transaction_id uint32, err error) {
	// Status
	return ses.UnderTransacting, nil
}

func (ses *MultiThreadSerialSession) GetContent() []byte {
	return ses.Content
}

func (ses *MultiThreadSerialSession) SetSigner(signer uiptypes.Signer) {
	ses.Signer = signer
}

func (ses *MultiThreadSerialSession) InitFromOpIntents(opIntents uiptypes.OpIntents) (bool, string, error) {
	intents, _, err := opintents.NewOpIntentInitializer().InitOpIntent(opIntents)
	if err != nil {
		return false, err.Error(), nil
	}
	ses.Transactions = make([][]byte, 0, len(intents))

	ses.Accounts = nil
	c := makeComparator()
	ses.Accounts = append(ses.Accounts, &account.PureAccount{ChainId: 4, Address: ses.Signer.GetPublicKey()})
	for _, intent := range intents {
		ses.Transactions = append(ses.Transactions, intent.Bytes())

		if c.Insert(intent.ChainID, intent.Src) {
			ses.Accounts = append(ses.Accounts, &account.PureAccount{ChainId: intent.ChainID, Address: intent.Src})
		}
		if intent.TransType == TransType.Payment && c.Insert(intent.ChainID, intent.Dst) {
			ses.Accounts = append(ses.Accounts, &account.PureAccount{ChainId: intent.ChainID, Address: intent.Dst})
		}
	}
	ses.TransactionCount = uint32(len(intents))
	ses.UnderTransacting = 0
	ses.Status = 0
	// ses.Acks = make([]byte, (len(ses.Accounts)+7)>>3)
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

func (ses *MultiThreadSerialSession) AfterInitGUID() error {
	return bitmap.PutBitMapLength(ses.ISCAddress, int64(len(ses.Accounts)), redispool.RedisCacheClient.Pool.Get())
}

func (ses *MultiThreadSerialSession) AckForInit(
	account uiptypes.Account,
	signature uiptypes.Signature,
) (success_or_not bool, help_info string, err error) {
	var addr, acks = account.GetAddress(), bitmap.GetBitMap(ses.ISCAddress, redispool.RedisCacheClient.Pool.Get())
	// fmt.Println(ses.Acks, len(ses.Acks))
	log.Printf("acked:")
	log.Println(acks.Count())
	log.Printf("ack:")
	log.Println(acks.Length())
	log.Printf("wanting:")
	log.Println(len(ses.GetAccounts()))
	for idx, ak := range ses.GetAccounts() {
		if bytes.Equal(ak.GetAddress(), addr) {
			if ok, err := acks.InLength(int64(idx)); err != nil {
				return false, "", fmt.Errorf("internal error: getting length %v", err)
			} else if !ok {
				return false, "", errors.New("wrong Acks bytes set..")
			}
			if setted, err := acks.Get(int64(idx)); err != nil {
				return false, "", fmt.Errorf("internal error: getting bit %v", err)
			} else if setted {
				return false, "have acked", nil
			}
			if !Verify(signature, ses.GetContent(), account.GetAddress()) {
				return false, "verify signature error...", nil
			}

			if lastbit, err := acks.Set(int64(idx)); err != nil {
				return false, "", fmt.Errorf("internal error: setting bit %v", err)
			} else if lastbit {
				return false, "conflict set at the same time", nil
			}

			if count, err := acks.Count(); err != nil {
				return false, "", fmt.Errorf("internal error: counting bit %v", err)
			} else if count == int64(len(ses.Accounts)) {
				log.Infoln("session setup finished")
			}
			// todo: NSB
			return true, "", nil
		}
	}
	return false, "account not found in this session", nil
}

func (ses *MultiThreadSerialSession) NotifyAttestation(
	nsb types.NSBInterface, bn uiptypes.BlockChainInterface, atte uiptypes.Attestation,
) (success_or_not bool, help_info string, err error) {
	// todo
	tid := atte.GetTid()

	if tid != uint64(ses.UnderTransacting) {
		return false, "this transaction is not undertransacting", nil
	}
	fmt.Println("notifying")

	switch atte.GetAid() {
	// case Unknown:
	// 	return nil, errors.New("transaction is of the status unknown")
	// case Initing:
	// 	return nil, errors.New("transaction is of the status initing")
	// case Inited:
	// 	return nil, errors.New("transaction is of the status inited")
	case TxState.Instantiating:
		// err := nsb.InsuranceClaim(ses.GetGUID(), iter(atte, ses.Signer))
		// if err != nil {
		// 	return false, "", err
		// }
		return true, "", nil
	case TxState.Instantiated:
		// chainID, tag, payload, err := serial_helper.UnserializeAttestationContent(atte.GetContent())
		//
		// if err != nil {
		// 	return false, err.Error(), nil
		// }

		// type = s.GetAtte().GetContent()
		// content = type.Content
		// s.BroadcastTxCommit(content)
		// if isRawTransaction(tag) {
		// 	cb, err := bn.RouteRaw(chainID, payload)
		// 	log.Infof("cbing %v", cb)
		// 	if err != nil {
		// 		return false, err.Error(), nil
		// 	}
		// } else {
		// 	cb, err := bn.Route(chainID, payload)
		// 	log.Infof("cbing %v", cb)
		// 	if err != nil {
		// 		return false, err.Error(), nil
		// 	}
		// }

		// err = nsb.InsuranceClaim(ses.GetGUID(), iter(atte, ses.Signer))
		// if err != nil {
		// 	return false, "", err
		// }

		return true, "", nil
	case TxState.Open:
		// err := nsb.InsuranceClaim(ses.GetGUID(), iter(atte, ses.Signer))
		// if err != nil {
		// 	return false, "", err
		// }
		return true, "", nil
	case TxState.Opened:
		// err := nsb.InsuranceClaim(ses.GetGUID(), iter(atte, ses.Signer))
		// if err != nil {
		// 	return false, "", err
		// }
		return true, "", nil
	case TxState.Closed:
		ses.UnderTransacting++
		if ses.UnderTransacting == ses.TransactionCount {
			err = nsb.SettleContract(ses.ISCAddress)

			if err != nil {
				return false, "", err
			}

		}
		return true, "closed session", nil
	default:
		return false, "", errors.New("unknown aid types")
	}
}

func (ses *MultiThreadSerialSession) ProcessAttestation(
	nsb types.NSBInterface, bn uiptypes.BlockChainInterface, atte uiptypes.Attestation,
) (success_or_not bool, help_info string, err error) {
	// todo

	tid := atte.GetTid()

	if tid != uint64(ses.UnderTransacting) {
		return false, "this transaction is not undertransacting", nil
	}

	switch atte.GetAid() {
	// case Unknown:
	// 	return nil, errors.New("transaction is of the status unknown")
	// case Initing:
	// 	return nil, errors.New("transaction is of the status initing")
	// case Inited:
	// 	return nil, errors.New("transaction is of the status inited")
	case TxState.Instantiating:
		err := nsb.InsuranceClaim(ses.GetGUID(), iter(atte, ses.Signer))
		if err != nil {
			return false, "", err
		}
		return true, "", nil
	case TxState.Instantiated:
		// chainID, tag, payload, err := serial_helper.UnserializeAttestationContent(atte.GetContent())
		//
		// if err != nil {
		// 	return false, err.Error(), nil
		// }

		// type = s.GetAtte().GetContent()
		// content = type.Content
		// s.BroadcastTxCommit(content)
		// if isRawTransaction(tag) {
		// 	cb, err := bn.RouteRaw(chainID, payload)
		// 	log.Infof("cbing %v", cb)
		// 	if err != nil {
		// 		return false, err.Error(), nil
		// 	}
		// } else {
		// 	cb, err := bn.Route(chainID, payload)
		// 	log.Infof("cbing %v", cb)
		// 	if err != nil {
		// 		return false, err.Error(), nil
		// 	}
		// }

		err = nsb.InsuranceClaim(ses.GetGUID(), iter(atte, ses.Signer))
		if err != nil {
			return false, "", err
		}

		return true, "", nil
	case TxState.Open:
		err := nsb.InsuranceClaim(ses.GetGUID(), iter(atte, ses.Signer))
		if err != nil {
			return false, "", err
		}
		return true, "", nil
	case TxState.Opened:
		err := nsb.InsuranceClaim(ses.GetGUID(), iter(atte, ses.Signer))
		if err != nil {
			return false, "", err
		}
		return true, "", nil
	case TxState.Closed:
		ses.UnderTransacting++
		if ses.UnderTransacting == ses.TransactionCount {
			err = nsb.SettleContract(ses.ISCAddress)
			if err != nil {
				return false, "", err
			}
		}
		return true, "", nil
	default:
		return false, "", errors.New("unknown aid types")
	}
}

func (ses *MultiThreadSerialSession) SyncFromISC() (err error) {
	return errors.New("TODO")
}

const MaxSessionCount = 500

type LabeledMutex struct {
	sync.Mutex
	Index int
}

// the database which used by others
type MultiThreadSerialSessionBase struct {
	mutex sync.Mutex
	// todo is it safe enough?
	alloc  map[uint64]*LabeledMutex
	resrc  chan int
	mutexs [MaxSessionCount]*LabeledMutex
	count  [MaxSessionCount]int
}

func NewMultiThreadSerialSessionBase() *MultiThreadSerialSessionBase {
	sb := &MultiThreadSerialSessionBase{
		alloc: make(map[uint64]*LabeledMutex),
		resrc: make(chan int, MaxSessionCount),
	}
	for idx := len(sb.mutexs) - 1; idx >= 0; idx-- {
		mutex := new(LabeledMutex)
		sb.mutexs[idx] = mutex
		mutex.Index = idx
		sb.resrc <- idx
	}
	return sb
}

// func RequestsStart() {
//     for{
//         select {
//         }
//     }
// }
//

func (sb *MultiThreadSerialSessionBase) InsertSessionInfo(
	db types.MultiIndex, idb types.Index, session types.Session,
) error {
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
	return db.Insert(session.(*MultiThreadSerialSession))
}

func (sb *MultiThreadSerialSessionBase) FindSessionInfo(
	db types.MultiIndex, idb types.Index, isc_address []byte,
) (session types.Session, err error) {
	var sessions interface{}
	sessions, err = db.Select(&MultiThreadSerialSession{ISCAddress: isc_address})
	if err != nil {
		return
	}
	f := sessions.([]MultiThreadSerialSession)
	if f == nil {
		return nil, errors.New("not found")
	}
	sb.FindSessionAccounts(idb, isc_address, func(arg1 uint64, arg2 []byte) error {
		f[0].Accounts = append(f[0].Accounts, &account.PureAccount{ChainId: arg1, Address: arg2})
		return nil
	})
	for idx := uint32(0); idx < f[0].TransactionCount; idx++ {
		sb.FindTransaction(idb, isc_address, uint64(idx), func(arg []byte) error {
			f[0].Transactions = append(f[0].Transactions, arg)
			return nil
		})
	}
	session = &f[0]
	return
}

func (sb *MultiThreadSerialSessionBase) UpdateSessionInfo(
	db types.MultiIndex, idb types.Index, session types.Session,
) (err error) {
	return db.Modify(session, session.ToKVMap())
}

func (sb *MultiThreadSerialSessionBase) DeleteSessionInfo(
	db types.MultiIndex, idb types.Index, isc_address []byte,
) (err error) {
	return db.Delete(&MultiThreadSerialSession{ISCAddress: isc_address})
}

func (sb *MultiThreadSerialSessionBase) InsertSessionAccounts(
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

func (sb *MultiThreadSerialSessionBase) FindSessionAccounts(
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

// func (sb *MultiThreadSerialSessionBase) InsertTransaction(
// 	db types.Index, transaction_id uint64, Transaction []byte
// ) (err error) (
//
// )
/***********************************Tx***********************/
func (sb *MultiThreadSerialSessionBase) InsertTransaction(
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

func (sb *MultiThreadSerialSessionBase) FindTransaction(
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

func (sb *MultiThreadSerialSessionBase) ActivateSession(isc_address []byte) {
	sb.mutex.Lock()
	idx := binary.BigEndian.Uint64(isc_address[0:8])
	if mutex := sb.alloc[idx]; mutex != nil {
		sb.count[mutex.Index]++
		sb.mutex.Unlock()
		mutex.Lock()
	} else {
		mutex = sb.mutexs[<-sb.resrc]
		sb.alloc[idx] = mutex
		sb.count[mutex.Index]++
		sb.mutex.Unlock()
		mutex.Lock()
	}
	return
}
func (sb *MultiThreadSerialSessionBase) InactivateSession(isc_address []byte) {
	sb.mutex.Lock()
	idx := binary.BigEndian.Uint64(isc_address[0:8])

	// for fast operation
	mutex := sb.alloc[idx]
	sb.count[mutex.Index]--
	if sb.count[mutex.Index] == 0 {
		sb.resrc <- mutex.Index
		sb.alloc[idx] = nil
	}
	sb.mutex.Unlock()
	mutex.Unlock()
	return
}
