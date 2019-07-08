package nsbcli

import (
	"encoding/json"
	"fmt"

	"github.com/tidwall/gjson"
)

type jsonMap = map[string]interface{}

type JsonError struct {
	errorx string
}

func (je JsonError) Error() string {
	return je.errorx
}

func fromJsonMapError(jm jsonMap) *JsonError {
	return &JsonError{
		errorx: fmt.Sprintf("jsonrpc error: %v(%v), %v", jm["message"], jm["code"], jm["data"]),
	}
}

func fromBytesError(b []byte) *JsonError {
	var jm jsonMap
	err := json.Unmarshal(b, &jm)
	if err != nil {
		return &JsonError{
			errorx: fmt.Sprintf("bad format of json error: %v", err),
		}
	}
	return fromJsonMapError(jm)
}

func fromGJsonResultError(b gjson.Result) *JsonError {
	return &JsonError{
		errorx: fmt.Sprintf("jsonrpc error: %v(%v), %v", b.Get("message"), b.Get("code"), b.Get("data")),
	}
}
