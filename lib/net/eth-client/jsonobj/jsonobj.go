package jsonobj

import (
	"bytes"
	"encoding/hex"
	"strconv"

	bytespool "github.com/Myriad-Dreamin/object-pool/bytes-pool"
)

const (
	maxBufferSize = 1024
	splitByte     = ','
	ssplitByte    = '"'
	cbx           = `","`
	endParamByte  = ']'
	endJSONByte   = '}'

	maxBytesSize = 1024
)

var (
	reqGetAccount      = []byte(`{"id":1,"jsonrpc":"2.0","method":"eth_accounts","params":[]}`)
	reqPersonalUnlock  = []byte(`{"id":64,"jsonrpc":"2.0","method":"personal_unlockAccount","params":["`)
	reqSendTransaction = []byte(`{"id":1,"jsonrpc":"2.0","method":"eth_sendTransaction","params":[`)
	reqGetStorageAt    = []byte(`{"id":1,"jsonrpc":"2.0","method":"eth_getStorageAt","params":[`)
	hexPrefix          = "0x"
	bp                 = bytespool.NewMultiThreadBytesPool(maxBytesSize)
)

// GetAccount return all accounts on eth
func GetAccount() []byte {
	return reqGetAccount
}

// GetPersonalUnlock return whether unlocked
// do not send too long passphrase
func GetPersonalUnlock(addr string, passphrase string, duration int) []byte {
	var b = bp.Get()
	var buf = bytes.NewBuffer(b)
	buf.Reset()

	buf.Write(reqPersonalUnlock)

	buf.WriteString(addr)

	buf.WriteString(cbx)

	buf.WriteString(passphrase)

	buf.WriteByte(ssplitByte)
	buf.WriteByte(splitByte)

	buf.WriteString(strconv.Itoa(duration))

	buf.WriteByte(endParamByte)
	buf.WriteByte(endJSONByte)

	return buf.Bytes()
}

// GetSendTransaction return whether unlocked
// do not send too long obj
func GetSendTransaction(obj []byte) []byte {
	var b = bp.Get()
	var buf = bytes.NewBuffer(b)
	buf.Reset()

	buf.Write(reqSendTransaction)

	buf.Write(obj)

	buf.WriteByte(endParamByte)
	buf.WriteByte(endJSONByte)

	return buf.Bytes()
}

// GetStorageAt return whether unlocked
// do not send too long obj
func GetStorageAt(address, pos []byte, tag string) []byte {
	var b = bp.Get()
	var buf = bytes.NewBuffer(b)
	buf.Reset()

	buf.Write(reqGetStorageAt)

	buf.WriteByte(ssplitByte)

	buf.WriteString(hexPrefix)
	buf.WriteString(hex.EncodeToString(address))

	buf.WriteString(cbx)

	buf.WriteString(hexPrefix)
	buf.WriteString(hex.EncodeToString(pos))

	buf.WriteString(cbx)

	buf.WriteString(tag)

	buf.WriteByte(ssplitByte)

	buf.WriteByte(endParamByte)
	buf.WriteByte(endJSONByte)

	return buf.Bytes()
}

// ReturnBytes to Pool
func ReturnBytes(b []byte) {
	bp.Put(b)
}
