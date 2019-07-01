package centeredves

import (
	"log"
	"net"

	uiprpc "github.com/Myriad-Dreamin/go-ves/grpc"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const port = ":23452"

type Server struct {
}

func (server *Server) SayHello(
	ctx context.Context,
	in *uiprpc.HelloRequest,
) (*uiprpc.HelloReply, error) {
	return &uiprpc.HelloReply{Message: "Hello " + in.Name}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	uiprpc.RegisterGreeterServer(s, &Server{})
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
