package main

import (
	"log"

	"encoding/hex"

	signaturer "github.com/Myriad-Dreamin/go-uip/signaturer"
	index "github.com/Myriad-Dreamin/go-ves/lib/database/index"
	multi_index "github.com/Myriad-Dreamin/go-ves/lib/database/multi_index"

	ves_server "github.com/Myriad-Dreamin/go-ves/ves"
)

const port = ":23351"
const centerAddress = "127.0.0.1:23352"

func main() {

	var err error

	//TODO: SetEnv
	var muldb *multi_index.XORMMultiIndexImpl
	muldb, err = multi_index.GetXORMMultiIndex("mysql", "ves:123456@tcp(127.0.0.1:3306)/ves?charset=utf8")
	if err != nil {
		log.Fatalf("failed to get muldb: %v", err)
		return
	}
	var sindb *index.LevelDBIndex
	sindb, err = index.GetIndex("./data")
	if err != nil {
		log.Fatalf("failed to get sindb: %v", err)
		return
	}

	b, err := hex.DecodeString("2333bbffffffffffffff2333bbffffffffffffff2333bbffffffffffffffffff2333bbffffffffffffff2333bbffffffffffffff2333bbffffffffffffffffff")
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

	if err := server.ListenAndServe(port, centerAddress); err != nil {
		log.Fatal(err)
	}
}
