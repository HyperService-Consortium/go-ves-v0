package types

type MultiIndex interface {
	Insert([]KVPairs) error
	Delete([]KVPairs) error
	Select([]KVPairs) error
	Modify([]KVPairs, []KVPairs) error
}
