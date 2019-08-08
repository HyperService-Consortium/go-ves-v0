package session

import (
	"fmt"
	"testing"

	index "github.com/Myriad-Dreamin/go-ves/lib/database/index"
	xorm_multi_index "github.com/Myriad-Dreamin/go-ves/lib/database/multi_index"
	mtest "github.com/Myriad-Dreamin/mydrest"
)

type TestHelper struct {
	mtest.TestHelper
	res   *xorm_multi_index.XORMMultiIndexImpl
	ser   *index.LevelDBIndex
	logic *SerialSessionBase
}

var s TestHelper

const path = "ves:123456@tcp(127.0.0.1:3306)/ves?charset=utf8"

func SetUpHelper() {
	var err error
	s.res, err = xorm_multi_index.GetXORMMultiIndex("mysql", path)
	s.OutAssertNoErr(err)
	err = s.res.Register(&SerialSession{})
	s.OutAssertNoErr(err)
	s.ser, err = index.GetIndex("./testdb")
	s.OutAssertNoErr(err)
	s.logic = new(SerialSessionBase)
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
	err = s.logic.InsertSessionInfo(s.res, s.ser, tt1)
	s.AssertNoErr(t, err)
	err = s.logic.InsertSessionInfo(s.res, s.ser, tt2)
	s.AssertNoErr(t, err)
	fmt.Println(s.logic.FindSessionInfo(s.res, s.ser, tt1.GetGUID()))
	fmt.Println(s.logic.FindSessionInfo(s.res, s.ser, tt2.GetGUID()))
	err = s.logic.DeleteSessionInfo(s.res, s.ser, tt1.GetGUID())
	s.AssertNoErr(t, err)
	err = s.logic.DeleteSessionInfo(s.res, s.ser, tt2.GetGUID())
	s.AssertNoErr(t, err)
	fmt.Println(s.logic.FindSessionInfo(s.res, s.ser, tt1.GetGUID()))
	fmt.Println(s.logic.FindSessionInfo(s.res, s.ser, tt2.GetGUID()))
}

// func TestMoreLogic(t *testing.T) {
// 	tt1 := randomSession()
// 	tt2 := randomSession()
//
// }
