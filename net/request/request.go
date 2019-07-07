package request

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/google/go-querystring/query"
	"github.com/imroc/req"
)

var (
	request_timeout = 5 * time.Second
)

type RequestClient struct {
	BaseURL string
	Header  req.Header
}

func NewRequestClient(url string) *RequestClient {
	return &RequestClient{BaseURL: url}
}

func (jc *RequestClient) SetHeader(strMap map[string]string) *RequestClient {
	jc.Header = req.Header(strMap)
	return jc
}

func (jc *RequestClient) SetHeaderWithReqHeader(header req.Header) *RequestClient {
	jc.Header = header
	return jc
}

func (jc *RequestClient) Group(sub string) *RequestClient {
	return &RequestClient{BaseURL: jc.BaseURL + sub, Header: jc.Header}
}

func get(s string, v ...interface{}) ([]byte, error) {
	resp, err := req.Get(s, v...)
	if err != nil {
		return nil, err
	}
	var sc = resp.Response().StatusCode
	if sc != 200 {
		return nil, fmt.Errorf("error code: %v", sc)
	}
	return resp.ToBytes()
}

func (jc *RequestClient) Get() ([]byte, error) {
	return get(jc.BaseURL, jc.Header)
}

func (jc *RequestClient) GetWithKVMap(request map[string]interface{}) ([]byte, error) {
	return get(jc.BaseURL, jc.Header, request)
}

func (jc *RequestClient) GetWithParams(params ...interface{}) ([]byte, error) {
	return get(jc.BaseURL, jc.Header, params)
}

func (jc *RequestClient) GetWithStruct(request interface{}) ([]byte, error) {
	v, err := query.Values(request)
	if err != nil {
		return nil, err
	}
	s := bytes.NewBufferString(jc.BaseURL)
	err = s.WriteByte('?')
	if err != nil {
		return nil, err
	}
	_, err = s.WriteString(v.Encode())
	if err != nil {
		return nil, err
	}
	return get(s.String(), jc.Header)
}

type RequestClientX struct {
	BaseURL string
	Header  req.Header
	path    string
}

func NewRequestClientX(url string) *RequestClientX {
	return &RequestClientX{BaseURL: url}
}

func getx(s string, v ...interface{}) (*req.Resp, error) {
	resp, err := req.Get(s, v...)
	if err != nil {
		return nil, err
	}
	var sc = resp.Response().StatusCode
	if sc != 200 {
		return nil, fmt.Errorf("error code: %v", sc)
	}
	return resp, err
}

func (jc *RequestClientX) SetHeader(i interface{}) *RequestClientX {
	switch s := i.(type) {
	case map[string]string:
		jc.Header = req.Header(s)
	case req.Header:
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

func (jc *RequestClientX) Get(params ...interface{}) ([]byte, error) {
	fin, r := false, jc.BaseURL
	for idx, param := range params {
		if reflect.TypeOf(param).Kind() == reflect.Struct {
			v, err := query.Values(param)
			if err != nil {
				return nil, err
			}
			s := bytes.NewBufferString(r)
			err = s.WriteByte('?')
			if err != nil {
				return nil, err
			}
			_, err = s.WriteString(v.Encode())
			if err != nil {
				return nil, err
			}
			fin = true
			r = s.String()
			params = append(params[:idx], params[idx+1:]...)
			continue
		}
		if reflect.TypeOf(param).Kind() == reflect.Ptr && reflect.ValueOf(param).Elem().Kind() == reflect.Struct {
			v, err := query.Values(param)
			if err != nil {
				return nil, err
			}
			s := bytes.NewBufferString(r)
			err = s.WriteByte('?')
			if err != nil {
				return nil, err
			}
			_, err = s.WriteString(v.Encode())
			if err != nil {
				return nil, err
			}
			fin = true
			r = s.String()
			params = append(params[:idx], params[idx+1:]...)
			continue
		}

		switch i := param.(type) {
		case *req.QueryParam:
			params[idx] = *i
			if fin {
				return nil, errors.New("struct mapping and map[string]interface{} cannot be used at the same time")
			}
		case *req.Param:
			params[idx] = *i
			if fin {
				return nil, errors.New("struct mapping and map[string]interface{} cannot be used at the same time")
			}
		case req.Param, req.QueryParam:
			if fin {
				return nil, errors.New("struct mapping and map[string]interface{} cannot be used at the same time")
			}
		default:
		}
	}
	params = append(params, jc.Header)
	return get(r, params...)
}

func (jc *RequestClientX) Use(handler func(*req.Resp) error) *Context {
	return &Context{BaseURL: jc.BaseURL + jc.path, Header: jc.Header, handler: handler}
}

type Context struct {
	BaseURL string
	Header  req.Header
	handler func(*req.Resp) error
}

func (jc *Context) Path(path string) *Context {
	jc.BaseURL += path
	return jc
}

func getc(s string, handler func(*req.Resp) error, v ...interface{}) error {
	resp, err := req.Get(s, v...)
	if err != nil {
		return err
	}
	var sc = resp.Response().StatusCode
	if sc != 200 {
		return fmt.Errorf("error code: %v", sc)
	}

	err = handler(resp)
	return err
}

func (jc *Context) Get(params ...interface{}) error {
	fin, r := false, jc.BaseURL
	for idx, param := range params {
		if reflect.TypeOf(param).Kind() == reflect.Struct {
			v, err := query.Values(param)
			if err != nil {
				return err
			}
			s := bytes.NewBufferString(r)
			err = s.WriteByte('?')
			if err != nil {
				return err
			}
			_, err = s.WriteString(v.Encode())
			if err != nil {
				return err
			}
			fin = true
			r = s.String()
			params = append(params[:idx], params[idx+1:]...)
			continue
		}
		if reflect.TypeOf(param).Kind() == reflect.Ptr && reflect.ValueOf(param).Elem().Kind() == reflect.Struct {
			v, err := query.Values(param)
			if err != nil {
				return err
			}
			s := bytes.NewBufferString(r)
			err = s.WriteByte('?')
			if err != nil {
				return err
			}
			_, err = s.WriteString(v.Encode())
			if err != nil {
				return err
			}
			fin = true
			r = s.String()
			params = append(params[:idx], params[idx+1:]...)
			continue
		}

		switch i := param.(type) {
		case *req.QueryParam:
			params[idx] = *i
			if fin {
				return errors.New("struct mapping and map[string]interface{} cannot be used at the same time")
			}
		case *req.Param:
			params[idx] = *i
			if fin {
				return errors.New("struct mapping and map[string]interface{} cannot be used at the same time")
			}
		case req.Param, req.QueryParam:
			if fin {
				return errors.New("struct mapping and map[string]interface{} cannot be used at the same time")
			}
		default:
		}
	}
	params = append(params, jc.Header)
	return getc(r, jc.handler, params...)
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
