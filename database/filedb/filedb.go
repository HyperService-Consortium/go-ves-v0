package filedb

import (
	"crypto/rand"
	"io"
	"io/ioutil"
	"os"

	"encoding/gob"
)

var (
	salt     []byte
	dataPath = "./data"
	saltPath = dataPath + "/salt.dat"
	perm     = os.FileMode(0755)
)

func readFile(path string) ([]byte, error) {
	file, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0755)
	if err != nil {
		return nil, err
	}
	var a []byte
	a, err = ioutil.ReadAll(file)
	file.Close()
	if err != nil {
		return nil, err
	}
	return a, nil
}

func preloadSalt() []byte {
	s, err := readFile(saltPath)
	if err != nil {
		panic(err)
	}
	if len(s) < 10 {
		s = make([]byte, 256)
		var n int
		n, err = io.ReadFull(rand.Reader, s)
		if err != nil {
			panic(err)
		}
		if n < 256 {
			panic("the size of data produced by rand.Reader is less than 256 bytes")
		}
		ioutil.WriteFile(saltPath, s, perm)
	}
	return s
}

type FileDB struct {
	dbpath string
}

var pobj = &FileDB{dbpath: dataPath}

func Register(v interface{}) {
	gob.Register(v)
}

func NewFileDB(path string) (*FileDB, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err = os.MkdirAll(path, perm); err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}
	return &FileDB{dbpath: path}, nil
}

func (fdb *FileDB) ReStatDir(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err = os.MkdirAll(path, perm); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}
	fdb.dbpath = path
	return nil
}

type ReadEvent struct {
	file *os.File
	dec  *gob.Decoder
}

func (fdb *FileDB) ReadWithPath(name string) (e *ReadEvent, err error) {
	e = new(ReadEvent)
	e.file, err = os.OpenFile(fdb.dbpath+name, os.O_RDONLY|os.O_CREATE, 0755)
	if err != nil {
		return nil, err
	}
	e.dec = gob.NewDecoder(e.file)
	return e, nil
}

func ReadWithPath(name string) (e *ReadEvent, err error) {
	e = new(ReadEvent)
	e.file, err = os.OpenFile(pobj.dbpath+name, os.O_RDONLY|os.O_CREATE, 0755)
	if err != nil {
		return nil, err
	}
	e.dec = gob.NewDecoder(e.file)
	return e, nil
}

func (e *ReadEvent) Decode(obj interface{}) error {
	return e.dec.Decode(obj)
}

func (e *ReadEvent) Settle() error {
	return e.file.Close()
}

type WriteEvent struct {
	file *os.File
	enc  *gob.Encoder
}

func (fdb *FileDB) WriteWithPath(name string) (e *WriteEvent, err error) {
	e = new(WriteEvent)
	e.file, err = os.OpenFile(fdb.dbpath+name, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return nil, err
	}
	e.enc = gob.NewEncoder(e.file)
	return e, nil
}

func WriteWithPath(name string) (e *WriteEvent, err error) {
	e = new(WriteEvent)
	e.file, err = os.Create(pobj.dbpath + name)
	if err != nil {
		return nil, err
	}
	e.enc = gob.NewEncoder(e.file)
	return e, nil
}

func (e *WriteEvent) Encode(obj interface{}) error {
	return e.enc.Encode(obj)
}

func (e *WriteEvent) Settle() error {
	return e.file.Close()
}

// func ReadWithReader()

func init() {
	if _, err := os.Stat(dataPath); os.IsNotExist(err) {
		if err = os.MkdirAll(dataPath, perm); err != nil {
			panic(err)
		}
	} else if err != nil {
		panic(err)
	}
	salt = preloadSalt()
}
