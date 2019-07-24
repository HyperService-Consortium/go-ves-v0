package types

import uiptypes "github.com/Myriad-Dreamin/go-uip/types"

type isc_address = []byte

type NSBInterface interface {
	SaveAttestation(isc_address, uiptypes.Attestation) error
	InsuranceClaim(isc_address, uiptypes.Attestation) error
	SettleContract(isc_address) error
}

type destination = uint64
type payload = []byte
type on_chain_transaction = []byte
type cb_info = []byte
type BNInterface interface {
	RouteRaw(destination, payload) (cb_info, error)
	Route(destination, on_chain_transaction) (cb_info, error)
}

type success_or_not = bool
type help_info = string
type Session interface {
	// session must has isc_address or other guid

	// session is a kv-object
	KVObject

	SetSigner(uiptypes.Signer)

	GetGUID() isc_address
	GetAccounts() []uiptypes.Account
	GetTransaction(transaction_local_id) transaction
	GetTransactions() []transaction

	GetAckCount() uint32
	GetTransactingTransaction() (transaction_local_id, error)

	// error reports Internal errors, help_info reports Logic errors
	InitFromOpIntents(uiptypes.OpIntents) (success_or_not, help_info, error)
	AckForInit(uiptypes.Account, uiptypes.Signature) (success_or_not, help_info, error)
	ProcessAttestation(NSBInterface, BNInterface, uiptypes.Attestation) (success_or_not, help_info, error)

	SyncFromISC() error
}

// the database which used by others

type transaction_id = uint64
type getter = func([]byte) error
type SessionBase interface {

	// insert accounts maps from guid to account
	InsertSessionInfo(MultiIndex, Index, Session) error

	// find accounts which guid is corresponding to user
	FindSessionInfo(MultiIndex, Index, isc_address) (Session, error)

	UpdateSessionInfo(MultiIndex, Index, Session) error

	DeleteSessionInfo(MultiIndex, Index, isc_address) error

	FindTransaction(Index, isc_address, transaction_id, getter) error
}
