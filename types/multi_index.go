package types

type MultiIndex interface {
	Insert([]KVPair) error
	Delete([]KVPair) error
	Select([]KVPair) ([]KVPair, error)
	Modify([]KVPair, []KVPair) error
}
