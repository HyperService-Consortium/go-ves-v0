package ves

import (
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/HyperService-Consortium/go-ves/config"
	helper "github.com/HyperService-Consortium/go-ves/lib/net/help-func"
	nsbcli "github.com/HyperService-Consortium/go-ves/lib/net/nsb-client"
	chain_dns "github.com/HyperService-Consortium/go-ves/types/chain-dns"
	"github.com/HyperService-Consortium/go-ves/types/kvdb"
	"github.com/HyperService-Consortium/go-ves/types/storage-handler"
	"github.com/HyperService-Consortium/go-ves/ves/vs"
	"github.com/Myriad-Dreamin/minimum-lib/logger"
	"io"
	"net"
	"time"

	"github.com/HyperService-Consortium/go-uip/signaturer"
	"github.com/HyperService-Consortium/go-ves/grpc/uiprpc"
	uipbase "github.com/HyperService-Consortium/go-ves/grpc/uiprpc-base"
	log "github.com/HyperService-Consortium/go-ves/lib/log"
	"github.com/HyperService-Consortium/go-ves/types"
	vesdb "github.com/HyperService-Consortium/go-ves/types/database"
	"github.com/HyperService-Consortium/go-ves/types/session"
	"github.com/HyperService-Consortium/go-ves/types/user"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gogo/protobuf/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

// Server provides the basic service of session
type Server vs.VServer

// func (s *Server) locmup(mupper uint16) {
// 	if !s.mup[mupper] {
// 		s.mutex.Lock()
// 		s.mup[mupper] = true
// 		s.mutex.Unlock()
// 	}
// }

// MigrateFunction is used to make migration by passing kv-objects
type MigrateFunction = func(types.MultiIndex, types.KVObject) error

func migrate(
	muldb types.MultiIndex,
	makeMigrate MigrateFunction,
) error {
	if err := makeMigrate(muldb, &user.XORMUserAdapter{}); err != nil {
		return err
	}
	if err := makeMigrate(muldb, &session.MultiThreadSerialSession{}); err != nil {
		return err
	}
	return nil
}

type NSBHostOption string

type ServerOptions struct {
	nsbHost NSBHostOption
	logger  logger.Logger
}

func defaultServerOptions() ServerOptions {
	return ServerOptions{
		logger:  logger.NewStdLogger(),
		nsbHost: "localhost:26657",
	}
}

func parseOptions(rOptions []interface{}) ServerOptions {
	var options = defaultServerOptions()
	for i := range rOptions {
		switch option := rOptions[i].(type) {
		case logger.Logger:
			options.logger = option
		case NSBHostOption:
			options.nsbHost = option
		}
	}
	return options
}

// NewServer return a pointer of Server
func NewServer(
	muldb types.MultiIndex,
	sindb types.Index,
	migrateFunction MigrateFunction,
	signer *signaturer.TendermintNSBSigner,
	rOptions ...interface{},
) (*Server, error) {
	var server = new(Server)
	options := parseOptions(rOptions)

	server.Logger = options.logger
	server.Signer = signer

	server.Resp = &uipbase.Account{Address: server.Signer.GetPublicKey(), ChainId: 3}
	server.Host = []byte{127, 0, 0, 1, ((23351) >> 8 & 0xff), 23351 & 0xff}
	err := migrate(muldb, migrateFunction)
	if err != nil {
		return nil, fmt.Errorf("failed to migrate: %v", err)
	}
	server.NsbHost, err = helper.HostFromString(string(options.nsbHost))
	if err != nil {
		return nil, fmt.Errorf("failed to decode host: %v", err)
	}
	// todo use instance
	server.DB = new(vesdb.Database)

	server.DB.SetMultiIndex(muldb)
	server.DB.SetIndex(sindb)

	server.DB.SetUserBase(new(user.XORMUserBase))
	server.DB.SetSessionBase(session.NewMultiThreadSerialSessionBase())
	server.DB.SetStorageHandler(new(storage_handler.Database))
	server.DB.SetSessionKVBase(new(kvdb.Database))
	server.DB.SetChainDNS(chain_dns.NewDatabase(config.HostMap))

	log.Println("will connect to remote nsb host", options.nsbHost)
	server.NsbClient = nsbcli.NewNSBClient(string(options.nsbHost))
	return server, nil
}

// ListenAndServe listen the port `port` and connect with remote central-ves with
// address `centerAddress`
func (server *Server) ListenAndServe(port, centerAddress string) error {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}
	// lis2, err := net.Listen("tcp", ":33334")
	// if err != nil {
	// 	return fmt.Errorf("failed to listen: %v", err)
	// }

	conn, err := grpc.Dial(centerAddress, grpc.WithInsecure(), grpc.WithKeepaliveParams(keepalive.ClientParameters{}))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	server.CVes = uiprpc.NewCenteredVESClient(conn)

	s := grpc.NewServer()

	uiprpc.RegisterVESServer(s, server)
	reflection.Register(s)

	fmt.Printf("prepare to serve on %v\n", port)

	go func() {
		if err = s.Serve(lis); err != nil {
			err = fmt.Errorf("failed to serve: %v", err)
			return
		}
	}()
	// fmt.Println(port)
	mux := GetHTTPServeMux(server)

	if err = mux.Run(":33004"); err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}
	return err
}

