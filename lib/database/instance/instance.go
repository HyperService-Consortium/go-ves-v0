package dbinstance

import (
	"fmt"
	"github.com/HyperService-Consortium/go-ves/lib/database/index"
	multi_index "github.com/HyperService-Consortium/go-ves/lib/database/multi_index"
	"github.com/HyperService-Consortium/go-ves/types"
	chain_dns "github.com/HyperService-Consortium/go-ves/types/chain-dns"
	vesdb "github.com/HyperService-Consortium/go-ves/types/database"
	"github.com/HyperService-Consortium/go-ves/types/kvdb"
	"github.com/HyperService-Consortium/go-ves/types/session"
	"github.com/HyperService-Consortium/go-ves/types/user"
	"github.com/HyperService-Consortium/go-uip/const/chain_type"
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
	db.SetChainDNS(chain_dns.NewDatabase(chain_dns.HostMap{
		1: chain_dns.ChainInfo{
			Host:      "127.0.0.1:8545",
			ChainType: ChainType.Ethereum,
		},
		2: chain_dns.ChainInfo{
			Host:      "127.0.0.1:8545",
			ChainType: ChainType.Ethereum,
		},
		3: chain_dns.ChainInfo{
			Host:      "127.0.0.1:26657",
			ChainType: ChainType.TendermintNSB,
		},
		4: chain_dns.ChainInfo{
			Host:      "127.0.0.1:26657",
			ChainType: ChainType.TendermintNSB,
		},
		5: chain_dns.ChainInfo{
			Host:      "127.0.0.1:26657",
			ChainType: ChainType.TendermintNSB,
		},
	}))
	return db
}
