package open

import (
	"encoding/base64"
	"math/big"
	"os"
	"testing"

	"github.com/boltdb/bolt"
)

func BenchmarkOpenBolt(b *testing.B) {

	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		b.Error(err)
		return
	}
	defer db.Close()
	var s, t, tt, _1, _2 = big.NewInt(5300).SetBytes(append(make([]byte, 30), 0x03, 0x53)), big.NewInt(5300).SetBytes(append(make([]byte, 30), 0x53, 0x53)), big.NewInt(5300).SetBytes(append(make([]byte, 70), 0x53, 0x53)), big.NewInt(1), big.NewInt(2)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {

		db.Batch(func(tx *bolt.Tx) error {
			for i := 0; i < 2000; i++ {
				buc, err := tx.CreateBucketIfNotExists(s.Bytes())
				if err != nil {
					return err
				}
				for j := 0; j < 20; j++ {
					err := buc.Put(t.Bytes(), tt.Bytes())
					if err != nil {
						return err
					}
					t.Add(t, _1)
					tt.Add(tt, _2)
				}
				s.Add(s, _1)
			}
			return nil
		})
	}
}

func BenchmarkOpenBolt2(b *testing.B) {

	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		b.Error(err)
		return
	}
	defer db.Close()
	var s, t, tt, _1, _2 = big.NewInt(5300).SetBytes(append(make([]byte, 30), 0x23, 0x53)), big.NewInt(5300).SetBytes(append(make([]byte, 30), 0x53, 0x53)), big.NewInt(5300).SetBytes(append(make([]byte, 70), 0x53, 0x53)), big.NewInt(1), big.NewInt(2)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {

		db.Batch(func(tx *bolt.Tx) error {
			for i := 0; i < 2000; i++ {
				buc, err := tx.CreateBucketIfNotExists(s.Bytes())
				if err != nil {
					return err
				}
				for j := 0; j < 20; j++ {
					err := buc.Put(t.Bytes(), tt.Bytes())
					if err != nil {
						return err
					}
					t.Add(t, _1)
					tt.Add(tt, _2)
				}
				s.Add(s, _1)
			}
			return nil
		})
	}
}

func BenchmarkOpenBolt3(b *testing.B) {

	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		b.Error(err)
		return
	}
	defer db.Close()
	var s, t, tt, _1, _2 = big.NewInt(5300).SetBytes(append(make([]byte, 30), 0x53, 0x53)), big.NewInt(5300).SetBytes(append(make([]byte, 30), 0x53, 0x53)), big.NewInt(5300).SetBytes(append(make([]byte, 70), 0x53, 0x53)), big.NewInt(1), big.NewInt(2)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {

		db.Batch(func(tx *bolt.Tx) error {
			for i := 0; i < 2000; i++ {
				buc, err := tx.CreateBucketIfNotExists(s.Bytes())
				if err != nil {
					return err
				}
				for j := 0; j < 20; j++ {
					err := buc.Put(t.Bytes(), tt.Bytes())
					if err != nil {
						return err
					}
					t.Add(t, _1)
					tt.Add(tt, _2)
				}
				s.Add(s, _1)
			}
			return nil
		})
	}
}

func BenchmarkOpenBolt4(b *testing.B) {

	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		b.Error(err)
		return
	}
	defer db.Close()
	var s, t, tt, _1, _2 = big.NewInt(5300).SetBytes(append(make([]byte, 30), 0x53, 0x53)), big.NewInt(5300).SetBytes(append(make([]byte, 30), 0x53, 0x53)), big.NewInt(5300).SetBytes(append(make([]byte, 70), 0x53, 0x53)), big.NewInt(1), big.NewInt(2)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {

		db.Batch(func(tx *bolt.Tx) error {
			for i := 0; i < 2000; i++ {
				buc, err := tx.CreateBucketIfNotExists(s.Bytes())
				if err != nil {
					return err
				}
				for j := 0; j < 20; j++ {
					err := buc.Put(t.Bytes(), tt.Bytes())
					if err != nil {
						return err
					}
					t.Add(t, _1)
					tt.Add(tt, _2)
				}
				s.Add(s, _1)
			}
			return nil
		})
	}
}

func BenchmarkOpenBolt5(b *testing.B) {

	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		b.Error(err)
		return
	}
	defer db.Close()
	var s, t, tt, _1, _2 = big.NewInt(5300).SetBytes(append(make([]byte, 30), 0x73, 0x53)), big.NewInt(5300).SetBytes(append(make([]byte, 30), 0x53, 0x53)), big.NewInt(5300).SetBytes(append(make([]byte, 70), 0x53, 0x53)), big.NewInt(1), big.NewInt(2)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {

		db.Batch(func(tx *bolt.Tx) error {
			for i := 0; i < 2000; i++ {
				buc, err := tx.CreateBucketIfNotExists(s.Bytes())
				if err != nil {
					return err
				}
				for j := 0; j < 20; j++ {
					err := buc.Put(t.Bytes(), tt.Bytes())
					if err != nil {
						return err
					}
					t.Add(t, _1)
					tt.Add(tt, _2)
				}
				s.Add(s, _1)
			}
			return nil
		})
	}
}

func BenchmarkBackOut(b *testing.B) {

	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		b.Error(err)
		return
	}
	defer db.Close()
	var s, _1 = big.NewInt(5300).SetBytes(append(make([]byte, 30), 0x53, 0x53)), big.NewInt(1)

	var buf = make([]byte, 130)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {

		err = db.Batch(func(tx *bolt.Tx) error {
			for i := 0; i < 2000; i++ {
				buc := tx.Bucket(s.Bytes())
				fi, err := os.Create("bac/" + s.String())
				if err != nil {
					return err
				}
				buc.ForEach(func(k, v []byte) error {
					base64.StdEncoding.Encode(buf, k)
					ll := base64.StdEncoding.EncodedLen(len(k))
					buf[ll] = '\x19'
					fi.Write(buf[:ll+1])
					base64.StdEncoding.Encode(buf, v)
					ll = base64.StdEncoding.EncodedLen(len(v))
					buf[ll] = '\n'
					fi.Write(buf[:ll+1])
					return nil
				})
				fi.Close()
				err = tx.DeleteBucket(s.Bytes())
				if err != nil {
					return err
				}
				s.Add(s, _1)
			}
			return nil
		})
		if err != nil {
			b.Error(err)
			return
		}
	}
}