func GetHTTPServeMux(server *Server) *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOriginFunc:  func(origin string) bool { return true },
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "X-Grpc-Web", "X-User-Agent"},
		ExposeHeaders:    []string{"grpc-status", "grpc-message"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	vGroup := r.Group("/uiprpc.VES")

	tieMethod(
		vGroup, "UserRegister",
		func() proto.Message { return new(uiprpc.UserRegisterRequest) },
		func(ctx context.Context, in proto.Message) (proto.Message, error) {
			return server.UserRegister(context.Background(), in.(*uiprpc.UserRegisterRequest))
		},
	)

	tieMethod(
		vGroup, "SessionStart",
		func() proto.Message { return new(uiprpc.SessionStartRequest) },
		func(ctx context.Context, in proto.Message) (proto.Message, error) {
			return server.SessionStart(context.Background(), in.(*uiprpc.SessionStartRequest))
		},
	)

	tieMethod(
		vGroup, "SessionAckForInit",
		func() proto.Message { return new(uiprpc.SessionAckForInitRequest) },
		func(ctx context.Context, in proto.Message) (proto.Message, error) {
			return server.SessionAckForInit(context.Background(), in.(*uiprpc.SessionAckForInitRequest))
		},
	)

	tieMethod(
		vGroup, "SessionRequireTransact",
		func() proto.Message { return new(uiprpc.SessionRequireTransactRequest) },
		func(ctx context.Context, in proto.Message) (proto.Message, error) {
			return server.SessionRequireTransact(context.Background(), in.(*uiprpc.SessionRequireTransactRequest))
		},
	)

	tieMethod(
		vGroup, "SessionRequireRawTransact",
		func() proto.Message { return new(uiprpc.SessionRequireRawTransactRequest) },
		func(ctx context.Context, in proto.Message) (proto.Message, error) {
			return server.SessionRequireRawTransact(context.Background(), in.(*uiprpc.SessionRequireRawTransactRequest))
		},
	)

	tieMethod(
		vGroup, "AttestationReceive",
		func() proto.Message { return new(uiprpc.AttestationReceiveRequest) },
		func(ctx context.Context, in proto.Message) (proto.Message, error) {
			return server.AttestationReceive(context.Background(), in.(*uiprpc.AttestationReceiveRequest))
		},
	)

	tieMethod(
		vGroup, "MerkleProofReceive",
		func() proto.Message { return new(uiprpc.MerkleProofReceiveRequest) },
		func(ctx context.Context, in proto.Message) (proto.Message, error) {
			return server.MerkleProofReceive(context.Background(), in.(*uiprpc.MerkleProofReceiveRequest))
		},
	)

	tieMethod(
		vGroup, "ShrotenMerkleProofReceive",
		func() proto.Message { return new(uiprpc.ShortenMerkleProofReceiveRequest) },
		func(ctx context.Context, in proto.Message) (proto.Message, error) {
			return server.ShrotenMerkleProofReceive(context.Background(), in.(*uiprpc.ShortenMerkleProofReceiveRequest))
		},
	)

	tieMethod(
		vGroup, "InformMerkleProof",
		func() proto.Message { return new(uiprpc.MerkleProofReceiveRequest) },
		func(ctx context.Context, in proto.Message) (proto.Message, error) {
			return server.InformMerkleProof(context.Background(), in.(*uiprpc.MerkleProofReceiveRequest))
		},
	)

	tieMethod(
		vGroup, "InformShortenMerkleProof",
		func() proto.Message { return new(uiprpc.ShortenMerkleProofReceiveRequest) },
		func(ctx context.Context, in proto.Message) (proto.Message, error) {
			return server.InformShortenMerkleProof(context.Background(), in.(*uiprpc.ShortenMerkleProofReceiveRequest))
		},
	)

	tieMethod(
		vGroup, "InformAttestation",
		func() proto.Message { return new(uiprpc.AttestationReceiveRequest) },
		func(ctx context.Context, in proto.Message) (proto.Message, error) {
			return server.InformAttestation(context.Background(), in.(*uiprpc.AttestationReceiveRequest))
		},
	)

	return r
}

