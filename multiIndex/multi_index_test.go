package types

import (
	"fmt"
	"testing"
)

func newTT(name string, balance float64) *TT {
	ret := new(TT)
	ret.Name = name
	ret.Balance = balance
	return ret
}
func (this *TT) GetSlice() interface{} {
	var ret []TT
	return &ret
}
func (this *TT) GetId() int64 {
	return this.Id
}

const path = "root:12345678@tcp(127.0.0.1:3306)/test?charset=utf8"

func TestGetDB(t *testing.T) {
	var fact ORMMultiIndexFatory = new(ORMFactoty_szh)
	res, err := fact.GetDB("mysql", path)
	if err != nil {
		t.Error(err, "\n")
	} else {
		fmt.Printf("RES %T\n", res)
	}
}

func TestInsert(t *testing.T) {
	fact := new(ORMFactoty_szh)
	res, err := fact.GetDB("mysql", path)
	err = res.RegisterObject(new(TT))
	if err != nil {
		t.Error("reg err", err)
		return
	}
	tt1 := newTT("szh1", 10)
	tt2 := newTT("szh2", 11)
	err = res.Insert(tt1)
	if err != nil {
		t.Error("tt1 err", err)
	}
	err = res.Insert(tt2)
	if err != nil {
		t.Error("tt2 err", err)
	}
}

func TestSelect(t *testing.T) {
	fact := new(ORMFactoty_szh)
	res, err := fact.GetDB("mysql", path)
	err = res.RegisterObject(new(TT))
	if err != nil {
		t.Error("reg err", err)
		return
	}
	tt1 := newTT("szh1", 10)
	tt2 := newTT("szh2", 11)
	err = res.Insert(tt1)
	err = res.Insert(tt2)
	condition := new(TT)
	condition.Balance = 11
	var result interface{}
	result, err = res.Select(condition)
	if err != nil {
		t.Error("SELECT ERR", err)
	}
	/*
		fmt.Println(result.([]TT), res.regTable[reflect.TypeOf(tt1).Name()])
		(result.([]TT))[0].Balance = 15
		fmt.Println(res.regTable[reflect.TypeOf(tt1).Name()])
		(*((res.regTable[reflect.TypeOf(tt1).Name()]).(*[]TT)))[0].Balance = 20
	*/
	fmt.Println(result.([]TT))
}

func TestDelete(t *testing.T) {
	fact := new(ORMFactoty_szh)
	res, err := fact.GetDB("mysql", path)
	err = res.RegisterObject(new(TT))
	if err != nil {
		t.Error("reg err", err)
		return
	}
	tt1 := newTT("szh1", 10)
	tt2 := newTT("szh2", 11)
	err = res.Insert(tt1)
	err = res.Insert(tt2)
	err = res.Delete(tt1)
	if err != nil {
		t.Error("DELETE", err)
	}
	sb := &TT{Balance: 11}
	err = res.Delete(sb)
	if err != nil {
		t.Error("DELETE", err)
	}
	err = res.Delete(newTT("SD", 1000))
	if err == nil {
		t.Error("SB DELETE")
	}

}

func TestMultiDelete(t *testing.T) {
	fact := new(ORMFactoty_szh)
	res, err := fact.GetDB("mysql", path)
	err = res.RegisterObject(new(TT))
	if err != nil {
		t.Error("reg err", err)
		return
	}
	tt1 := newTT("szh1", 10)
	tt2 := newTT("szh2", 11)
	tt3 := newTT("szh3", 10)
	err = res.Insert(tt1)
	err = res.Insert(tt2)
	err = res.Insert(tt3)
	err = res.MultiDelete(&TT{Balance: 10})
	if err != nil {
		t.Error("MULTIDELETE", err)
	}
}

func TestModify(t *testing.T) {
	fact := new(ORMFactoty_szh)
	res, err := fact.GetDB("mysql", path)
	err = res.RegisterObject(new(TT))
	if err != nil {
		t.Error("reg err", err)
		return
	}
	tt1 := newTT("szh1", 10)
	tt2 := newTT("szh2", 11)
	tt3 := newTT("szh3", 10)
	err = res.Insert(tt1)
	err = res.Insert(tt2)
	err = res.Insert(tt3)
	mod := map[string]interface{}{"Balance": 17}
	err = res.Modify(tt2, mod)
	if err != nil {
		t.Error("MODIFY", err)
	}
}
