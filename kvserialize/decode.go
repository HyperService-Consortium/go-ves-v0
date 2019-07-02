package kvs

import (
	"fmt"
	"reflect"

	"github.com/Myriad-Dreamin/go-ves/types"
)

func isBlank(value reflect.Value) bool {
	switch value.Kind() {
	case reflect.String:
		return value.Len() == 0
	case reflect.Bool:
		return !value.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return value.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return value.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return value.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return value.IsNil()
	}
	return reflect.DeepEqual(value.Interface(), reflect.Zero(value.Type()).Interface())
}

func SetValue(kvv interface{}, segValue reflect.Value) error {
	genuineSegValue := segValue
	if genuineSegValue.Kind() == reflect.Ptr && !genuineSegValue.Elem().CanSet() {
		fmt.Println("xxx")
		return fmt.Errorf("bad interface can't be set")
	} else if genuineSegValue.Kind() == reflect.Ptr {
		genuineSegValue = genuineSegValue.Elem()
	}

	switch kvvValue := kvv.(type) {
	// TODO: SliceType and ArrayType
	case string:
		if genuineSegValue.Kind() == reflect.String {
			genuineSegValue.SetString(kvvValue)
			return nil
		} else if segValue, ok := segValue.Interface().(VString); ok {
			segValue.UnmarshalKV(kvvValue)
			return nil
		}
		return fmt.Errorf("cannot convert string to segValue")
	case []byte:
		// TODO: Check this judge
		// var typeOfBytes = reflect.TypeOf([]byte(nil))
		if genuineSegValue.Kind() == reflect.Array && genuineSegValue.Elem().Kind() == reflect.Uint8 {
			genuineSegValue.SetBytes(kvvValue)
			return nil
		} else if segValue, ok := segValue.Interface().(VBytes); ok {
			segValue.UnmarshalKV(kvvValue)
			return nil
		}
		return fmt.Errorf("cannot convert []byte to segValue")
	case bool:
		if genuineSegValue.Kind() == reflect.Bool {
			genuineSegValue.SetBool(kvvValue)
			return nil
		} else if segValue, ok := segValue.Interface().(VBool); ok {
			segValue.UnmarshalKV(kvvValue)
			return nil
		}
		return fmt.Errorf("cannot convert bool to segValue")
	case int64:
		if genuineSegValue.Kind() == reflect.Int64 {
			genuineSegValue.SetInt(MaxIntType(kvvValue))
			return nil
		} else if segValue, ok := segValue.Interface().(VInt64); ok {
			segValue.UnmarshalKV(MaxIntType(kvvValue))
			return nil
		}
		return fmt.Errorf("cannot convert int64 to segValue")
	case int32:
		if genuineSegValue.Kind() == reflect.Int32 {
			genuineSegValue.SetInt(MaxIntType(kvvValue))
			return nil
		} else if segValue, ok := segValue.Interface().(VInt32); ok {
			segValue.UnmarshalKV(kvvValue)
			return nil
		}
		return fmt.Errorf("cannot convert int32 to segValue")
	case int16:
		if genuineSegValue.Kind() == reflect.Int16 {
			genuineSegValue.SetInt(MaxIntType(kvvValue))
			return nil
		} else if segValue, ok := segValue.Interface().(VInt64); ok {
			segValue.UnmarshalKV(MaxIntType(kvvValue))
			return nil
		}
		return fmt.Errorf("cannot convert int16 to segValue")
	case int8:
		if genuineSegValue.Kind() == reflect.Int8 {
			genuineSegValue.SetInt(MaxIntType(kvvValue))
			return nil
		} else if segValue, ok := segValue.Interface().(VInt8); ok {
			segValue.UnmarshalKV(kvvValue)
			return nil
		}
		return fmt.Errorf("cannot convert int8 to segValue")
	case uint64:
		if genuineSegValue.Kind() == reflect.Uint64 {
			genuineSegValue.SetUint(MaxUintType(kvvValue))
			return nil
		} else if segValue, ok := segValue.Interface().(VUint64); ok {
			segValue.UnmarshalKV(MaxUintType(kvvValue))
			return nil
		}
		return fmt.Errorf("cannot convert uint64 to segValue")
	case uint32:
		if genuineSegValue.Kind() == reflect.Uint32 {
			genuineSegValue.SetUint(MaxUintType(kvvValue))
			return nil
		} else if segValue, ok := segValue.Interface().(VUint32); ok {
			segValue.UnmarshalKV(kvvValue)
			return nil
		}
		return fmt.Errorf("cannot convert uint32 to segValue")
	case uint16:
		if genuineSegValue.Kind() == reflect.Uint16 {
			genuineSegValue.SetUint(MaxUintType(kvvValue))
			return nil
		} else if segValue, ok := segValue.Interface().(VUint64); ok {
			segValue.UnmarshalKV(MaxUintType(kvvValue))
			return nil
		}
		return fmt.Errorf("cannot convert uint16 to segValue")
	case uint8:
		if genuineSegValue.Kind() == reflect.Uint8 {
			genuineSegValue.SetUint(MaxUintType(kvvValue))
			return nil
		} else if segValue, ok := segValue.Interface().(VUint8); ok {
			segValue.UnmarshalKV(kvvValue)
			return nil
		}
		return fmt.Errorf("cannot convert uint8 to segValue")
	default:
		return fmt.Errorf("v of kv is not basic type")
	}
}

func Unmarshal(kvs []types.KVPair, xx interface{}) error {
	objt := reflect.TypeOf(xx)
	if objt.Kind() != reflect.Ptr {
		fmt.Println("vx")
		return fmt.Errorf("bad interface can't be set")
	} else {
		objt = objt.Elem()
	}

	kvmp := make(map[string]interface{})
	for _, kv := range kvs {
		kvmp[kv.Key] = kv.Value
	}
	fmt.Println(objt.NumField())
	obj := reflect.ValueOf(xx).Elem()
	for idx := 0; idx < objt.NumField(); idx++ {
		if tg, ok := objt.Field(idx).Tag.Lookup("kv"); ok {
			if v, ok := kvmp[tg]; ok {
				f := obj.Field(idx)
				SetValue(v, f)
			} else {
				return fmt.Errorf("kv tag set but missing kvpair")
			}
		}

	}
	return nil
}
