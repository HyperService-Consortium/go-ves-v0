package types

type KVPair struct {
	Key   string
	Value interface{}
}

type KVMap = map[string]interface{}

type KVObject interface {
	GetObjectPtr() interface{}
	GetSlicePtr() interface{}
	GetId() int64
}

type MultiIndex interface {
	// RegisterObject(KVObject) error

	Insert(KVObject) error

	Select(KVObject) (interface{}, error)

	SelectAll(KVObject) (interface{}, error)

	// 要求只Delete到一个
	Delete(KVObject) error

	// 可以Delete多个
	MultiDelete(KVObject) error

	Modify(KVObject, KVMap) error

	MultiModify(KVObject, KVMap) error
}

type KVPMultiIndex interface {
	Insert(...KVPair) error

	Select([]interface{}, ...KVPair) error

	// 要求只Delete到一个
	Delete(...KVPair) error
	// 可以Delete多个
	MultiDelete(...KVPair) error

	// 要求只Update到一个
	Modify([]KVPair, ...KVPair) error
	// 可以Update到多个
	MultiModify([]KVPair, ...KVPair) error
}

type ORMMultiIndex interface {
	MultiIndex
	// 要求只Update到一个
	// Modify(ORMObject, ORMObject) error
	// 可以Update到多个
	// MultiModify(ORMObject, ORMObject) error
}
