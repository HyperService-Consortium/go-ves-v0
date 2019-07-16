package request

import (
	"fmt"
	"io"

	"github.com/imroc/req"
)

func post(s string, v ...interface{}) (io.ReadCloser, error) {
	resp, err := req.Post(s, v...)
	if err != nil {
		return nil, err
	}
	var sc = resp.Response().StatusCode
	if sc != 200 {
		return nil, fmt.Errorf("error code: %v", sc)
	}
	return resp.Response().Body, nil
}

func PostWithBody() {

}

func (jc *RequestClient) PostWithJsonObj(obj interface{}) (io.ReadCloser, error) {
	return post(jc.BaseURL, jc.Header, req.BodyJSON(obj))
}
