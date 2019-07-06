package session

import types "github.com/Myriad-Dreamin/go-ves/types"

type Session struct {
	ISCAddress   []byte
	Accounts     []types.Account
	Transactions [][]byte
}

func (ses *Session) GetGUID() (isc_address []byte) {
	return ses.ISCAddress
}

func (ses *Session) GetAccounts() (accounts []types.Account) {
	// ?
}

func (ses *Session) GetTransaction(transaction_id uint32) (transaction []byte) {

}

func (ses *Session) GetTransactions() (transactions [][]byte) {

}

func (ses *Session) IsSyncing() (syncing bool) {

}

func (ses *Session) GetTransactingTransaction() (transaction_id uint32, err error) {

}

func (ses *Session) AckForInit(
	account types.Account,
	signature types.Signature,
) (success_or_not bool, help_info string, err error) {

}

func (ses *Session) ProcessAttestation(
	attestation types.Attestation,
) (success_or_not bool, help_info string, err error) {

}

func (ses *Session) SyncFromISC() (err error) {

}

// the database which used by others
type SessionBase struct {
}

func (sb *SessionBase) InsertSession(
	db types.MultiIndex, session types.Session,
) (err error) {

}

func (sb *SessionBase) FindSession(
	db types.MultiIndex, isc_address []byte,
) (session types.Session, err error) {

}

func (sb *SessionBase) UpdateSession(
	db types.MultiIndex, session types.Session,
) (err error) {

}

func (sb *SessionBase) DeleteSession(
	db types.MultiIndex, isc_address []byte,
) (err error) {

}
