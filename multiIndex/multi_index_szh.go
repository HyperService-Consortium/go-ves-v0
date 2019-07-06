package types

import (
	"errors"
	"fmt"
	"reflect"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

type TT struct {
	Id      int64  `xorm:"pk autoincr"`
	Name    string `xorm:"unique"`
	Balance float64
}
type ORMObj_szh interface {
	GetSlice() interface{}
	GetId() int64
}

type ORMMultiIndex_szh interface {
	RegisterObject(ORMObj_szh) error

	Insert(ORMObj_szh) error

	Select(ORMObj_szh) (interface{}, error)

	SelectAll(ORMObj_szh) (interface{}, error)

	Delete(ORMObj_szh) error

	MultiDelete(ORMObj_szh) error

	Modify(ORMObj_szh, map[string]interface{}) error

	MultiModify(ORMObj_szh, map[string]interface{}) error
}

type XORMImplementation struct {
	db       *xorm.Engine
	regTable map[string]interface{}
}

func (this *XORMImplementation) RegisterObject(obj ORMObj_szh) error {
	if this.db == nil {
		return errors.New("No Database")
	}
	if obj == nil {
		return errors.New("Invalid obj")
	}
	tp := reflect.TypeOf(obj).Name()
	if this.regTable[tp] == nil {
		err := this.db.Sync(obj)
		if err != nil {
			return err
		}
		this.regTable[tp] = obj.GetSlice()
	}
	return nil
}
func (this *XORMImplementation) Insert(obj ORMObj_szh) error {
	affected, err := this.db.Insert(obj)
	fmt.Println("AFFECTED", affected)
	return err
}
func (this *XORMImplementation) Select(condition ORMObj_szh) (interface{}, error) {
	if condition == nil {
		return nil, errors.New("nil condition")
	}
	tp := reflect.TypeOf(condition).Name()
	sli := this.regTable[tp]
	if sli == nil {
		return nil, errors.New("unregistered object")
	}
	err := this.db.Find(sli, condition)
	if err != nil {
		return nil, err
	}
	return reflect.ValueOf(sli).Elem().Interface(), nil
}
func (this *XORMImplementation) SelectAll(obj ORMObj_szh) (interface{}, error) {
	if obj == nil {
		return nil, errors.New("nil object")
	}
	tp := reflect.TypeOf(obj).Name()
	sli := this.regTable[tp]
	if sli == nil {
		return nil, errors.New("unregistered object")
	}
	err := this.db.Find(sli)
	if err != nil {
		return nil, err
	}
	return reflect.ValueOf(sli).Elem().Interface(), nil
}
func (this *XORMImplementation) Delete(obj ORMObj_szh) error {
	has, err := this.db.Get(obj)
	if err != nil {
		return err
	}
	if !has {
		return errors.New("Database does not have the object")
	}
	_, err = this.db.Delete(obj)
	if err != nil {
		return err
	}
	return nil
}
func (this *XORMImplementation) MultiDelete(obj ORMObj_szh) error {
	modified, err := this.db.Delete(obj)
	if err != nil {
		return err
	}
	fmt.Println("DEL SUCC", modified)
	return nil
}
func (this *XORMImplementation) Modify(oldObj ORMObj_szh, newValue map[string]interface{}) error {
	has, err := this.db.Get(oldObj)
	if err != nil {
		return err
	}
	if !has {
		return errors.New("Database does not have the object")
	}
	affected, err := this.db.Table(oldObj).Id(oldObj.GetId()).Update(newValue)
	if err != nil {
		return err
	}
	fmt.Println("MODI SUCC", affected)
	return nil
}
func (this *XORMImplementation) MultiModify(condition ORMObj_szh, newValue map[string]interface{}) error {
	sli, err := this.SelectAll(condition)
	if err != nil {
		return err
	}
	for _, obj := range sli.([]ORMObj_szh) {
		affected, err := this.db.Table(condition).Id(obj.GetId()).Update(newValue)
		if err != nil {
			return err
		}
		fmt.Println("MULMOD", affected)
	}
	return nil
}

type ORMMultiIndexFatory interface {
	GetDB(string, string) (ORMMultiIndex_szh, error)
}

type ORMFactoty_szh struct {
}

func (this *ORMFactoty_szh) GetDB(tp string, pth string) (ORMMultiIndex_szh, error) {
	ret := new(XORMImplementation)
	db, err := xorm.NewEngine(tp, pth)
	if err != nil {
		return nil, err
	}
	ret.db = db
	ret.regTable = make(map[string]interface{})
	return ret, nil
}
