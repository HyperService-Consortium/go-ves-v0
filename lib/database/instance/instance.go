package dbinstance

import (
	"fmt"
	"github.com/HyperService-Consortium/go-ves/config"
	"github.com/HyperService-Consortium/go-ves/lib/database/index"
	multi_index "github.com/HyperService-Consortium/go-ves/lib/database/multi_index"
	"github.com/HyperService-Consortium/go-ves/types"
	chain_dns "github.com/HyperService-Consortium/go-ves/types/chain-dns"
	vesdb "github.com/HyperService-Consortium/go-ves/types/database"
	"github.com/HyperService-Consortium/go-ves/types/kvdb"
	"github.com/HyperService-Consortium/go-ves/types/session"
	"github.com/HyperService-Consortium/go-ves/types/storage-handler"
	"github.com/HyperService-Consortium/go-ves/types/user"
)

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

func MakeDB() types.VESDB {

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
	db.SetChainDNS(chain_dns.NewDatabase(config.GetHostMap()))
	return db
}
