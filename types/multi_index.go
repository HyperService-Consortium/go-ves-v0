package types

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
	Insert(ORMObject) error

	Select(ORMObject, []ORMObject) error

	// 要求只Delete到一个
	Delete(ORMObject) error
	// 可以Delete多个
	MultiDelete(ORMObject) error

	// 要求只Update到一个
	Modify(ORMObject, ORMObject) error
	// 可以Update到多个
	MultiModify(ORMObject, ORMObject) error
}

type MultiIndex interface {
	/*
	 * interface{}: []KVPair or ORMObject
	 */
	Insert(interface{}) error

	/*
	 * interface{} A: []KVPair or ORMObject
	 * []interface{} B: []ORMObject
	 */
	Select(interface{}, []interface{}) error

	/*
	 * interface{}: []KVPair or ORMObject
	 */
	// 要求只Delete到一个
	Delete(interface{}) error
	// 可以Delete多个
	MultiDelete(interface{}) error

	/*
	 * interface{}: []KVPair or ORMObject
	 */
	// 要求只Update到一个
	Modify(interface{}, interface{}) error
	// 可以Update到多个
	MultiModify(interface{}, interface{}) error
}
