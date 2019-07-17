package wsrpc

import (
	"bytes"
	"encoding/binary"
	"sync"

	"github.com/gogo/protobuf/proto"
)

const maxSize = 65536

type BufferPool struct {
	*sync.Pool
	maxBufferSize int
}

func NewBufferPool(maxBufferSize int) *BufferPool {
	return &BufferPool{Pool: &sync.Pool{
		New: func() interface{} { return bytes.NewBuffer(make([]byte, 0, maxBufferSize)) },
	},
		maxBufferSize: maxBufferSize,
	}
}

type Serializer struct {
	bufferPool *BufferPool
}

var serial *Serializer

func NewSerializer(maxBufferSize int) *Serializer {
	return &Serializer{bufferPool: NewBufferPool(maxBufferSize)}
}

func (ser *Serializer) Serial(msgid uint16, msg proto.Message) (*bytes.Buffer, error) {
	var qwq = ser.bufferPool.Get().(*bytes.Buffer)
	qwq.Reset()

	binary.Write(qwq, binary.BigEndian, msgid)

	b, err := proto.Marshal(msg)
	if err != nil {
		return nil, err
	}
	qwq.Write(b)
	return qwq, nil
}

func (ser *Serializer) Put(buf *bytes.Buffer) bool {
	if buf.Cap() < ser.bufferPool.maxBufferSize {
		return false
	}
	ser.bufferPool.Put(buf)
	return true
}

func GetDefaultSerializer() *Serializer {
	return serial
}

func init() {
	serial = NewSerializer(maxSize)
}
