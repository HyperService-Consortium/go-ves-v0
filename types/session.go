package types

type isc_address = []byte

type success_or_not = bool
type help_info = string
type Session interface {
	// session must has isc_address or other guid

	// session is a kv-object
	KVObject

	GetGUID() isc_address
	GetAccounts() []Account
	GetTransaction(transaction_id) transaction
	GetTransactions() []transaction

	IsSyncing() bool
	GetTransactingTransaction() (transaction_id, error)

	// error reports Internal errors, help_info reports Logic errors
	InitFromOpIntents(OpIntents) (success_or_not, help_info, error)
	AckForInit(Account, Signature) (success_or_not, help_info, error)
	ProcessAttestation(Attestation) (success_or_not, help_info, error)

	SyncFromISC() error
}

// the database which used by others
type SessionBase interface {
	// insert accounts maps from guid to account
	InsertSession(MultiIndex, Session) error

	// find accounts which guid is corresponding to user
	FindSession(MultiIndex, isc_address) (Session, error)

	UpdateSession(MultiIndex, Session) error

	DeleteSession(MultiIndex, isc_address) error
}
