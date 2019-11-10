package types

type SessionKVBase interface {
	SetKV(Index, isc_address, provedKey, provedValue) error
	GetKV(Index, isc_address, provedKey) (provedValue, error)
	GetSetter(Index, isc_address) KVSetter
	GetGetter(Index, isc_address) KVGetter
}

type provedKey = []byte
type provedValue = []byte

type KVSetter interface {
	Set(provedKey, provedValue) error
}

type KVGetter interface {
	Get([]byte) ([]byte, error)
}
