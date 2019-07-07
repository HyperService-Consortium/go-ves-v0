package session

import "github.com/Myriad-Dreamin/go-ves/types"

type SerialSession struct {
	SessionId []byte
	CurLoad   int32
	TxCount   int32
}

func (ses *SerialSession) InsertToDB(db types.MultiIndex) error {
	return nil
}
func (ses *SerialSession) SelectFromDB(sessionId []byte, db types.MultiIndex) error {
	return nil
}

func NewSerialSessionFromDB(
	sessionId []byte,
	db types.MultiIndex,
) (ses *SerialSession, err error) {
	ses = new(SerialSession)
	err = ses.SelectFromDB(sessionId, db)
	return
}
