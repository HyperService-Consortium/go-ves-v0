package index

import (
	"fmt"
	"testing"

	"github.com/syndtr/goleveldb/leveldb"
)

var __x_db *leveldb.DB
var __x_lo *LevelDBIndex

func TestOpen(t *testing.T) {
	var err error
	__x_db, err = leveldb.OpenFile("./testdb", nil)
	if err != nil {
		panic(err)
	}
	__x_lo = &LevelDBIndex{db: __x_db}
}

func TestSetGetDel(t *testing.T) {
	var err error
	err = __x_lo.Set([]byte("k"), []byte("123456"))
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(__x_lo.Get([]byte("k")))
	err = __x_lo.Delete([]byte("k"))
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(__x_lo.Get([]byte("k")))
}

func TestClose(t *testing.T) {
	__x_db.Close()
}
