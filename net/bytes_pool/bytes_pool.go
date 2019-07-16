package nsbcli

import "sync"

type BytesPool struct {
	*sync.Pool
}
type poolnewer func() interface{}

func MakeNewBytesFunc(maxBytesSize int) poolnewer {
	return func() interface{} {
		return make([]byte, maxBytesSize)
	}
}

func NewBytesPool(maxBytesSize int) *BytesPool {
	return &BytesPool{
		Pool: &sync.Pool{
			New: MakeNewBytesFunc(maxBytesSize),
		},
	}
}
