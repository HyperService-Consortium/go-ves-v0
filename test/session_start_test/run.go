package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"time"

	uiprpc "github.com/HyperService-Consortium/go-ves/grpc/uiprpc"
	uipbase "github.com/HyperService-Consortium/go-ves/grpc/uiprpc-base"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	m_port    = ":23351"
	m_address = "127.0.0.1:23351"
)

type obj map[string]interface{}

func main() {
	var opintent = obj{
		"name":    "Op1",
		"op_type": "Payment",
		"src": obj{
			"domain":    2,
			"user_name": "a1",
		},
		"dst": obj{
			"domain":    1,
			"user_name": "a2",
		},
		"amount": "0x2e0",
		"unit":   "wei",
	}

	// Set up a connection to the server.
	conn, err := grpc.Dial(m_address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := uiprpc.NewVESClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	var b []byte
	b, err = json.Marshal(opintent)
	if err != nil {
		log.Fatalf("Marshal failed: %v", err)
	}
	fmt.Println(string(b))
	r, err := c.SessionStart(
		ctx,
		&uiprpc.SessionStartRequest{
			Opintents: &uipbase.OpIntents{
				Dependencies: nil,
				Contents: [][]byte{
					b,
				},
			},
		})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	fmt.Printf("Session Start: %v, %v\n", r.GetOk(), hex.EncodeToString(r.GetSessionId()))
}
