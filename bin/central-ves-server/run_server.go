package main

import (
	"flag"
	"fmt"
	"github.com/HyperService-Consortium/go-ves/config"
	chain_dns "github.com/HyperService-Consortium/go-ves/types/chain-dns"
	storage_handler "github.com/HyperService-Consortium/go-ves/types/storage-handler"
	"log"

	types "github.com/HyperService-Consortium/go-ves/types"
	vesdb "github.com/HyperService-Consortium/go-ves/types/database"
	kvdb "github.com/HyperService-Consortium/go-ves/types/kvdb"
	session "github.com/HyperService-Consortium/go-ves/types/session"
	user "github.com/HyperService-Consortium/go-ves/types/user"

	index "github.com/HyperService-Consortium/go-ves/lib/database/index"
	multi_index "github.com/HyperService-Consortium/go-ves/lib/database/multi_index"

	centered_ves_server "github.com/HyperService-Consortium/go-ves/central-ves"
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
	db.SetSessionKVBase(new(kvdb.Database))
	db.SetStorageHandler(new(storage_handler.Database))
	db.SetChainDNS(chain_dns.NewDatabase(config.HostMap))
	return db
}

func main() {
	flag.Parse()
	srv, err := centered_ves_server.NewServer(port, *addr, makeDB())
	if err != nil {
		log.Fatalf("Create Server: %v\n", err)
	}
	if err = srv.Start(); err != nil {
		log.Fatalf("ListenAndServe: %v\n", err)
	}
}
