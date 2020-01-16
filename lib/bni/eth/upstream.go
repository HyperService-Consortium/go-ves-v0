package bni

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/HyperService-Consortium/go-uip/const/value_type"
	"github.com/HyperService-Consortium/go-uip/uiptypes"
	"github.com/Myriad-Dreamin/minimum-lib/sugar"
	"github.com/tidwall/gjson"
	"math/big"
	"reflect"
	"testing"
)

type mcs struct{}
type _serializer struct {
	Meta struct {
		Contract mcs
	}
}

var serializer = _serializer{}

func (mcs) Unmarshal(b []byte, meta *uiptypes.ContractInvokeMeta) error {
	return json.Unmarshal(b, meta)
}

func (mcs) Marshal(meta *uiptypes.ContractInvokeMeta) ([]byte, error) {
	return json.Marshal(meta)
}




type kv struct {
	k string
	v interface{}
}

type gJSONAssertion struct {
	kvs []kv
}

func (g gJSONAssertion) AssertBytes(object []byte) (err error) {
	for _, assertKeyValue := range g.kvs {
		k, v := assertKeyValue.k, assertKeyValue.v
		if err = g.compare(gjson.GetBytes(object, k), v); err != nil {
			return fmt.Errorf("compared failed on %v, assertion error %v", k, err)
		}
	}
	return
}

var int64T = reflect.TypeOf(int64(1))

func (g gJSONAssertion) compare(bytes gjson.Result, v interface{}) error {
	t := reflect.TypeOf(v)
	switch bytes.Type {
	case gjson.Null:
		if v != nil {
			return fmt.Errorf("compare failed: %v %v", bytes, v)
		}
	case gjson.False:
		if t.Kind() != reflect.Bool || v != false {
			return fmt.Errorf("compare failed: %v %v", bytes, v)
		}
	case gjson.True:
		if t.Kind() != reflect.Bool || v != true {
			return fmt.Errorf("compare failed: %v %v", bytes, v)
		}
	case gjson.Number:
		if !t.ConvertibleTo(int64T) ||
			reflect.ValueOf(v).Convert(int64T).Int() != bytes.Int() {
			return fmt.Errorf("compare failed: %v %v", bytes, v)
		}
	case gjson.String:
		if t.Kind() != reflect.String || v != bytes.String() {
			return fmt.Errorf("compare failed: %v %v", bytes, v)
		}
	case gjson.JSON:
		return fmt.Errorf("not basic comparable data")
	default:
		panic("unknown g-json type")
	}
	return nil
}

func gJSONWant(kvs ...kv) gJSONAssertion {
	return gJSONAssertion{kvs:kvs}
}

type mockBNIStorage struct {
	data []mockData
}

func (m mockBNIStorage) GetTransactionProof(chainID uiptypes.ChainID, blockID uiptypes.BlockID, color []byte) (uiptypes.MerkleProof, error) {
	panic("implement me")
}

func (m mockBNIStorage) GetStorageAt(chainID uiptypes.ChainID, typeID uiptypes.TypeID, contractAddress uiptypes.ContractAddress, pos []byte, description []byte) (uiptypes.Variable, error) {
	for _, d := range m.data {
		if d.chainID == chainID && d.typeID == typeID &&
			bytes.Equal(d.contractAddress, contractAddress) &&
			bytes.Equal(d.pos, pos) &&
			bytes.Equal(d.description, description) {
			return d.v, nil
		}
	}
	return nil, errors.New("not found")
}

func (m *mockBNIStorage) insertMockData(data []mockData) {
	m.data = append(m.data, data...)
}

type testFunc = func(t *testing.T)

type bNIStorageTestSet struct {
	s uiptypes.Storage
}

type mockData struct {
	chainID uiptypes.ChainID
	typeID uiptypes.TypeID
	contractAddress uiptypes.ContractAddress
	pos []byte
	description []byte
	v uiptypes.Variable
}

type mockValue struct {
	t value_type.Type
	v interface{}

}

func (m mockValue) GetType() uiptypes.TypeID {
	return m.t
}

func (m mockValue) GetValue() interface{} {
	return m.v
}

