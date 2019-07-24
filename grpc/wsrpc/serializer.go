package wsrpc

import (
	"bytes"
	"encoding/binary"
	"sync"

	"github.com/gogo/protobuf/proto"
)

const maxSize = 65536

// BufferPool Object
// accept temporary buffer object whose size is larger than maxBufferSize
// provide the steady source of temporary buffer.
type BufferPool struct {
	*sync.Pool
	maxBufferSize int
}

// NewBufferPool return a pointer of BufferPool
func NewBufferPool(maxBufferSize int) *BufferPool {
	return &BufferPool{Pool: &sync.Pool{
		New: func() interface{} {
			// buf :=
			// fmt.Printf("gee %p", buf)
			return bytes.NewBuffer(make([]byte, 0, maxBufferSize))
		},
	},
		maxBufferSize: maxBufferSize,
	}
}

// Put a buffer into this BufferPool. the buffer whose size is smaller than the
// bufferPool.maxBufferSize will be ignored
func (bufpool *BufferPool) Put(buf interface{}) {
	if buf, ok := buf.(*bytes.Buffer); ok {
		// fmt.Printf("puu %p", buf)
		if buf.Cap() >= bufpool.maxBufferSize {
			buf.Reset()
			bufpool.Pool.Put(buf)
		}
	}
}

// Serializer works for provide the environment of serializing protobuf
// messages
type Serializer struct {
	bufferPool *BufferPool
}

var serial *Serializer

// NewSerializer return a pointer of Serializer
func NewSerializer(maxBufferSize int) *Serializer {
	return &Serializer{bufferPool: NewBufferPool(maxBufferSize)}
}

// Serial concat the msgid and serialized msg
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

// Put a buffer into its buffer pool
func (ser *Serializer) Put(buf *bytes.Buffer) {
	ser.bufferPool.Put(buf)
}

// GetDefaultSerializer gets the default serializer singleton
func GetDefaultSerializer() *Serializer {
	return serial
}

func init() {
	serial = NewSerializer(maxSize)
}
