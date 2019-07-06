package session

import (
	"fmt"
	"testing"

	xorm_multi_index "github.com/Myriad-Dreamin/go-ves/database/multi_index"
	mtest "github.com/Myriad-Dreamin/mydrest"
)

type TestHelper struct {
	mtest.TestHelper
	fact  *xorm_multi_index.XORMMultiIndexFatory
	res   *xorm_multi_index.XORMMultiIndexImpl
	logic *SerialSessionBase
}

var s TestHelper

const path = "ves:123456@tcp(127.0.0.1:3306)/ves?charset=utf8"

func SetUpHelper() {
	var err error
	s.fact = new(xorm_multi_index.XORMMultiIndexFatory)
	s.res, err = s.fact.GetDB("mysql", path)
	s.OutAssertNoErr(err)
	err = s.res.Register(&SerialSession{})
	fmt.Println(err)
	s.OutAssertNoErr(err)
	s.logic = new(SerialSessionBase)
}

func TestMain(m *testing.M) {
	SetUpHelper()
	m.Run()
}

func TestGetDB(t *testing.T) {
	var fact *xorm_multi_index.XORMMultiIndexFatory = new(xorm_multi_index.XORMMultiIndexFatory)
	_, err := fact.GetDB("mysql", path)
	s.AssertNoErr(t, err)
}

func TestRaw(t *testing.T) {
	var err error
	tt1 := randomSession()
	tt2 := randomSession()
	err = s.res.Insert(tt1)
	s.AssertNoErr(t, err)
	err = s.res.Insert(tt2)
	s.AssertNoErr(t, err)
	fmt.Println(s.res.SelectAll(tt1))
	err = s.res.Delete(tt1)
	s.AssertNoErr(t, err)
	err = s.res.Delete(tt2)
	s.AssertNoErr(t, err)
	fmt.Println(s.res.SelectAll(tt1))
}

func TestLogic(t *testing.T) {
	var err error
	tt1 := randomSession()
	tt2 := randomSession()
	fmt.Println(tt1)
	fmt.Println(tt2)
	err = s.logic.InsertSessionInfo(s.res, tt1)
	s.AssertNoErr(t, err)
	err = s.logic.InsertSessionInfo(s.res, tt2)
	s.AssertNoErr(t, err)
	fmt.Println(s.logic.FindSessionInfo(s.res, tt1.GetGUID()))
	fmt.Println(s.logic.FindSessionInfo(s.res, tt2.GetGUID()))
	err = s.logic.DeleteSessionInfo(s.res, tt1.GetGUID())
	s.AssertNoErr(t, err)
	err = s.logic.DeleteSessionInfo(s.res, tt2.GetGUID())
	s.AssertNoErr(t, err)
	fmt.Println(s.logic.FindSessionInfo(s.res, tt1.GetGUID()))
	fmt.Println(s.logic.FindSessionInfo(s.res, tt2.GetGUID()))
}
