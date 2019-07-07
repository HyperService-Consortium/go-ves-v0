package request

import (
	"fmt"
	"testing"

	mtest "github.com/Myriad-Dreamin/mydrest"
	"github.com/imroc/req"
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

type service struct {
	b3 []byte
}

func (s *service) myController(r *req.Resp) (err error) {
	s.b3, err = r.ToBytes()
	return
}

func TestGetParamAndGroup(t *testing.T) {
	var params = &mParam3{"closed"}
	githubApi := NewRequestClient("https://api.github.com")
	b, err := githubApi.Group("/repos/solomonxie/solomonxie.github.io/issues").GetWithStruct(params)
	s.AssertNoErr(t, err)

	reposApi := githubApi.Group("/repos")

	b2, err := reposApi.Group("/solomonxie/solomonxie.github.io/issues").GetWithStruct(params)
	s.AssertNoErr(t, err)
	s.AssertEqual(t, string(b), string(b2))

	solomonxieRepo := reposApi.Group("/solomonxie/solomonxie.github.io")

	b2, err = solomonxieRepo.Group("/issues").GetWithStruct(params)
	s.AssertNoErr(t, err)
	s.AssertEqual(t, string(b), string(b2))

	githubApiX := NewRequestClientX("https://api.github.com")
	b2, err = githubApiX.Group("/repos/solomonxie/solomonxie.github.io/issues").Get(params)
	s.AssertNoErr(t, err)
	s.AssertEqual(t, string(b), string(b2))

	b2, err = githubApiX.Group("/repos/solomonxie/solomonxie.github.io/issues").Get(&req.QueryParam{
		"status": "closed",
	})
	s.AssertNoErr(t, err)
	s.AssertEqual(t, string(b), string(b2))
	// req.Debug = true
	yandeApi := NewRequestClientX("https://yande.re")

	_, err = yandeApi.Group("/post").Get(&req.QueryParam{
		"tags": "dress",
	})
	s.AssertNoErr(t, err)
	serve := new(service)
	err = yandeApi.Group("/post").Use(serve.myController).Get(&req.Param{
		"tags": "asian_clothes",
	})
	s.AssertNoErr(t, err)
	fmt.Println(string(serve.b3))
	// s.AssertEqual(t, string(b), string(b2))
}