func tieMethod(
	vGroup *gin.RouterGroup, method string,
	objectFactory func() proto.Message,
	serveFunc func(ctx context.Context, in proto.Message) (proto.Message, error),
) {
	vGroup.OPTIONS(method, func(c *gin.Context) {
		c.Status(200)
	})

	vGroup.POST(method, func(c *gin.Context) {
		p, err := c.GetRawData()
		if err != nil && err != io.EOF {
			c.AbortWithError(400, err)
			return
		}
		var in = objectFactory()
		err = decode(p, in)
		if err != nil {
			c.AbortWithError(400, err)
			return
		}

		ret, err := serveFunc(context.Background(), in)
		if err != nil {
			c.AbortWithError(400, err)
			return
		}
		b, err := encode(c, ret)
		if err != nil {
			c.AbortWithError(500, err)
			return
		}
		c.Data(200, "application/grpc-web-text", b)
		// fmt.Println(c.Request.Method)
		// fmt.Println(c.Request.URL)
		// fmt.Println(c.Request.Proto)
		// fmt.Println(c.Request.ProtoMajor)
		// fmt.Println(c.Request.Header)
		// fmt.Println(c.Request.ProtoMajor)
	})
}

// now only decode one message
func decode(b []byte, in proto.Message) error {
	l := base64.StdEncoding.DecodedLen(len(b))
	if l < 5 {
		return errors.New("too short")
	}
	g := make([]byte, l)
	_, err := base64.StdEncoding.Decode(g, b)
	if err != nil {
		return err
	}

	if g[0] != 0 {
		return errors.New("needing data frame")
	}

	var gg = binary.BigEndian.Uint32(g[1:5])
	if len(g) != 5+int(gg) {
		return errors.New("malformed")
	}
	return proto.Unmarshal(g[5:5+gg], in)
}

func encode(c *gin.Context, ret proto.Message) (t []byte, err error) {
	s, err := proto.Marshal(ret)
	if err != nil {
		return nil, err
	}
	t = make([]byte, len(s)+5)
	t[0] = 0
	binary.BigEndian.PutUint32(t[1:5], uint32(len(s)))
	copy(t[5:], s)
	l := base64.StdEncoding.EncodedLen(len(t))
	g := make([]byte, l)
	base64.StdEncoding.Encode(g, t)
	return g, nil
}
