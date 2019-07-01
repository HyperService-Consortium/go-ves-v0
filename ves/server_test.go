package ves

import (
	"fmt"
	"time"

	uiprpc "github.com/Myriad-Dreamin/go-ves/grpc"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"testing"
)

const (
	address     = "127.0.0.1:23351"
	defaultName = "world"
)

func TestUserRegister(t *testing.T) {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
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
