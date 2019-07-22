package main

import (
	uiprpc "github.com/Myriad-Dreamin/go-ves/grpc/uiprpc"
	"golang.org/x/net/context"
)

const port = ":23351"

type Server struct {
}

func (server *Server) SayHello(
	ctx context.Context,
	in *uiprpc.UserRegisterRequest,
) (*uiprpc.UserRegisterReply, error) {
	return &uiprpc.UserRegisterReply{Ok: true}, nil
}

func main() {

}

// func main() {
// 	lis, err := net.Listen("tcp", port)
// 	if err != nil {
// 		log.Fatalf("failed to listen: %v", err)
// 	}
//
// 	s := grpc.NewServer()
//
// 	uiprpc.RegisterVESServer(s, &Server{})
// 	reflection.Register(s)
// 	if err := s.Serve(lis); err != nil {
// 		log.Fatalf("failed to serve: %v", err)
// 	}
// }
