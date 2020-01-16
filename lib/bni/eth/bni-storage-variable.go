package bni

import (
	"github.com/HyperService-Consortium/go-uip/uiptypes"
)

type variable struct {
	Type  uiptypes.TypeID
	Value interface{}
}

func (v variable) GetType() uiptypes.TypeID {
	return v.Type
}

func (v variable) GetValue() interface{} {
	return v.Value
}
