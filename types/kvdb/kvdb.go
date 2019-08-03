package kvdb

import (
	"bytes"

	uiptypes "github.com/Myriad-Dreamin/go-uip/types"

	types "github.com/Myriad-Dreamin/go-ves/types"
)

type Database struct {
}

type Getter struct {
	idb         types.Index
	isc_address []byte
}

type Setter struct {
	idb         types.Index
	isc_address []byte
}

func clone(b []byte) []byte {
	var c = make([]byte, len(b))
	copy(c, b)
	return c
}

func cloneWithLen(b []byte, l int) []byte {
	var c = make([]byte, l)
	copy(c, b)
	return c
}

func decorate(p, b []byte) []byte {
	var x = bytes.NewBuffer(make([]byte, len(p)+len(b)))
	x.Reset()
	x.Write(p)
	x.WriteByte('k')
	x.WriteByte('k')
	x.Write(b)
	return x.Bytes()
}

func (g *Getter) Get(k []byte) ([]byte, error) {
	return g.idb.Get(decorate(g.isc_address, k))
}

func (g *Setter) Set(k, v []byte) error {
	return g.idb.Set(decorate(g.isc_address, k), v)
}

func (g *Database) GetKV(idb types.Index, isc_address, b []byte) ([]byte, error) {
	return idb.Get(decorate(isc_address, b))
}

func (g *Database) SetKV(idb types.Index, isc_address, k, v []byte) error {
	return idb.Set(decorate(isc_address, k), v)
}

func (db *Database) GetGetter(idb types.Index, isc_address []byte) uiptypes.KVGetter {
	return &Getter{idb, isc_address}
}
func (db *Database) GetSetter(idb types.Index, isc_address []byte) types.KVSetter {
	return &Setter{idb, isc_address}
}
