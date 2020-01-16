package bencher

import (
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"sync/atomic"
	"time"

	uiptypes "github.com/HyperService-Consortium/go-uip/uiptypes"
	nsbclient "github.com/HyperService-Consortium/go-ves/lib/net/nsb-client"

	signaturer "github.com/HyperService-Consortium/go-uip/signaturer"
)

type task struct {
	badSession  Int
	addedAction Int
	action      Int
	tps         Int
}

const signSize = 65
const validated = false

var idleSignature = make([]byte, signSize)

func (tasker *task) nsbRoutine(
	cli *nsbclient.NSBClient,
	SignContentSize, ActionLong int,
	batch bool,
	signer uiptypes.Signer, index int, txpadding []byte) {

	// info, err := cli.GetAbciInfo()
	iscAddress, err := cli.CreateISC(signer, []uint32{0}, [][]byte{signer.GetPublicKey()}, nil, txpadding)
	if err != nil {
		tasker.badSession.Add(1)
		fmt.Println(err)
		return
	}
	tasker.tps.Add(1)
	var badSessionFlag = false
	var localAddedAction int32

	var cont = make([]byte, SignContentSize-1, SignContentSize)

	var bytesIdx = make([]byte, 8)

	if batch {
		for idf := 0; idf < 60; idf++ {
			tasker.action.Add(int32(ActionLong))
			batcher := cli.AddActions(signer, nil, ActionLong)
			for idx := 0; idx < ActionLong; idx++ {
				binary.BigEndian.PutUint64(bytesIdx, uint64(idx+1))
				batcher.Insert(iscAddress, uint64(index), 0, 1, append(cont, bytesIdx...), idleSignature)
			}

			_, err := batcher.Commit()
			if err != nil {
				badSessionFlag = true
				fmt.Println(err)
			} else {
				tasker.tps.Add(1)
				localAddedAction += int32(ActionLong)
			}
		}
		// localAddedAction += int32(ActionLong)
	} else {
		for idf := 0; idf < 60; idf++ {
			tasker.action.Add(int32(ActionLong))
			for idx := 0; idx < ActionLong; idx++ {
				binary.BigEndian.PutUint64(bytesIdx, uint64(idx+1))
				_, err := cli.AddAction(signer, nil,
					iscAddress, uint64(index), 0, 1, append(cont, bytesIdx...), idleSignature)
				if err != nil {
					badSessionFlag = true
					fmt.Println(err)
				} else {
					localAddedAction++
					tasker.tps.Add(1)
				}
			}
		}
	}

	if badSessionFlag {
		tasker.badSession.Add(1)
	}

	tasker.addedAction.Add(localAddedAction)
}

// Header csv when file is being created
func Header(f *os.File) error {
	_, err := f.Write([]byte("action size, validated, bad session, total session, batched, added action, planning action, tps, upload (KB/s), download (KB/s), uploaded (KB), downloaded (KB), time used (s)\n"))
	return err
}

// Main bencher task
func Main(cli *nsbclient.NSBClient, SessionLimit, SignContentSize, ActionLong int, Batched bool, f *os.File) error {

	var _TxpaddingSize int

	var txpadding []byte
	_TxpaddingSize = 1

	txpadding = make([]byte, _TxpaddingSize)

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
	var tasker = new(task)

	for idx := 0; idx < SessionLimit; idx++ {
		time.Sleep(20 * time.Millisecond)
		go func(index int) {
			tasker.nsbRoutine(cli, SignContentSize, ActionLong, Batched, signer, index, txpadding)
			U <- true
		}(idx)
	}
	for idx := 0; idx < SessionLimit; idx++ {
		<-U
	}
	var consumed = time.Now().Sub(costing).Seconds()
	var base = 1024 * consumed
	fmt.Println("================================================================================")
	fmt.Printf(
		"\naction size: %v, validated: %v\n bad session count: %v/%v\n batched: %v, addedAction: %v/%v\n tps: %v\n UpLoaded: %vKB/s, Downloaded: %vKB/s\n UpLoaded: %vKB, Downloaded: %vKB, base %vs\n",
		32+8+8+1+SignContentSize+64, validated, tasker.badSession.value, SessionLimit, Batched, tasker.addedAction.value, tasker.action.value, float64(tasker.tps.value)/consumed,
		float64(nsbclient.SentBytes)/base, float64(nsbclient.ReceivedBytes)/base,
		float64(nsbclient.SentBytes)/1024.0, float64(nsbclient.ReceivedBytes)/1024.0,
		consumed,
	)

	if f != nil {
		f.Write([]byte(
			fmt.Sprintf(
				"%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v\n",
				32+8+8+1+SignContentSize+64, validated, tasker.badSession.value, SessionLimit, Batched, tasker.addedAction.value, tasker.action.value, float64(tasker.tps.value)/consumed,
				float64(nsbclient.SentBytes)/base, float64(nsbclient.ReceivedBytes)/base,
				float64(nsbclient.SentBytes)/1024.0, float64(nsbclient.ReceivedBytes)/1024.0,
				consumed,
			),
		))
	}
	return nil
}

// MaybeWithFile Help OpenFile
func MaybeWithFile(filePath string, Task func(f *os.File) error, IfNewFile func(f *os.File) error) error {
	if filePath != "" {
		var f *os.File
		var err error
		var nf = false
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			f, err = os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
			nf = true
		} else if err != nil {
			log.Fatal(err)
			return nil
		} else {
			f, err = os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND, 0666)
		}
		if err != nil {
			log.Fatal(err)
			return nil
		}
		// will close whether Task panic or not
		defer f.Close()
		if nf {
			if err = IfNewFile(f); err != nil {
				return err
			}
		}
		err = Task(f)
		return err
	}
	return Task(nil)

}

type Int struct {
	value int32
}

func (i *Int) Add(inc int32) {
	atomic.AddInt32(&i.value, inc)
}
