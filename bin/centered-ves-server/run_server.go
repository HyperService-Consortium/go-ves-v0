package main

import (
	"flag"
	"fmt"
	"log"

	centered_ves_server "github.com/Myriad-Dreamin/go-ves/centered_ves"
	"github.com/Myriad-Dreamin/go-ves/database/index"
	multi_index "github.com/Myriad-Dreamin/go-ves/database/multi_index"
	"github.com/Myriad-Dreamin/go-ves/types"
	vesdb "github.com/Myriad-Dreamin/go-ves/types/database"
	"github.com/Myriad-Dreamin/go-ves/types/session"
	"github.com/Myriad-Dreamin/go-ves/types/user"
)

const port = ":23352"

var addr = flag.String("port", ":23452", "http service address")

func XORMMigrate(muldb types.MultiIndex) (err error) {
	var xorm_muldb = muldb.(*multi_index.XORMMultiIndexImpl)
	err = xorm_muldb.Register(&user.XORMUserAdapter{})
	if err != nil {
		return
	}
	err = xorm_muldb.Register(&session.SerialSession{})
	if err != nil {
		return
	}
	return nil
}

func makeDB() types.VESDB {

	var db = new(vesdb.Database)
	var err error

	//TODO: SetEnv
	var muldb *multi_index.XORMMultiIndexImpl
	muldb, err = multi_index.GetXORMMultiIndex("mysql", "ves:123456@tcp(127.0.0.1:3306)/ves?charset=utf8")
	if err != nil {
		panic(fmt.Errorf("failed to get muldb: %v", err))
	}
	err = XORMMigrate(muldb)
	if err != nil {
		panic(fmt.Errorf("failed to migrate: %v", err))
	}

	var sindb *index.LevelDBIndex
	sindb, err = index.GetIndex("./index_data")
	if err != nil {
		panic(fmt.Errorf("failed to get sindb: %v", err))
	}

	db.SetIndex(sindb)
	db.SetMultiIndex(muldb)

	db.SetUserBase(new(user.XORMUserBase))
	db.SetSessionBase(new(session.SerialSessionBase))
	return db
}

func main() {
	flag.Parse()
	if err := centered_ves_server.NewServer(port, *addr, makeDB()).Start(); err != nil {
		log.Fatalf("ListenAndServe: %v\n", err)
	}
}
