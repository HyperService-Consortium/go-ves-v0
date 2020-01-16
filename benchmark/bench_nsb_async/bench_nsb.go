package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"sync/atomic"
	"time"

	uiptypes "github.com/HyperService-Consortium/go-uip/uiptypes"
	nsbclient "github.com/HyperService-Consortium/go-ves/lib/net/nsb-client"

	signaturer "github.com/HyperService-Consortium/go-uip/signaturer"
)

var (
	// host = flag.String("host", "127.0.0.1:26657", "aim nsb sever")
	host = "47.251.2.73:26657"
	cli  *nsbclient.NSBClient

	badSession = new(Int)
)

type Int struct {
	value int32
}

func (i *Int) Add(inc int32) {
	atomic.AddInt32(&i.value, inc)
}

const SessionLimit = 1

var bb = make([]byte, 65)
var bg = bb[0:64]
var idleProof, _ = json.Marshal(map[string]interface{}{
	"r": bg,
	"h": [][]byte{make([]byte, 64*8+400*7)},
})

var txpadding = make([]byte, 5*(64*5+100))

func NSBRoutine(signer uiptypes.Signer, index int) {
	// info, err := cli.GetAbciInfo()
	iscAddress, err := cli.CreateISCAsync(signer, []uint32{0}, [][]byte{signer.GetPublicKey()}, nil, txpadding, nil)
	if err != nil {
		badSession.Add(1)
		fmt.Println(err)
		return
	}
	_ = iscAddress
	// var bbg = make([]byte, 8)
	// for idx := 0; idx < 20; idx++ {
	// _, err := cli.AddAction(signer, nil,
	// 	iscAddress, uint64(index), 0, 1, []byte{uint8(1 + idx)}, bb, nil)
	// if err != nil {
	// 	badSession.Add(1)
	// 	fmt.Println(err)
	// 	return
	// }
	// }

	// for idx := 0; idx < 5; idx++ {
	// 	binary.BigEndian.PutUint64(bbg, uint64(index<<8|idx))
	// 	_, err = cli.AddMerkleProof(signer, nil,
	// 		1, bg, idleProof, bbg, bg,
	// 	)
	// 	if err != nil {
	// 		badSession.Add(1)
	// 		fmt.Println(err)
	// 		return
	// 	}

	// 	_, err = cli.AddBlockCheck(signer, nil,
	// 		1, bbg, bg, 1,
	// 	)
	// 	if err != nil {
	// 		badSession.Add(1)
	// 		fmt.Println(err)
	// 		return
	// 	}
	// }

	// AddMerkleProof
}

func main() {

	costing := time.Now()

	var privatekey = make([]byte, 64)
	for i := 0; i < 64; i++ {
		privatekey[i] = uint8(i)
	}

	signer, err := signaturer.NewTendermintNSBSigner(privatekey)
	if err != nil {
		log.Fatal(err)
		return
	}

	var U = make(chan bool, SessionLimit)
	for idx := 0; idx < SessionLimit; idx++ {
		go func(index int) {
			NSBRoutine(signer, index)
			U <- true
		}(idx)
	}
	for idx := 0; idx < SessionLimit; idx++ {
		<-U
	}

	var consumed = time.Now().Sub(costing).Seconds()
	var base = 1024 * consumed
	fmt.Printf(
		"bad session count: %v\n UpLoaded: %vKB/s, Downloaded: %vKB/s\n UpLoaded: %vKB, Downloaded: %vKB, base %vs\n",
		badSession.value,
		float64(nsbclient.SentBytes)/base, float64(nsbclient.ReceivedBytes)/base,
		float64(nsbclient.SentBytes)/1024.0, float64(nsbclient.ReceivedBytes)/1024.0,
		consumed,
	)

}

func init() {
	flag.Parse()
	cli = nsbclient.NewNSBClient(host)

}
