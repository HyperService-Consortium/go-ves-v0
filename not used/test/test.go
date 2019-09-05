package main

import (
	"fmt"
	"reflect"

	kvs "github.com/HyperService-Consortium/go-ves/kvserialize"
	"github.com/HyperService-Consortium/go-ves/types"
)

type OAO struct {
	MyHH  *ww   `kv:"FirstElement" omitempty:"true"`
	MyHH2 int16 `kv:"SecondElement" omitempty:"true"`
}

type ww struct {
	A int8
	B int8
}

func (c *ww) String() string {
	return fmt.Sprintf("%v", *c)
}

func (w *ww) MarshalKV() interface{} {
	return int16(w.A)<<8 | int16(w.B)
}

func (w *ww) UnmarshalKV(cw int16) {
	w.A = int8((cw >> 8) & 0xff)
	w.B = int8(cw & 0xff)
}

func RawKVSet(kvv interface{}, xx interface{}) error {
	v := reflect.ValueOf(xx)
	fmt.Println(v.Kind() == reflect.Ptr, v)
	if v.Kind() == reflect.Ptr && !v.Elem().CanSet() {
		fmt.Println("xxx")
		return fmt.Errorf("bad interface can't be set")
	} else if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	switch c := kvv.(type) {
	case int16:
		if v.Kind() == reflect.Int16 {
			v.SetInt(kvs.MaxIntType(c))
			return nil
		} else if xx, ok := xx.(kvs.VInt16); ok {
			xx.UnmarshalKV(c)
			return nil
		}
		return fmt.Errorf("cannot convert int16 to v interface")
	default:
		return fmt.Errorf("v of kv is not basic type")
	}
}

func main() {
	var x ww
	fmt.Println("magic string method: ", &x, x.String())

	fmt.Println("marshal kv", (&ww{A: 0x3, B: 0x4}).MarshalKV(), 0x3<<8|0x4)

	var cw = int16(0x33ff)
	var y = new(ww)

	y.UnmarshalKV(cw)
	fmt.Printf("ww's unmarshal method %+v\n", *y)

	fmt.Println("set kv", "error:", RawKVSet(int16(33<<8|22), y), "content:", *y)
	fmt.Println("set kv", "error:", RawKVSet(int16(22), &cw), "content:", cw)

	var wwww = &OAO{new(ww), 1}
	fmt.Println("unmarshal kv", "error:", kvs.Unmarshal(
		[]types.KVPair{
			types.KVPair{Key: "FirstElement", Value: int16(22)},
			types.KVPair{Key: "SecondElement", Value: int16(33<<8 | 22)},
		},
		wwww), wwww,
	)
	wwww = &OAO{MyHH2: 1}
	fmt.Println(kvs.Marshal(wwww))

	wwww = &OAO{MyHH: &ww{A: 127, B: 0}}
	fmt.Println(kvs.Marshal(wwww))
}
