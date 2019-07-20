package jsonrpc_client

import (
	"errors"
	"io"
	"strings"

	bytes_pool "github.com/Myriad-Dreamin/go-ves/net/bytes_pool"
	request "github.com/Myriad-Dreamin/go-ves/net/request"
	"github.com/tidwall/gjson"
)

const (
	httpPrefix   = "http://"
	httpsPrefix  = "https://"
	maxBytesSize = 64 * 1024
)

func decorateHost(host string) string {
	if strings.HasPrefix(host, httpPrefix) || strings.HasPrefix(host, httpsPrefix) {
		return host
	}
	return "http://" + host
}

type JsonRpcClient struct {
	handler    *request.RequestClient
	bufferPool *bytes_pool.BytesPool
}

func NewJsonRpcClient(host string) *JsonRpcClient {
	return &JsonRpcClient{
		handler:    request.NewRequestClient(decorateHost(host)),
		bufferPool: bytes_pool.NewBytesPool(maxBytesSize),
	}
}

// todo: test invalid json
func (nc *JsonRpcClient) preloadJsonResponse(bb io.ReadCloser) ([]byte, error) {

	var b = nc.bufferPool.Get().([]byte)
	defer nc.bufferPool.Put(b)

	_, err := bb.Read(b)
	if err != nil && err != io.EOF {
		return nil, err
	}
	bb.Close()
	var jm = gjson.ParseBytes(b)
	if s := jm.Get("jsonrpc"); !s.Exists() || s.String() != "2.0" {
		return nil, errors.New("reject ret that is not jsonrpc: 2.0")
	}
	if s := jm.Get("error"); s.Exists() {
		return nil, FromGJsonResultError(s)
	}
	if s := jm.Get("result"); s.Exists() {
		if s.Index > 0 {
			return b[s.Index : s.Index+len(s.Raw)], nil
		}
	}
	return nil, errors.New("bad format of jsonrpc")
}

func (nc *JsonRpcClient) GetRequestParams(params map[string]interface{}) ([]byte, error) {
	b, err := nc.handler.GetWithParams(params)
	if err != nil {
		return nil, err
	}
	return nc.preloadJsonResponse(b)
}

func (nc *JsonRpcClient) PostRequestWithJsonObj(jsonObj interface{}) ([]byte, error) {
	b, err := nc.handler.PostWithJsonObj(jsonObj)
	if err != nil {
		return nil, err
	}
	return nc.preloadJsonResponse(b)
}