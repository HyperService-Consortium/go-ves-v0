package main

import (
	"fmt"
	"log"
	"time"

	uiprpc "github.com/Myriad-Dreamin/go-ves/grpc"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	m_port    = ":23351"
	m_address = "127.0.0.1:23351"
)

func main() {

	// Set up a connection to the server.
	conn, err := grpc.Dial(m_address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := uiprpc.NewVESClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.UserRegister(
		ctx,
		&uiprpc.UserRegisterRequest{Account: &uiprpc.Account{
			ChainId: 1,
			Address: []byte{1},
		},
		})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	fmt.Printf("Register: %v\n", r.Ok)
}
