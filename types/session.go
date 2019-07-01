package types

type SessionInterface interface {
	InsertToDB(db MultiIndex) error
	SelectFromDB(sessionId []byte, db MultiIndex) error
}
