package kvs

import (
	"fmt"
	"reflect"

	"github.com/HyperService-Consortium/go-ves/types"
)

func GetValue(rawValue interface{}) (interface{}, error) {
	switch rawValue := rawValue.(type) {
	case bool, string, byte, []byte, float32, float64,
		int, int8, int16, int32, int64,
		uint, uint16, uint32, uint64: /*uint8 == byte*/
		return rawValue, nil
	case Marshalable:
		return rawValue.MarshalKV(), nil
	default:
		return nil, fmt.Errorf("interface to get is not kv-marshalable type")
	}
}

func Marshal(xx interface{}) (kvs []types.KVPair, err error) {
	var amt interface{}

	t := reflect.TypeOf(xx)
	if t.Kind() != reflect.Ptr {
		fmt.Println("vx")
		return nil, fmt.Errorf("bad interface can't be set")
	} else {
		t = t.Elem()
	}

	s := reflect.ValueOf(xx).Elem()
	for idx := 0; idx < t.NumField(); idx++ {
		ele := s.Field(idx)
		if _, ok := t.Field(idx).Tag.Lookup("omitempty"); ok && isBlank(ele) {
			continue
		}
		if tg, ok := t.Field(idx).Tag.Lookup("kv"); ok {

			amt, err = GetValue(ele.Interface())
			if err != nil {
				return nil, err
			}
			kvs = append(kvs, types.KVPair{Key: tg, Value: amt})
		}
	}
	return
}
