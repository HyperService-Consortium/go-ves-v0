package ves

import (
	"encoding/hex"
	"fmt"
	"log"
	"time"

	signaturer "github.com/HyperService-Consortium/go-uip/signaturer"
	uiprpc "github.com/HyperService-Consortium/go-ves/grpc/uiprpc"
	uipbase "github.com/HyperService-Consortium/go-ves/grpc/uiprpc-base"
	index "github.com/HyperService-Consortium/go-ves/lib/database/index"
	multi_index "github.com/HyperService-Consortium/go-ves/lib/database/multi_index"
	"github.com/gogo/protobuf/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"testing"
)

const (
	mPort         = ":23351"
	mAddress      = "127.0.0.1:23351"
	centerAddress = "127.0.0.1:23452"
)

func TestUserRegister(t *testing.T) {
	var userReq uiprpc.UserRegisterRequest
	userReq.Account = new(uipbase.Account)
	userReq.Account.Address = []byte{0x01, 0x23}
	userReq.Account.ChainId = 1
	userReq.UserName = "000123"
	fmt.Println(proto.Marshal(&userReq))

	var err error

	//TODO: SetEnv
	var muldb *multi_index.XORMMultiIndexImpl
	muldb, err = multi_index.GetXORMMultiIndex("mysql", "ves:123456@tcp(127.0.0.1:3306)/ves?charset=utf8")
	if err != nil {
		t.Errorf("failed to get muldb: %v", err)
		return
	}
	var sindb *index.LevelDBIndex
	sindb, err = index.GetIndex("./data")
	if err != nil {
		t.Errorf("failed to get sindb: %v", err)
		return
	}

	b, err := hex.DecodeString("2333bbffffffffffffff2333bbffffffffffffff2333bbffffffffffffffffff2333bbffffffffffffff2333bbffffffffffffff2333bbffffffffffffffffff")
	if err != nil {
		t.Error(err)
		return
	}
	signer, err := signaturer.NewTendermintNSBSigner(b)
	if err != nil {
		t.Error(err)
		return
	}

	var server *Server
	if server, err = NewServer(
		muldb, sindb, multi_index.XORMMigrate, signer,
	); err != nil {
		log.Fatal(err)
	}
	go func() {
		if err := server.ListenAndServe(mAddress, centerAddress); err != nil {
			log.Fatal(err)
		}
	}()

	// Set up a connection to the server.
	conn, err := grpc.Dial(mAddress, grpc.WithInsecure())

	if err != nil {
		t.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := uiprpc.NewVESClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
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
