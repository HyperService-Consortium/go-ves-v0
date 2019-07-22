package main

import (
	"fmt"
	"time"

	uiprpc "github.com/Myriad-Dreamin/go-ves/grpc/uiprpc"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"testing"
)

const (
	m_address = "127.0.0.1:23351"
)

func TestUserRegister(t *testing.T) {
	go main()

	// Set up a connection to the server.
	conn, err := grpc.Dial(m_address, grpc.WithInsecure())
	if err != nil {
		t.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := uiprpc.NewVESClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.UserRegister(
		ctx,
		&uiprpc.UserRegisterRequest{User: &uiprpc.Account{
			ChainType: 1,
			Address:   []byte{1},
		},
		})
	if err != nil {
		t.Fatalf("could not greet: %v", err)
	}
	fmt.Printf("Register: %v\n", r.Ok)
}
