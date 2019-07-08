package nsbcli

import "sync"

type BytesPool struct {
	*sync.Pool
}

func newBytes() interface{} {
	return make([]byte, maxBytesSize)
}

func NewBytesPool() *BytesPool {
	return &BytesPool{
		Pool: &sync.Pool{
			New: newBytes,
		},
	}
}
