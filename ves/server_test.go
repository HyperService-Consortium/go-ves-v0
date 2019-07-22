package ves

import (
	"fmt"
	"log"
	"time"

	uiprpc "github.com/Myriad-Dreamin/go-ves/grpc/uiprpc"
	uipbase "github.com/Myriad-Dreamin/go-ves/grpc/uiprpc-base"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"testing"
)

const (
	m_port        = ":23351"
	m_address     = "127.0.0.1:23351"
	centerAddress = "127.0.0.1:23452"
)

func TestUserRegister(t *testing.T) {
	go func() {
		if err := ListenAndServe(m_port, centerAddress); err != nil {
			log.Fatal(err)
		}
	}()

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
