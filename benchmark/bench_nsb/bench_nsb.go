package main

import (
	"encoding/binary"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"sync/atomic"
	"time"

	uiptypes "github.com/HyperService-Consortium/go-uip/uiptypes"
	nsbclient "github.com/Myriad-Dreamin/go-ves/lib/net/nsb-client"

	signaturer "github.com/HyperService-Consortium/go-uip/signaturer"
)

var (
	// host = flag.String("host", "127.0.0.1:26657", "aim nsb sever")
	host = "47.251.2.73:26657"
	cli  *nsbclient.NSBClient

	badSession = new(Int)
)

var SessionLimit int
var sessionLimit = flag.Int("ses", 1, "max count of go-routine")

var SignContentSize int
var signContentSize = flag.Int("con", 400, "signature content size")

var OffchainTransactionSize int
var offchainTransactionSize = flag.Int("oc", 200, "off-chain transaction size(in op-intent)")

var NodeSize int
var nodeSize = flag.Int("node.siz", 200, "average size of merkle proof nodes")

var ProofDepth int
var proofDepth = flag.Int("node.dep", 4, "arverage depth of merkle proof nodes")

var AverageCountOfTxInEachOpIntent int
var averageCountOfTxInEachOpIntent = flag.Int("txcount", 1, "arverage count of tx in each op-intent")

const signSize = 65
const hashSize = 64

var ProofSize, TxIntentSize, TxpaddingSize int

var bb, bg, idleProof, txpadding []byte

func NSBRoutine(signer uiptypes.Signer, index int) {
	// info, err := cli.GetAbciInfo()
	iscAddress, err := cli.CreateISC(signer, []uint32{0}, [][]byte{signer.GetPublicKey()}, nil, txpadding)
	if err != nil {
		badSession.Add(1)
		fmt.Println(err)
		return
	}
	_ = iscAddress
	var bbg, cont = make([]byte, 8), make([]byte, SignContentSize-1, SignContentSize)
	for idx := 0; idx < 20; idx++ {
		_, err := cli.AddAction(signer, nil,
			iscAddress, uint64(index), 0, 1, append(cont, uint8(1+idx)), bb)
		if err != nil {
			badSession.Add(1)
			fmt.Println(err)
			return
		}
	}

	for idx := 0; idx < 5; idx++ {
		binary.BigEndian.PutUint64(bbg, uint64(index<<8|idx))
		_, err = cli.AddMerkleProof(signer, nil,
			1, bg, idleProof, bbg, bg,
		)
		if err != nil {
			badSession.Add(1)
			fmt.Println(err)
			return
		}

		_, err = cli.AddBlockCheck(signer, nil,
			1, bbg, bg, 1,
		)
		if err != nil {
			badSession.Add(1)
			fmt.Println(err)
			return
		}
	}

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
	}

	var U = make(chan bool, SessionLimit)
	for idx := 0; idx < SessionLimit; idx++ {
		time.Sleep(20 * time.Millisecond)
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
		"bad session count: %v/%v\n UpLoaded: %vKB/s, Downloaded: %vKB/s\n UpLoaded: %vKB, Downloaded: %vKB, base %vs\n",
		badSession.value, SessionLimit,
		float64(nsbclient.SentBytes)/base, float64(nsbclient.ReceivedBytes)/base,
		float64(nsbclient.SentBytes)/1024.0, float64(nsbclient.ReceivedBytes)/1024.0,
		consumed,
	)

}

func init() {
	flag.Parse()
	SessionLimit = *sessionLimit
	SignContentSize = *signContentSize
	ProofDepth = *proofDepth
	NodeSize = *nodeSize
	OffchainTransactionSize = *offchainTransactionSize
	AverageCountOfTxInEachOpIntent = *averageCountOfTxInEachOpIntent

	ProofSize = hashSize*ProofDepth + NodeSize*(ProofDepth-1)
	const chainidSize, numberSize = 8, 32
	TxIntentSize = hashSize*4 + chainidSize*2 + numberSize*2 + OffchainTransactionSize
	TxpaddingSize = AverageCountOfTxInEachOpIntent * TxIntentSize

	// fmt.Println(signSize, hashSize, ProofSize, TxpaddingSize)

	bb = make([]byte, signSize)
	bg = bb[0:hashSize]
	idleProof, _ = json.Marshal(map[string]interface{}{
		"r": bg,
		"h": [][]byte{make([]byte, ProofSize)},
	})
	txpadding = make([]byte, TxpaddingSize)

	cli = nsbclient.NewNSBClient(host)
}

type Int struct {
	value int32
}

func (i *Int) Add(inc int32) {
	atomic.AddInt32(&i.value, inc)
}
