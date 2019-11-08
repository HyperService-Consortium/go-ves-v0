package transfer_test

import (
	"encoding/hex"
	"fmt"
	"github.com/HyperService-Consortium/go-uip/signaturer"
	"github.com/HyperService-Consortium/go-ves/config"
	"github.com/HyperService-Consortium/go-ves/lib/database/index"
	dbinstance "github.com/HyperService-Consortium/go-ves/lib/database/instance"
	vesclient "github.com/HyperService-Consortium/go-ves/lib/net/ves-client"
	"github.com/HyperService-Consortium/go-ves/types"
	"github.com/Myriad-Dreamin/minimum-lib/sugar"
	"log"
	"testing"
	"time"

	centered_ves_server "github.com/HyperService-Consortium/go-ves/central-ves"
	multi_index "github.com/HyperService-Consortium/go-ves/lib/database/multi_index"
	ves_server "github.com/HyperService-Consortium/go-ves/ves"
)

const testServer = "localhost:23452"
const cVesPort, cVesAddr = ":23352", ":23452"
const cfgPath = "./ves-server-config.toml"
const nsbHost = "127.0.0.1:26657"


func Prepare() (muldb types.MultiIndex, sindb types.Index) {
	var cfg = config.Config()
	var err error

	switch cfg.DatabaseConfig.Engine {
	case "xorm":
		var dbConfig = cfg.DatabaseConfig
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

	switch cfg.KVDBConfig.Type {
	case "leveldb":
		sindb, err = index.GetIndex(cfg.KVDBConfig.Path)
		if err != nil {
			log.Fatalf("failed to get sindb: %v", err)
			return
		}
	default:
		log.Fatal("unrecognized kvdb type")
	}

	return muldb, sindb
}

type Logger interface {
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
}
type handlerError struct {
	logger Logger
}

func (h handlerError) HandlerError(err error) {
	if err != nil {
		sugar.PrintStack()
		h.logger.Fatal(err)
	}
}

func (h handlerError) HandlerError2(i interface{}, err error) interface{} {
	if err != nil {
		sugar.PrintStack()
		h.logger.Fatal(err)
	}
	return i
}


func TestTransfer(t *testing.T) {
	h := handlerError{logger: t}
	config.ResetPath(cfgPath)
	var cfg = config.Config()
	db := dbinstance.MakeDB()
	go func() {
		if err := h.HandlerError2(
			centered_ves_server.NewServer(
				cVesPort, cVesAddr, db, centered_ves_server.NSBHostOption(nsbHost),
			)).(*centered_ves_server.Server).Start(); err != nil {
			t.Fatalf("ListenAndServe: %v\n", err)
		}
	}()
	signer := signaturer.NewTendermintNSBSigner(
		h.HandlerError2(hex.DecodeString("2333bfffffffffffffff2333bbffffffffffffff2333bbffffffffffffffffff2333bfffffffffffffff2333bbffffffffffffff2333bbffffffffffffffffff"),
	).([]byte))

	muldb, sindb := Prepare()
	var server = h.HandlerError2(ves_server.NewServer(
		muldb, sindb, multi_index.XORMMigrate, signer,
		ves_server.NSBHostOption(nsbHost))).(*ves_server.Server)

	go func() {
		if err := server.ListenAndServe(cfg.ServerConfig.Port, cfg.ServerConfig.CentralVesAddress); err != nil {
			t.Fatal(err)
		}
	}()

	time.Sleep(time.Millisecond * 100)
	ves, aws, qwq :=
		h.HandlerError2(vesclient.VanilleMakeClient("ves", testServer)).(*vesclient.VesClient),
		h.HandlerError2(vesclient.VanilleMakeClient("awsl", testServer)).(*vesclient.VesClient),
		h.HandlerError2(vesclient.VanilleMakeClient("qwq", testServer)).(*vesclient.VesClient)

	var b = make([]byte, 65536)
	h.HandlerError(ves.ConfigEth("./json/veth.json", b))
	h.HandlerError(aws.ConfigEth("./json/leth.json", b))
	h.HandlerError(qwq.ConfigEth("./json/qeth.json", b))
	h.HandlerError(ves.ConfigKey("./json/vesa.json", b))
	h.HandlerError(aws.ConfigKey("./json/lswa.json", b))
	h.HandlerError(qwq.ConfigKey("./json/qwq.json", b))


	h.HandlerError(ves.SendOpIntents("./json/intent.json", b))
}

