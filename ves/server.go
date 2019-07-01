package ves

import (
	"fmt"
	"net"

	uiprpc "github.com/Myriad-Dreamin/go-ves/grpc"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
}

func (server *Server) UserRegister(
	ctx context.Context,
	in *uiprpc.UserRegisterRequest,
) (*uiprpc.UserRegisterReply, error) {
	return &uiprpc.UserRegisterReply{Ok: true}, nil
}

func (server *Server) SessionStart(
	ctx context.Context,
	in *uiprpc.SessionStartRequest,
) (*uiprpc.SessionStartReply, error) {
	return &uiprpc.SessionStartReply{Ok: true}, nil
}

func (server *Server) SessionAckForInit(
	ctx context.Context,
	in *uiprpc.SessionAckForInitRequest,
) (*uiprpc.SessionAckForInitReply, error) {
	return &uiprpc.SessionAckForInitReply{Ok: true}, nil
}

func (server *Server) SessionRequireTransact(
	ctx context.Context,
	in *uiprpc.SessionRequireTransactRequest,
) (*uiprpc.SessionRequireTransactReply, error) {
	return &uiprpc.SessionRequireTransactReply{Ok: true}, nil
}

func (server *Server) AttestationReceive(
	ctx context.Context,
	in *uiprpc.AttestationReceiveRequest,
) (*uiprpc.AttestationReceiveReply, error) {
	return &uiprpc.AttestationReceiveReply{Ok: true}, nil
}

func ListenAndServe(port string) error {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	uiprpc.RegisterVESServer(s, &Server{})
	reflection.Register(s)

	fmt.Printf("prepare to serve on %v\n", port)

	if err := s.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}

	return nil
}
