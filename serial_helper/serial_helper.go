package SerialHelper

import (
	"bytes"
	"encoding/binary"
	"errors"

	types "github.com/Myriad-Dreamin/go-ves/types"
)

func SerializeAccountInterface(account types.Account) []byte {
	var buf = new(bytes.Buffer)
	var bx = make([]byte, 8)
	binary.PutUvarint(bx, account.GetChainType())
	buf.Write(bx)
	var bc = account.GetAddress()
	binary.PutVarint(bx, int64(len(bc)))
	buf.Write(bx)
	buf.Write(bc)
	return buf.Bytes()
}

func UnserializeAccountInterface(b []byte) (n int64, chain_type uint64, address []byte, err error) {
	var buf = bytes.NewBuffer(b)
	var ilen int64
	chain_type, err = binary.ReadUvarint(buf)
	if err != nil {
		return
	}
	ilen, err = binary.ReadVarint(buf)
	if err != nil {
		return
	}
	var nn int
	address = make([]byte, ilen)
	nn, err = buf.Read(address)
	if err != nil {
		return
	}
	n = int64(nn)
	if n < ilen {
		err = errors.New("insufficient bytes to read")
		return
	}
	n += 16
	return
}
