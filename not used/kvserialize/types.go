package kvs

type Marshalable interface {
	MarshalKV() interface{}
}

type MaxIntType = int64
type MaxUintType = uint64

type VInt interface {
	UnmarshalKV(int)
}

type VInt8 interface {
	UnmarshalKV(int8)
}

type VInt16 interface {
	UnmarshalKV(int16)
}

type VInt32 interface {
	UnmarshalKV(int32)
}

type VInt64 interface {
	UnmarshalKV(int64)
}

type VUint8 interface {
	UnmarshalKV(uint8)
}

type VUint16 interface {
	UnmarshalKV(uint16)
}

type VUint32 interface {
	UnmarshalKV(uint32)
}

type VUint64 interface {
	UnmarshalKV(uint64)
}

type VUint interface {
	UnmarshalKV(uint)
}

type VFloat64 interface {
	UnmarshalKV(float64)
}

type VString interface {
	UnmarshalKV(string)
}

type VBytes interface {
	UnmarshalKV([]byte)
}

type VBool interface {
	UnmarshalKV(bool)
}
