package user

import (
	"fmt"
	"testing"

	xorm_multi_index "github.com/HyperService-Consortium/go-ves/database/multi_index"
	mtest "github.com/HyperService-Consortium/mydrest"
)

type TestHelper struct {
	mtest.TestHelper
	res   *xorm_multi_index.XORMMultiIndexImpl
	logic *XORMUserBase
}

var s TestHelper

const path = "ves:123456@tcp(127.0.0.1:3306)/ves?charset=utf8"

type mA struct {
	x uint64
	y []byte
}

func (a *mA) GetChainId() uint64 { return a.x }
func (a *mA) GetAddress() []byte { return a.y }

func SetUpHelper() {
	var err error
	s.res, err = xorm_multi_index.GetXORMMultiIndex("mysql", path)
	s.OutAssertNoErr(err)
	err = s.res.Register(&XORMUserAdapter{})
	fmt.Println(err)
	s.OutAssertNoErr(err)
	s.logic = new(XORMUserBase)
}

func TestMain(m *testing.M) {
	SetUpHelper()
	m.Run()
}

func TestGetDB(t *testing.T) {
	_, err := xorm_multi_index.GetXORMMultiIndex("mysql", path)
	s.AssertNoErr(t, err)
}

func TestRaw(t *testing.T) {
	var err error
	tt1 := &mA{x: 0x20202020, y: []byte{3, 4, 5}}
	tt2 := &mA{x: 0x02020202, y: []byte{3, 5, 4}}
	ww1 := NewXORMUserAdapter("xxx", tt1)
	ww2 := NewXORMUserAdapter("xxx", tt2)
	err = s.res.Insert(ww1)
	s.AssertNoErr(t, err)
	err = s.res.Insert(ww2)
	s.AssertNoErr(t, err)
	fmt.Println(s.res.SelectAll(ww1))
	err = s.res.Delete(ww1)
	s.AssertNoErr(t, err)
	err = s.res.Delete(ww2)
	s.AssertNoErr(t, err)
	fmt.Println(s.res.SelectAll(ww1))
}

func TestLogic(t *testing.T) {
	var err error
	tt1 := &mA{x: 0x20202020, y: []byte{3, 4, 5}}
	tt2 := &mA{x: 0x02020202, y: []byte{3, 5, 4}}
	ww1 := NewXORMUserAdapter("xxx", tt1)
	ww2 := NewXORMUserAdapter("xxx", tt2)
	fmt.Println(ww1)
	fmt.Println(ww2)
	err = s.logic.InsertAccount(s.res, "xxx", tt1)
	s.AssertNoErr(t, err)
	err = s.logic.InsertAccount(s.res, "xxx", tt2)
	s.AssertNoErr(t, err)
	fmt.Println(s.logic.FindUser(s.res, "xxx"))

	tt3 := &mA{x: 0x20202020, y: []byte{3, 4, 6}}
	tt4 := &mA{x: 0x02020202, y: []byte{3, 5, 6}}
	ww3 := NewXORMUserAdapter("xxx", tt3)
	ww4 := NewXORMUserAdapter("yyy", tt4)
	err = s.logic.InsertAccount(s.res, "xxx", tt3)
	s.AssertNoErr(t, err)
	err = s.logic.InsertAccount(s.res, "yyy", tt4)
	s.AssertNoErr(t, err)

	fmt.Println(s.logic.FindAccounts(s.res, "xxx", 0x20202020))

	fmt.Println(s.logic.FindUser(s.res, "xxx"))

	fmt.Println(s.logic.InvertFind(s.res, ww2))

	fmt.Println(s.logic.InvertFind(s.res, ww4))

	err = s.res.Delete(ww1)
	s.AssertNoErr(t, err)
	err = s.res.Delete(ww2)
	s.AssertNoErr(t, err)

	err = s.res.Delete(ww3)
	s.AssertNoErr(t, err)
	err = s.res.Delete(ww4)
	s.AssertNoErr(t, err)
}

func TestMoreLogic(t *testing.T) {
	// tt1 := randomSession()
	// tt2 := randomSession()

}
