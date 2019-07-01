package types

type KVPair interface {
	Key() Stringable
	Value() Stringable
}
