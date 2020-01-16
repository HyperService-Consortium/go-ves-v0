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
	"github.com/Myriad-Dreamin/minimum-lib/logger"
	"github.com/Myriad-Dreamin/minimum-lib/sugar"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"net/http"
	"os"
	"sync"
	"testing"
	"time"

	centered_ves_server "github.com/HyperService-Consortium/go-ves/central-ves"
	multi_index "github.com/HyperService-Consortium/go-ves/lib/database/multi_index"
	ves_server "github.com/HyperService-Consortium/go-ves/ves"

	_ "net/http/pprof"
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

func TestTransfer(t *testing.T) {
	go func() {
		err := http.ListenAndServe("0.0.0.0:22239", nil)
		if err != nil {
			panic(err)
		}
	}()

	go vesclient.StartDaemon()
	h := sugar.NewHandlerErrorLogger(t)

	f, err := os.Create("gin-server.log")
	if err != nil {
		t.Fatal(err)
	}
	gin.DefaultWriter = f
	defer f.Close()

	cfg0, cfg1, cfg2 := zap.NewDevelopmentConfig(), zap.NewDevelopmentConfig(), zap.NewDevelopmentConfig()
	cfg3, cfg4 := zap.NewDevelopmentConfig(), zap.NewDevelopmentConfig()

	//cfg0.OutputPaths = []string{"ves-client.log"}
	//cfg1.OutputPaths = []string{"aws-client.log"}
	//cfg2.OutputPaths = []string{"qwq-client.log"}
	//cfg3.OutputPaths = []string{"central-server.log"}
	//cfg4.OutputPaths = []string{"ves-server.log"}

	cfg0.OutputPaths = []string{"stdout"}
	cfg1.OutputPaths = []string{"stdout"}
	cfg2.OutputPaths = []string{"stdout"}
	cfg3.OutputPaths = []string{"stdout"}
	cfg4.OutputPaths = []string{"stdout"}

	//vesLogger, awsLogger, qwqLogger, cVesServerLogger, vesServerLogger :=
	//	h.HandlerError(logger.NewZapLogger(cfg0, zapcore.DebugLevel)).(logger.Logger),
	//	h.HandlerError(logger.NewZapLogger(cfg1, zapcore.DebugLevel)).(logger.Logger),
	//	h.HandlerError(logger.NewZapLogger(cfg2, zapcore.DebugLevel)).(logger.Logger),
	//	h.HandlerError(logger.NewZapLogger(cfg1, zapcore.DebugLevel)).(logger.Logger),
	//	h.HandlerError(logger.NewZapLogger(cfg2, zapcore.DebugLevel)).(logger.Logger)

	vesLogger, awsLogger, qwqLogger, cVesServerLogger, vesServerLogger :=
		h.HandlerError(logger.NewZapLogger(cfg0, zapcore.DebugLevel)).(logger.Logger).With("from", "ves"),
		h.HandlerError(logger.NewZapLogger(cfg1, zapcore.DebugLevel)).(logger.Logger).With("from", "user1"),
		h.HandlerError(logger.NewZapLogger(cfg2, zapcore.DebugLevel)).(logger.Logger).With("from", "user2"),
		h.HandlerError(logger.NewZapLogger(cfg1, zapcore.DebugLevel)).(logger.Logger).With("from", "cves"),
		h.HandlerError(logger.NewZapLogger(cfg2, zapcore.DebugLevel)).(logger.Logger).With("from", "ves")

	config.ResetPath(cfgPath)
	var cfg = config.Config()
	db := dbinstance.MakeDB()
	go func() {
		if err := h.HandlerError(
			centered_ves_server.NewServer(
				cVesPort, cVesAddr, db, centered_ves_server.NSBHostOption(nsbHost), cVesServerLogger,
			)).(*centered_ves_server.Server).Start(); err != nil {
			t.Fatalf("ListenAndServe: %v\n", err)
		}
	}()
	signer := h.HandlerError(signaturer.NewTendermintNSBSigner(
		h.HandlerError(hex.DecodeString("2333bfffffffffffffff2333bbffffffffffffff2333bbffffffffffffffffff2333bfffffffffffffff2333bbffffffffffffff2333bbffffffffffffffffff")).([]byte))).(*signaturer.TendermintNSBSigner)

	muldb, sindb := Prepare()
	var server = h.HandlerError(ves_server.NewServer(
		muldb, sindb, multi_index.XORMMigrate, signer,
		ves_server.NSBHostOption(nsbHost), vesServerLogger)).(*ves_server.Server)

	go func() {
		if err := server.ListenAndServe(cfg.ServerConfig.Port, cfg.ServerConfig.CentralVesAddress); err != nil {
			t.Fatal(err)
		}
	}()

	time.Sleep(time.Millisecond * 100)

	ves, aws, qwq :=
		h.HandlerError(vesclient.VanilleMakeClient("ves", testServer, vesLogger)).(*vesclient.VesClient),
		h.HandlerError(vesclient.VanilleMakeClient("awsl", testServer, awsLogger)).(*vesclient.VesClient),
		h.HandlerError(vesclient.VanilleMakeClient("qwq", testServer, qwqLogger)).(*vesclient.VesClient)

	var b = make([]byte, 65536)
	h.HandlerError0(ves.ConfigEth("./json/veth.json", b))
	h.HandlerError0(aws.ConfigEth("./json/leth.json", b))
	h.HandlerError0(qwq.ConfigEth("./json/qeth.json", b))
	h.HandlerError0(ves.ConfigKey("./json/vesa.json", b))
	h.HandlerError0(aws.ConfigKey("./json/lswa.json", b))
	h.HandlerError0(qwq.ConfigKey("./json/qwq.json", b))
	var wg = false
	var mx sync.Mutex
	var ch = make(chan bool)
	ves.SubscribeCloseSession(func(sess []byte) {
		fmt.Println("session closed", "session", sess, "ves")
		mx.Lock()
		defer mx.Unlock()
		if wg == false {
			ch <- true
			wg = true
		}
	})
	aws.SubscribeCloseSession(func(sess []byte) {
		fmt.Println("session closed", "session", sess, "aws")
		mx.Lock()
		defer mx.Unlock()
		if wg == false {
			ch <- true
			wg = true
		}
	})
	qwq.SubscribeCloseSession(func(sess []byte) {
		fmt.Println("session closed", "session", sess, "qwq")
		mx.Lock()
		defer mx.Unlock()
		if wg == false {
			ch <- true
			wg = true
		}
	})

	h.HandlerError0(ves.SendOpIntents("./json/intent.json", b))
	<-ch
}
