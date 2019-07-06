package types

type KVObject interface {
	GetObjectPtr() interface{}
	GetSlicePtr() interface{}
	GetId() int64
}
