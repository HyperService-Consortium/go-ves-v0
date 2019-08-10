package ves

import (
	"fmt"
	"net"

	signaturer "github.com/Myriad-Dreamin/go-uip/signaturer"
	uiprpc "github.com/Myriad-Dreamin/go-ves/grpc/uiprpc"
	uipbase "github.com/Myriad-Dreamin/go-ves/grpc/uiprpc-base"
	log "github.com/Myriad-Dreamin/go-ves/lib/log"
	types "github.com/Myriad-Dreamin/go-ves/types"
	vesdb "github.com/Myriad-Dreamin/go-ves/types/database"
	kvdb "github.com/Myriad-Dreamin/go-ves/types/kvdb"
	session "github.com/Myriad-Dreamin/go-ves/types/session"
	user "github.com/Myriad-Dreamin/go-ves/types/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

// Server provides the basic service of session
type Server struct {
	db     types.VESDB
	resp   *uipbase.Account
	signer *signaturer.TendermintNSBSigner
	cves   uiprpc.CenteredVESClient
	// mutex sync.Mutex
	// mup  map[uint16]bool
}

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

// NewServer return a pointer of Server
func NewServer(
	muldb types.MultiIndex,
	sindb types.Index,
	migrateFunction MigrateFunction,
	signer *signaturer.TendermintNSBSigner,
) (*Server, error) {
	var server = new(Server)

	server.signer = signer
	server.resp = &uipbase.Account{Address: server.signer.GetPublicKey(), ChainId: 3}

	err := migrate(muldb, migrateFunction)
	if err != nil {
		return nil, fmt.Errorf("failed to migrate: %v", err)
	}

	server.db = new(vesdb.Database)

	server.db.SetMultiIndex(muldb)
	server.db.SetIndex(sindb)

	server.db.SetUserBase(new(user.XORMUserBase))
	server.db.SetSessionBase(session.NewMultiThreadSerialSessionBase())
	server.db.SetSessionKVBase(new(kvdb.Database))

	return server, nil
}

// ListenAndServe listen the port `port` and connect with remote central-ves with
// address `centerAddress`
func (server *Server) ListenAndServe(port, centerAddress string) error {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	conn, err := grpc.Dial(centerAddress, grpc.WithInsecure(), grpc.WithKeepaliveParams(keepalive.ClientParameters{}))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	server.cves = uiprpc.NewCenteredVESClient(conn)

	s := grpc.NewServer()

	uiprpc.RegisterVESServer(s, server)
	reflection.Register(s)

	fmt.Printf("prepare to serve on %v\n", port)

	if err := s.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}

	return nil
}
