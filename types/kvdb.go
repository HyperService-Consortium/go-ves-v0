package types

import uiptypes "github.com/HyperService-Consortium/go-uip/uiptypes"

type SessionKVBase interface {
	SetKV(Index, isc_address, provedKey, provedValue) error
	GetKV(Index, isc_address, provedKey) (provedValue, error)
	GetSetter(Index, isc_address) KVSetter
	GetGetter(Index, isc_address) uiptypes.KVGetter
}

type provedKey = []byte
type provedValue = []byte

type KVSetter interface {
	Set(provedKey, provedValue) error
}
