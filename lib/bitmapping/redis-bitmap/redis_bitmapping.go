package bitmap

import (
	"bytes"
	"encoding/binary"
	"unsafe"

	"github.com/gomodule/redigo/redis"
)

var referbit [8]byte

type BitMap struct {
	name   []byte
	Conn   redis.Conn
	length int64
}

func NewBitMap(Name []byte, length int64, conn redis.Conn) (*BitMap, error) {
	var b = make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(length))
	s := bytes.NewBuffer(Name)
	s.WriteByte(':')
	s.WriteByte('l')
	_, err := conn.Do("set", s.Bytes(), b)
	if err != nil {
		return nil, err
	}
	return &BitMap{
		name:   Name,
		Conn:   conn,
		length: length,
	}, nil
}

func GetBitMap(Name []byte, conn redis.Conn) *BitMap {
	return &BitMap{
		name: Name,
		Conn: conn,
	}
}

func PutBitMapLength(Name []byte, length int64, conn redis.Conn) error {
	var b = make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(length))
	s := bytes.NewBuffer(Name)
	s.WriteByte(':')
	s.WriteByte('l')
	_, err := conn.Do("set", s.Bytes(), b)
	if err != nil {
		return err
	}
	return nil
}

func (b *BitMap) Get(idx int64) (bool, error) {
	a, err := b.Conn.Do("getbit", b.name, idx)
	if err != nil {
		return false, err
	}
	return a.(int64) != 0, nil
}

func (b *BitMap) Set(idx int64) (bool, error) {
	a, err := b.Conn.Do("setbit", b.name, idx, 1)
	if err != nil {
		return false, err
	}
	return a.(int64) != 0, nil
}

func (b *BitMap) Reset(idx int64) (bool, error) {
	a, err := b.Conn.Do("setbit", b.name, idx, 0)
	if err != nil {
		return false, err
	}
	return a.(int64) != 0, nil
}

func (b *BitMap) Clear() (bool, error) {
	a, err := b.Conn.Do("del", b.name)
	if err != nil {
		return false, err
	}
	return a.(int64) != 0, nil
}

func (b *BitMap) Delete() (bool, error) {
	a, err := b.Conn.Do("del", b.name)
	if err != nil {
		return false, err
	}
	a, err = b.Conn.Do("del", *(*string)(unsafe.Pointer(&b.name))+":l")
	if err != nil {
		return false, err
	}
	return a.(int64) != 0, nil
}

func (b *BitMap) InLength(idx int64) (bool, error) {
	if b.length == 0 {
		a, err := b.Conn.Do("get", *(*string)(unsafe.Pointer(&b.name))+":l")
		if err != nil {
			return false, err
		}
		b.length = int64(binary.BigEndian.Uint64(a.([]byte)))
	}
	return idx < b.length, nil
}

func (b *BitMap) Length() (int64, error) {
	if b.length == 0 {
		a, err := b.Conn.Do("get", *(*string)(unsafe.Pointer(&b.name))+":l")
		if err != nil {
			return 0, err
		}
		b.length = int64(binary.BigEndian.Uint64(a.([]byte)))
	}
	return b.length, nil
}

func (b *BitMap) Count() (int64, error) {
	a, err := b.Conn.Do("bitcount", b.name)
	if err != nil {
		return 0, err
	}
	return a.(int64), nil
}

func init() {
	for idx := uint(0); idx < 8; idx++ {
		referbit[idx] = (byte(1) << idx)
	}
}
