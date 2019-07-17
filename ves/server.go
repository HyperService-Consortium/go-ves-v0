package ves

import (
	"fmt"
	"net"

	signaturer "github.com/Myriad-Dreamin/go-uip/signaturer"
	multi_index "github.com/Myriad-Dreamin/go-ves/database/multi_index"
	uiprpc "github.com/Myriad-Dreamin/go-ves/grpc/uip-rpc"
	types "github.com/Myriad-Dreamin/go-ves/types"
	vesdb "github.com/Myriad-Dreamin/go-ves/types/database"
	session "github.com/Myriad-Dreamin/go-ves/types/session"
	user "github.com/Myriad-Dreamin/go-ves/types/user"
	service "github.com/Myriad-Dreamin/go-ves/ves/service"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	migrate_function = XORMMigrate
)

type Server struct {
	db     types.VESDB
	signer *signaturer.TendermintNSBSigner
}

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

func (server *Server) migrate(muldb types.MultiIndex, mfunc func(types.MultiIndex) error) error {
	return mfunc(muldb)
}

func (server *Server) UserRegister(
	ctx context.Context,
	in *uiprpc.UserRegisterRequest,
) (*uiprpc.UserRegisterReply, error) {
	return service.UserRegisterService{
		VESDB:               server.db,
		Context:             ctx,
		UserRegisterRequest: in,
	}.Serve()
}

func (server *Server) SessionStart(
	ctx context.Context,
	in *uiprpc.SessionStartRequest,
) (*uiprpc.SessionStartReply, error) {
	return (&service.SessionStartService{
		Signer:              server.signer,
		VESDB:               server.db,
		Context:             ctx,
		SessionStartRequest: in,
	}).Serve()
}

func (server *Server) SessionAckForInit(
	ctx context.Context,
	in *uiprpc.SessionAckForInitRequest,
) (*uiprpc.SessionAckForInitReply, error) {
	return service.SessionAckForInitService{
		VESDB:                    server.db,
		Context:                  ctx,
		SessionAckForInitRequest: in,
	}.Serve()
}

func (server *Server) SessionRequireTransact(
	ctx context.Context,
	in *uiprpc.SessionRequireTransactRequest,
) (*uiprpc.SessionRequireTransactReply, error) {
	return service.SessionRequireTransactService{
		VESDB:                         server.db,
		Context:                       ctx,
		SessionRequireTransactRequest: in,
	}.Serve()
}

func (server *Server) AttestationReceive(
	ctx context.Context,
	in *uiprpc.AttestationReceiveRequest,
) (*uiprpc.AttestationReceiveReply, error) {
	return service.AttestationReceiveService{
		VESDB:                     server.db,
		Context:                   ctx,
		AttestationReceiveRequest: in,
	}.Serve()
}

func ListenAndServe(port string) error {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	var server = new(Server)
	server.db = new(vesdb.Database)
	server.signer = signaturer.NewTendermintNSBSigner([]byte{
		233, 66, 233, 66, 233, 66, 233, 66, 233, 66, 233, 66, 233, 66, 233, 66,
		233, 66, 233, 66, 233, 66, 233, 66, 233, 66, 233, 66, 233, 66, 233, 66,
		233, 66, 233, 66, 233, 66, 233, 66, 233, 66, 233, 66, 233, 66, 233, 66,
		233, 66, 233, 66, 233, 66, 233, 66, 233, 66, 233, 66, 233, 66, 233, 66,
	})

	//TODO: SetEnv
	var muldb *multi_index.XORMMultiIndexImpl
	muldb, err = multi_index.GetXORMMultiIndex("mysql", "ves:123456@tcp(127.0.0.1:3306)/ves?charset=utf8")
	if err != nil {
		return fmt.Errorf("failed to get muldb: %v", err)
	}
	err = server.migrate(muldb, migrate_function)
	if err != nil {
		return fmt.Errorf("failed to migrate: %v", err)
	}

	server.db.SetMultiIndex(muldb)

	server.db.SetUserBase(new(user.XORMUserBase))
	server.db.SetSessionBase(new(session.SerialSessionBase))

	s := grpc.NewServer()

	uiprpc.RegisterVESServer(s, server)
	reflection.Register(s)

	fmt.Printf("prepare to serve on %v\n", port)

	if err := s.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}

	return nil
}
