package types

type KVPairs interface {
	Key() Stringable
	Value() Stringable
}
