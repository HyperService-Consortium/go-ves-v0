package open

import (
	_ "database/sql"
	_ "testing"

	_ "github.com/boltdb/bolt"
	_ "github.com/mattn/go-sqlite3"
	_ "github.com/prologic/bitcask"
	_ "github.com/syndtr/goleveldb/leveldb"
)

// func BenchmarkOpenBolt(b *testing.B) {
// 	for i := 0; i < b.N; i++ {
// 		db, err := bolt.Open("my.db", 0600, nil)
// 		if err != nil {
// 			b.Error(err)
// 			return
// 		}
// 		db.Close()
// 	}
// }
//
// func BenchmarkOpenLevel(b *testing.B) {
// 	for i := 0; i < b.N; i++ {
// 		db, err := leveldb.OpenFile("./mylevel.db", nil)
// 		if err != nil {
// 			b.Error(err)
// 			return
// 		}
// 		db.Close()
// 	}
// }
//
// func BenchmarkOpenSqlite(b *testing.B) {
// 	db, err := sql.Open("sqlite3", "./mysqlite.db")
// 	if err != nil {
// 		b.Error(err)
// 		return
// 	}
// 	var sqlTable = `
//     CREATE TABLE IF NOT EXISTS "userinfo" (
//       "uid" INTEGER PRIMARY KEY AUTOINCREMENT,
//       "username" VARCHAR(64) NULL,
//       "departname" VARCHAR(64) NULL,
//       "created" TIMESTAMP default (datetime('now', 'localtime'))
//     );
//     CREATE TABLE IF NOT EXISTS "userdeatail" (
//       "uid" INT(10) NULL,
//       "intro" TEXT NULL,
//       "profile" TEXT NULL,
//       PRIMARY KEY (uid)
//     );
//   `
// 	db.Exec(sqlTable)
// 	db.Close()
// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		db, err := sql.Open("sqlite3", "./mysqlite.db")
// 		if err != nil {
// 			b.Error(err)
// 			return
// 		}
// 		db.Close()
// 	}
// }
//
// func BenchmarkOpenBitcask(b *testing.B) {
// 	for i := 0; i < b.N; i++ {
// 		db, err := bitcask.Open("./mybitcask")
// 		if err != nil {
// 			b.Error(err)
// 			return
// 		}
// 		db.Close()
// 	}
// }

// func BenchmarkOpenBadger(b *testing.B) {
// 	for i := 0; i < b.N; i++ {
// 		db, err := badger.Open(badger.DefaultOptions("./mybadger"))
// 		if err != nil {
// 			b.Error(err)
// 			return
// 		}
// 		db.Close()
// 	}
// }
