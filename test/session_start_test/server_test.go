package main

import (
	"fmt"
	"time"

	uiprpc "github.com/Myriad-Dreamin/go-ves/grpc/uiprpc"
	uipbase "github.com/Myriad-Dreamin/go-ves/grpc/uiprpc-base"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"testing"
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
		&uiprpc.UserRegisterRequest{Account: &uipbase.Account{
			ChainId: 1,
			Address: []byte{1},
		},
		})
	if err != nil {
		t.Fatalf("could not greet: %v", err)
	}
	fmt.Printf("Register: %v\n", r.Ok)
}
