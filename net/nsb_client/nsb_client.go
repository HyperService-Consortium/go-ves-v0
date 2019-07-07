package nsbcli

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	request "github.com/Myriad-Dreamin/go-ves/net/request"
)

const (
	httpPrefix  = "http://"
	httpsPrefix = "https://"
)

type NSBClient struct {
	handler *request.RequestClient
}

func decorateHost(host string) string {
	if strings.HasPrefix(host, httpPrefix) || strings.HasPrefix(host, httpsPrefix) {
		return host
	}
	return "http://" + host
}

func NewNSBClient(host string) *NSBClient {
	return &NSBClient{
		handler: request.NewRequestClient(decorateHost(host)),
	}
}

type jsonMap = map[string]interface{}

type preLoadJsonStruct struct {
	JSONVersion string                 `json:"jsonrpc"`
	ID          string                 `json:"id"`
	Error       map[string]interface{} `json:"error"`
	Result      interface{}            `json:"result"`
}

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

func preloadJson(b []byte) ([]byte, error) {
	var jm preLoadJsonStruct
	if err := json.Unmarshal(b, &jm); err == nil {
		if jm.JSONVersion != "2.0" {
			return nil, errors.New("reject ret that is not jsonrpc: 2.0")
		}
		if jm.Error != nil {
			return nil, fromJsonMapError(jm.Error)
		}
		if jm.Result != nil {
			return json.Marshal(jm.Result)
		}
	} else {
		return nil, err
	}
	return nil, errors.New("bad format of jsonrpc")
}

func PretiJson(minterface interface{}) string {
	je, _ := json.MarshalIndent(minterface, "", "\t")
	return string(je)
}

type AbciInfoResponse struct {
	Data       string `json:"data"`
	Version    string `json:"version"`
	AppVersion string `json:"app_version"`
}
type AbciInfo struct {
	Response *AbciInfoResponse `json:"response"`
}

func (nc *NSBClient) GetAbciInfo() (*AbciInfoResponse, error) {
	b, err := nc.handler.Group("/abci_info").GetWithParams(request.Param{})
	if err != nil {
		return nil, err
	}
	b, err = preloadJson(b)
	if err != nil {
		return nil, err
	}
	var a AbciInfo
	err = json.Unmarshal(b, &a)
	if err != nil {
		return nil, err
	}
	return a.Response, nil
}

func (nc *NSBClient) GetBlock(id int64) (*AbciInfo, error) {
	b, err := nc.handler.Group("/block").GetWithParams(request.Param{
		"height": id,
	})
	if err != nil {
		return nil, err
	}
	var i interface{}
	i, err = preloadJson(b)
	if err != nil {
		return nil, err
	}
	fmt.Println(i, err)
	return nil, nil
}