func (b bNIStorageTestSet) MockingData() []mockData {
	return []mockData{
		{
			chainID:         2,
			typeID:          value_type.Int64,
			contractAddress: make([]byte, 32),
			pos:             make([]byte, 2),
			description:     make([]byte, 2),
			v:               mockValue{value_type.Int64, int64(10)},
		},
		{
			chainID:         2,
			typeID:          value_type.Int32,
			contractAddress: make([]byte, 32),
			pos:             make([]byte, 2),
			description:     make([]byte, 2),
			v:               mockValue{value_type.Int32, int32(11)},
		},
		{
			chainID:         2,
			typeID:          value_type.Int128,
			contractAddress: make([]byte, 32),
			pos:             make([]byte, 2),
			description:     make([]byte, 2),
			v:               mockValue{value_type.Int128, bigInt3},
		},
		{
			chainID:         2,
			typeID:          value_type.Int256,
			contractAddress: make([]byte, 32),
			pos:             make([]byte, 2),
			description:     make([]byte, 2),
			v:               mockValue{value_type.Int256, bigInt3},
		},
	}
}

func (b bNIStorageTestSet) RunTests(t *testing.T) {
	t.Run("testGetInt32", b.testGetInt32)
	t.Run("testGetInt64", b.testGetInt64)
	t.Run("testGetInt128", b.testGetInt128)
	t.Run("testGetInt256", b.testGetInt256)
}

func assertType(l *testing.T, x uiptypes.Variable, t value_type.Type, k reflect.Kind) bool {
	l.Helper()
	if x.GetType() != t {
		l.Fatal("bad type")
		return false
	}
	v0 := x.GetValue()
	v := reflect.ValueOf(v0)
	if v.Type().Kind() != k {
		l.Fatal("bad value type")
		return false
	}
	return true
}

func assertTypeOf(l *testing.T, x uiptypes.Variable, t value_type.Type, r reflect.Type) bool {
	l.Helper()
	if x.GetType() != t {
		l.Fatal("bad type")
		return false
	}
	v0 := x.GetValue()
	v := reflect.ValueOf(v0)
	if v.Type() != r {
		l.Fatal("bad value type")
		return false
	}
	return true
}


func (b bNIStorageTestSet) testGetInt32(t *testing.T) {
	x := sugar.HandlerError(b.s.GetStorageAt(2, value_type.Int32, make([]byte, 32), make([]byte, 2), make([]byte, 2))).(uiptypes.Variable)
	if !assertType(t, x, value_type.Int32, reflect.Int32) {
		return
	}
	if x.GetValue().(int32) != 11 {
		t.Fatal("bad value")
	}
}



func (b bNIStorageTestSet) testGetInt64(t *testing.T) {
	x := sugar.HandlerError(b.s.GetStorageAt(2, value_type.Int64, make([]byte, 32), make([]byte, 2), make([]byte, 2))).(uiptypes.Variable)
	if !assertType(t, x, value_type.Int64, reflect.Int64) {
		return
	}
	if x.GetValue().(int64) != 10 {
		t.Fatal("bad value")
	}
}

var bigInt3 = big.NewInt(3)

func (b bNIStorageTestSet) testGetInt128(t *testing.T) {
	x := sugar.HandlerError(b.s.GetStorageAt(2, value_type.Int128, make([]byte, 32), make([]byte, 2), make([]byte, 2))).(uiptypes.Variable)
	if !assertTypeOf(t, x, value_type.Int128, reflect.TypeOf(bigInt3)) {
		return
	}
	if x.GetValue().(*big.Int).Cmp(bigInt3) != 0 {
		t.Fatal("bad value")
	}
}



func (b bNIStorageTestSet) testGetInt256(t *testing.T) {
	x := sugar.HandlerError(b.s.GetStorageAt(2, value_type.Int256, make([]byte, 32), make([]byte, 2), make([]byte, 2))).(uiptypes.Variable)
	if !assertTypeOf(t, x, value_type.Int256, reflect.TypeOf(bigInt3)) {
		return
	}
	if x.GetValue().(*big.Int).Cmp(bigInt3) != 0 {
		t.Fatal("bad value")
	}
}


