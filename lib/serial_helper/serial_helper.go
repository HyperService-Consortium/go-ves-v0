package SerialHelper

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"

	types "github.com/HyperService-Consortium/go-uip/uiptypes"
)

var (
	errInsufficientBytes = errors.New("insufficient bytes to read")
)

// SerializeAccountsInterfaceBuffer write bytes with following bit view to buffer
//                                  1  1  1  1  1  1
//    0  1  2  3  4  5  6  7  8  9  0  1  2  3  4  5
//  +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//  |                                               |
//  |                    chainID                    |
//  |                                               |
//  |                                               |
//  +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//  |                                               |
//  |                    length                     |
//  |                                               |
//  |                                               |
//  +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//  |                                               |
//  /                    address                    /
//  /                                               /
//  +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
func SerializeAccountsInterfaceBuffer(buf io.ReadWriter, account types.Account) error {
	err := binary.Write(buf, binary.BigEndian, account.GetChainId())
	if err != nil {
		return err
	}

	var bc = account.GetAddress()
	err = binary.Write(buf, binary.BigEndian, int64(len(bc)))
	if err != nil {
		return err
	}
	buf.Write(bc)
	return nil
}

// SerializeAccountInterface return bytes with same bit view as
// SerializeAccountsInterfaceBuffer
func SerializeAccountInterface(account types.Account) ([]byte, error) {
	var buf = new(bytes.Buffer)
	err := SerializeAccountsInterfaceBuffer(buf, account)
	return buf.Bytes(), err
}

// SerializeAccountsInterface return bytes with same bit view as
// SerializeAccountsInterfaceBuffer
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

// UnserializeAccountInterface return bytes with same bit view as
// SerializeAccountsInterfaceBuffer
func UnserializeAccountInterface(b []byte) (n int64, chainID uint64, address []byte, err error) {
	var buf = bytes.NewBuffer(b)
	var ilen int64
	err = binary.Read(buf, binary.BigEndian, &chainID)
	if err != nil {
		return
	}
	err = binary.Read(buf, binary.BigEndian, &ilen)
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
		err = errInsufficientBytes
		return
	}
	n += 16
	return
}

// SerializeAttestationContent make bytes with following bit view
//                                  1  1  1  1  1  1
//    0  1  2  3  4  5  6  7  8  9  0  1  2  3  4  5
//  +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//  |                                               |
//  |                    chainID                    |
//  |                                               |
//  |                                               |
//  +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//  |                      tag                      |
//  +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//  |                                               |
//  /                    payload                    /
//  /                                               /
//  +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
func SerializeAttestationContent(chainID uint64, tag uint8, payload []byte) ([]byte, error) {
	var buf = new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, chainID)
	if err != nil {
		return nil, err
	}
	err = binary.Write(buf, binary.BigEndian, tag)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(payload)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// UnserializeAttestationContent recover information from bytes
func UnserializeAttestationContent(content []byte) (chainID uint64, tag uint8, payload []byte, err error) {
	var buf = bytes.NewBuffer(content)
	err = binary.Read(buf, binary.BigEndian, &chainID)
	if err != nil {
		return
	}
	err = binary.Read(buf, binary.BigEndian, &tag)
	if err != nil {
		return
	}
	payload = buf.Bytes()
	return
}

// DecoratePrefix concat two bytes, i.e. result = pre + b.
func DecoratePrefix(pre, b []byte) ([]byte, error) {
	var buf = bytes.NewBuffer(pre)
	_, err := buf.Write(b)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
