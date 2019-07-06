package types

import (
	"errors"
	"reflect"

	"github.com/Myriad-Dreamin/go-ves/types"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

var (
	errorObjectNotFound      = errors.New("Database does not have the object")
	notHandleDBError         = errors.New("XORMMultiIndex does not handle a database engine")
	notRegisteredObjectError = errors.New("unregistered object")
	nilObjError              = errors.New("null object can't be register")
	nilConditionError        = errors.New("null condition is not allowed")
)

type XORMMultiIndexImpl struct {
	db *xorm.Engine
	// regTable map[string]types.KVObject
}

// func (this *XORMMultiIndexImpl) RegisterObject(obj types.KVObject) error {
// 	if this.db == nil {
// 		return notHandleDBError
// 	}
// 	if obj == nil {
// 		return nilObjError
// 	}
// 	tp := reflect.TypeOf(obj).Name()
// 	if this.regTable[tp] == nil {
// 		err := this.db.Sync(obj)
// 		if err != nil {
// 			return err
// 		}
// 		this.regTable[tp] = obj
// 	}
// 	return nil
// }

// func (this *XORMMultiIndexImpl) GetSliceFactory(condition types.KVObject) (types.KVObject, error) {
// 	if condition == nil {
// 		return nil, nilConditionError
// 	}
// 	tp := reflect.TypeOf(condition).Name()
// 	sliFactory := this.regTable[tp]
// 	if sliFactory == nil {
// 		return nil, notRegisteredObjectError
// 	}
// 	return sliFactory, nil
// }

func (this *XORMMultiIndexImpl) Insert(obj types.KVObject) (err error) {
	_, err = this.db.Insert(obj)
	return
}
func (this *XORMMultiIndexImpl) Get(condition types.KVObject) (bool, error) {
	has, err := this.db.Get(condition)
	return has, err
}

func (this *XORMMultiIndexImpl) Select(condition types.KVObject) (interface{}, error) {
	sli := condition.GetSlicePtr()
	err := this.db.Find(sli, condition)
	if err != nil {
		return nil, err
	}
	return reflect.ValueOf(sli).Elem().Interface(), nil
}
func (this *XORMMultiIndexImpl) SelectAll(obji types.KVObject) (interface{}, error) {
	sli := obji.GetSlicePtr()
	err := this.db.Find(sli)
	if err != nil {
		return nil, err
	}
	return reflect.ValueOf(sli).Elem().Interface(), nil
}

func (this *XORMMultiIndexImpl) Delete(obj types.KVObject) error {
	has, err := this.db.Get(obj)
	if err != nil {
		return err
	}
	if !has {
		return errorObjectNotFound
	}
	_, err = this.db.Delete(obj)
	if err != nil {
		return err
	}
	return nil
}
func (this *XORMMultiIndexImpl) MultiDelete(obj types.KVObject) (err error) {
	_, err = this.db.Delete(obj)
	// fmt.Println("DEL SUCC", modified)
	return
}
func (this *XORMMultiIndexImpl) Modify(oldObj types.KVObject, newValue types.KVMap) error {
	has, err := this.db.Get(oldObj)
	if err != nil {
		return err
	}
	if !has {
		return errorObjectNotFound
	}
	_, err = this.db.Table(oldObj).Id(oldObj.GetId()).Update(newValue)
	// fmt.Println("MODI SUCC", affected)
	return err
}
func (this *XORMMultiIndexImpl) MultiModify(condition types.KVObject, newValue types.KVMap) error {
	sli, err := this.SelectAll(condition)
	if err != nil {
		return err
	}
	for _, obj := range sli.([]types.KVObject) {
		_, err := this.db.Table(condition).Id(obj.GetId()).Update(newValue)
		if err != nil {
			return err
		}
		// fmt.Println("MULMOD", affected)
	}
	return nil
}

type ORMMultiIndexFatory interface {
	GetDB(string, string) (types.MultiIndex, error)
}

type XORMMultiIndexFatory struct {
}

func (this *XORMMultiIndexFatory) GetDB(tp string, pth string) (types.MultiIndex, error) {
	ret := new(XORMMultiIndexImpl)
	db, err := xorm.NewEngine(tp, pth)
	if err != nil {
		return nil, err
	}
	ret.db = db
	// ret.regTable = make(map[string]types.KVObject)
	return ret, nil
}
