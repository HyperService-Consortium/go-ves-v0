package xorm_multi_index

import (
	"fmt"
	"testing"

	"github.com/Myriad-Dreamin/go-ves/types"
	mtest "github.com/Myriad-Dreamin/mydrest"
)

type TT struct {
	Id      int64  `xorm:"pk autoincr"`
	Name    string `xorm:"unique"`
	Balance float64
}

type TestHelper struct {
	mtest.TestHelper
	res types.MultiIndex
}

var s TestHelper

func (this *TT) ToKVMap() map[string]interface{} {
	return map[string]interface{}{
		"id":      this.Id,
		"name":    this.Name,
		"balance": this.Balance,
	}
}

func (this *TT) GetSlicePtr() interface{} {
	return new([]TT)
}

func (this *TT) GetObjectPtr() interface{} {
	return &TT{}
}

func (this *TT) GetID() int64 {
	return this.Id
}

func newTT(name string, balance float64) *TT {
	ret := new(TT)
	ret.Name = name
	ret.Balance = balance
	return ret
}

const path = "ves:123456@tcp(127.0.0.1:3306)/ves?charset=utf8"

func SetUpHelper() {
	var err error
	s.res, err = GetXORMMultiIndex("mysql", path)
	s.OutAssertNoErr(err)
	// err = s.res.RegisterObject(new(TT))
	s.OutAssertNoErr(err)
}

func TestMain(m *testing.M) {
	SetUpHelper()
	m.Run()
}

func TestGetDB(t *testing.T) {
	_, err := GetXORMMultiIndex("mysql", path)
	s.AssertNoErr(t, err)
}

func TestInsert(t *testing.T) {
	var err error
	tt1 := newTT("szh1", 10)
	tt2 := newTT("szh2", 11)
	err = s.res.Insert(tt1)
	s.AssertNoErr(t, err)
	err = s.res.Insert(tt2)
	s.AssertNoErr(t, err)
	fmt.Println(s.res.SelectAll(tt1))
}

func TestSelect(t *testing.T) {
	var err error
	condition := new(TT)
	condition.Balance = 11
	var result interface{}
	result, err = s.res.Select(condition)
	s.AssertNoErr(t, err)
	fmt.Println(result.([]TT))
}

func TestDelete(t *testing.T) {
	var err error
	tt1 := newTT("szh1", 10)
	err = s.res.Delete(tt1)
	s.AssertNoErr(t, err)
	sb := &TT{Balance: 11}
	err = s.res.Delete(sb)
	s.AssertNoErr(t, err)
	err = s.res.Delete(newTT("SD", 1000))
	s.AssertEqual(t, err, errorObjectNotFound)

}

func TestMultiDelete(t *testing.T) {
	fmt.Println("8")
	var err error
	tt1 := newTT("szh1", 10)
	tt2 := newTT("szh2", 11)
	tt3 := newTT("szh3", 10)
	err = s.res.Insert(tt1)
	s.AssertNoErr(t, err)
	err = s.res.Insert(tt2)
	s.AssertNoErr(t, err)
	err = s.res.Insert(tt3)
	s.AssertNoErr(t, err)
	err = s.res.MultiDelete(&TT{Balance: 10})
	s.AssertNoErr(t, err)
	var result interface{}
	result, err = s.res.SelectAll(tt1)
	s.AssertNoErr(t, err)
	fmt.Println(result.([]TT))
}

func TestModify(t *testing.T) {
	fmt.Println("9")
	var err error
	tt2 := newTT("szh2", 11)
	var result interface{}
	result, err = s.res.Select(tt2)
	s.AssertNoErr(t, err)
	tt2 = &result.([]TT)[0]
	mod := map[string]interface{}{"Balance": 17}
	err = s.res.Modify(tt2, mod)
	s.AssertNoErr(t, err)
	result, err = s.res.SelectAll(tt2)
	s.AssertNoErr(t, err)
	fmt.Println(result.([]TT))
}

func TestClearTable(t *testing.T) {
	fmt.Println("10")
	var err error
	err = s.res.MultiDelete(&TT{Balance: 10})
	s.AssertNoErr(t, err)
	err = s.res.MultiDelete(&TT{Balance: 17})
	s.AssertNoErr(t, err)
}
