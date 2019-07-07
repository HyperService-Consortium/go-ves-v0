package SerialHelper

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"

	types "github.com/Myriad-Dreamin/go-ves/types"
)

var (
	insufficientBytes = errors.New("insufficient bytes to read")
)

func SerializeAccountsInterfaceBuffer(buf io.ReadWriter, account types.Account) error {
	err := binary.Write(buf, binary.LittleEndian, account.GetChainId())
	if err != nil {
		return err
	}

	var bc = account.GetAddress()
	err = binary.Write(buf, binary.LittleEndian, int64(len(bc)))
	if err != nil {
		return err
	}
	buf.Write(bc)
	return nil
}

func SerializeAccountInterface(account types.Account) ([]byte, error) {
	var buf = new(bytes.Buffer)
	err := SerializeAccountsInterfaceBuffer(buf, account)
	return buf.Bytes(), err
}

func SerializeAccountsInterface(accounts []types.Account) ([]byte, error) {
	var buf = new(bytes.Buffer)
	for _, account := range accounts {
		err := SerializeAccountsInterfaceBuffer(buf, account)
		if err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}

func UnserializeAccountInterface(b []byte) (n int64, chain_type uint64, address []byte, err error) {
	var buf = bytes.NewBuffer(b)
	var ilen int64
	err = binary.Read(buf, binary.LittleEndian, &chain_type)
	if err != nil {
		return
	}
	err = binary.Read(buf, binary.LittleEndian, &ilen)
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
		err = insufficientBytes
		return
	}
	n += 16
	return
}

func DecoratePrefix(pre, b []byte) ([]byte, error) {
	var buf = bytes.NewBuffer(pre)
	_, err := buf.Write(b)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
