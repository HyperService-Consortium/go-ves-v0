package index

import (
	"errors"

	"github.com/syndtr/goleveldb/leveldb"
)

type LevelDBIndex struct {
	db *leveldb.DB
}

func (ldx *LevelDBIndex) Get(b []byte) ([]byte, error) {
	return ldx.db.Get(b, nil)
}

func (ldx *LevelDBIndex) Set(k, v []byte) error {
	return ldx.db.Put(k, v, nil)
}

func (ldx *LevelDBIndex) Delete(b []byte) error {
	return ldx.db.Delete(b, nil)
}

func (ldx *LevelDBIndex) Batch(ks, vs [][]byte) error {
	if len(ks) != len(vs) {
		return errors.New("not match length of keys and values")
	}
	var batch = new(leveldb.Batch)
	for idx, key := range ks {
		batch.Put(key, vs[idx])
	}
	return ldx.db.Write(batch, nil)
}

func GetIndex(filePath string) (*LevelDBIndex, error) {
	db, err := leveldb.OpenFile(filePath, nil)
	if err != nil {
		return nil, err
	}
	return &LevelDBIndex{db: db}, nil
}
