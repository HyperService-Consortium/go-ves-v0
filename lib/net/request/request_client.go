package request

import (
	"net/http"
	"time"

	"github.com/imroc/req"
)

var (
	request_timeout = 50 * time.Second
)

type RequestClient struct {
	BaseURL string
	Header  Header
}

func NewRequestClient(url string) *RequestClient {
	return &RequestClient{BaseURL: url}
}

func (jc *RequestClient) SetHeader(strMap map[string]string) *RequestClient {
	jc.Header = Header(strMap)
	return jc
}

func (jc *RequestClient) SetHeaderWithReqHeader(header Header) *RequestClient {
	jc.Header = header
	return jc
}

func (jc *RequestClient) Group(sub string) *RequestClient {
	return &RequestClient{BaseURL: jc.BaseURL + sub, Header: jc.Header}
}

type RequestClientX struct {
	BaseURL string
	Header  Header
	path    string
}

func NewRequestClientX(url string) *RequestClientX {
	return &RequestClientX{BaseURL: url}
}

func (jc *RequestClientX) SetHeader(i interface{}) *RequestClientX {
	switch s := i.(type) {
	case map[string]string:
		jc.Header = Header(s)
	case Header:
		jc.Header = s
	default:
	}
	return jc
}

func (jc *RequestClientX) Group(sub string) *RequestClientX {
	return &RequestClientX{BaseURL: jc.BaseURL + sub, Header: jc.Header}
}

func (jc *RequestClientX) Path(path string) *RequestClientX {
	jc.path = path
	return jc
}

func (jc *RequestClientX) Use(handler func(*Resp) error) *Context {
	return &Context{BaseURL: jc.BaseURL + jc.path, Header: jc.Header, handler: handler}
}

type Context struct {
	BaseURL string
	Header  Header
	handler func(*Resp) error
}

func (jc *Context) Path(path string) *Context {
	jc.BaseURL += path
	return jc
}

func SetConnPool() {
	client := &http.Client{}
	client.Transport = &http.Transport{
		MaxIdleConnsPerHost: 500,
	}

	req.SetClient(client)
	req.SetTimeout(request_timeout)
}

func init() {
	SetConnPool()
}
