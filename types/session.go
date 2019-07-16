package types

import uiptypes "github.com/Myriad-Dreamin/go-uip/types"

type isc_address = []byte

type success_or_not = bool
type help_info = string
type Session interface {
	// session must has isc_address or other guid

	// session is a kv-object
	KVObject

	GetGUID() isc_address
	GetAccounts() []uiptypes.Account
	GetTransaction(transaction_local_id) transaction
	GetTransactions() []transaction

	GetTransactingTransaction() (transaction_local_id, error)

	// error reports Internal errors, help_info reports Logic errors
	InitFromOpIntents(uiptypes.OpIntents) (success_or_not, help_info, error)
	AckForInit(uiptypes.Account, uiptypes.Signature) (success_or_not, help_info, error)
	ProcessAttestation(uiptypes.Attestation) (success_or_not, help_info, error)

	SyncFromISC() error
}

// the database which used by others
type SessionBase interface {

	// insert accounts maps from guid to account
	InsertSessionInfo(MultiIndex, Session) error

	// find accounts which guid is corresponding to user
	FindSessionInfo(MultiIndex, isc_address) (Session, error)

	UpdateSessionInfo(MultiIndex, Session) error

	DeleteSessionInfo(MultiIndex, isc_address) error
}
