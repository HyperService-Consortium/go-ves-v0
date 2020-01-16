package bni

import (
	"encoding/json"
	"fmt"
	"github.com/HyperService-Consortium/go-uip/uiptypes"
	"github.com/tidwall/gjson"
	"reflect"
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
