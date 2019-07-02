package types

type ORMObject interface {
	Fetch(MultiIndex, ...KVPair) error
	Update(MultiIndex) error
	Delete(MultiIndex) error
}
