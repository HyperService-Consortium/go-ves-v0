package request

import (
	"testing"

	mtest "github.com/Myriad-Dreamin/mydrest"
)

var s mtest.TestHelper

func TestGet(t *testing.T) {
	_, err := NewRequestClient("http://www.baidu.com").Get()
	s.AssertNoErr(t, err)
}

type mParam struct {
	Keyword    string `url:"keyword"`
	FromSource string `url:"from_source"`
}

type mParam2 struct {
	Q    string `url:"q"`
	Type string `url:"type"`
}

type mParam3 struct {
	State string `url:"state"`
}

type mParam4 struct {
}

type service struct {
	b3 []byte
}

func (s *service) myController(r *Resp) (err error) {
	s.b3, err = r.ToBytes()
	return
}

func TestGetParamAndGroup(t *testing.T) {
	var params = &mParam4{}
	NSBApi := NewRequestClient("http://47.251.2.73:26657")
	b, err := NSBApi.Group("/abci_info").GetWithStruct(params)
	s.AssertNoErr(t, err)

	abciInfoApi := NSBApi.Group("/abci_info")

	b2, err := abciInfoApi.GetWithStruct(params)
	s.AssertNoErr(t, err)
	s.AssertEqual(t, string(b), string(b2))

	NSBApiX := NewRequestClientX("http://47.251.2.73:26657")
	b2, err = NSBApiX.Group("/abci_info").Get(params)
	s.AssertNoErr(t, err)
	s.AssertEqual(t, string(b), string(b2))

	b2, err = NSBApiX.Group("/abci_info").Get(&QueryParam{})
	s.AssertNoErr(t, err)
	s.AssertEqual(t, string(b), string(b2))
	// req.Debug = true
	yandeApi := NewRequestClientX("https://yande.re")

	_, err = yandeApi.Group("/post").Get(&QueryParam{
		"tags": "dress",
	})
	s.AssertNoErr(t, err)
	serve := new(service)
	err = yandeApi.Group("/post").Use(serve.myController).Get(&Param{
		"tags": "dress",
	})
	s.AssertNoErr(t, err)
	// s.AssertEqual(t, string(b), string(b2))
}
