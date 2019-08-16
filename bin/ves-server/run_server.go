package main

import (
	"flag"
	"fmt"
	"log"

	"encoding/hex"

	signaturer "github.com/Myriad-Dreamin/go-uip/signaturer"
	index "github.com/Myriad-Dreamin/go-ves/lib/database/index"
	multi_index "github.com/Myriad-Dreamin/go-ves/lib/database/multi_index"
	types "github.com/Myriad-Dreamin/go-ves/types"

	ves_server "github.com/Myriad-Dreamin/go-ves/ves"
)

var (
	cfgPath = flag.String("config", "ves-server-config.toml", "configuration of ves server")
)

func main() {
	if *cfgPath != cfgContext {
		ResetPath(*cfgPath)
	}
	var config = Config()
	var err error

	var muldb types.MultiIndex
	switch config.DatabaseConfig.Engine {
	case "xorm":
		var dbConfig = config.DatabaseConfig
		var reqString = fmt.Sprintf(
			"%s:%s@%s(%s)/%s?charset=%s",
			dbConfig.UserName, dbConfig.Password,
			dbConfig.ConnectionType, dbConfig.RemoteHost,
			dbConfig.BaseName, dbConfig.Encoding,
		)

		muldb, err = multi_index.GetXORMMultiIndex(dbConfig.Type, reqString)
		if err != nil {
			log.Fatalf("failed to get muldb: %v", err)
			return
		}
	default:
		log.Fatal("unrecognized database engine")
		return
	}

	var sindb types.Index
	switch config.KVDBConfig.Type {
	case "leveldb":
		sindb, err = index.GetIndex(config.KVDBConfig.Path)
		if err != nil {
			log.Fatalf("failed to get sindb: %v", err)
			return
		}
	default:
		log.Fatal("unrecognized kvdb type")
	}

	b, err := hex.DecodeString("2333bfffffffffffffff2333bbffffffffffffff2333bbffffffffffffffffff2333bfffffffffffffff2333bbffffffffffffff2333bbffffffffffffffffff")
	if err != nil {
		log.Fatal(err)
		return
	}
	signer := signaturer.NewTendermintNSBSigner(b)

	var server *ves_server.Server
	if server, err = ves_server.NewServer(
		muldb, sindb, multi_index.XORMMigrate, signer,
	); err != nil {
		log.Fatal(err)
	}

	if err := server.ListenAndServe(config.ServerConfig.Port, config.ServerConfig.CentralVesAddress); err != nil {
		log.Fatal(err)
	}
}

func init() {
	flag.Parse()
}
