package main

import (
	"flag"
	"os"

	bencher "github.com/HyperService-Consortium/go-ves/benchmark/bench_nsb_action/bencher"

	nsbclient "github.com/HyperService-Consortium/go-ves/lib/net/nsb-client"
)

var (
	host = flag.String("host", "127.0.0.1:26657", "aim nsb sever")
	cli  *nsbclient.NSBClient

	_SessionLimit int
	sessionLimit  = flag.Int("ses", 1, "max count of go-routine")

	_SignContentSize int
	signContentSize  = flag.Int("con", 400, "signature content size")

	_ActionLong int
	actionLong  = flag.Int("accs", 20, "number of actions each routine must add")

	_Batched bool
	batched  = flag.Bool("batch", false, "use addActions or addaAtion")

	_WithFileToAppend string
	withFileToAppend  = flag.String("o", "", "to append csv files")
)

func main() {
	bencher.MaybeWithFile(_WithFileToAppend, func(f *os.File) error {
		return bencher.Main(cli, _SessionLimit, _SignContentSize, _ActionLong, _Batched, f)
	}, bencher.Header)
}

func init() {
	flag.Parse()
	_SessionLimit = *sessionLimit
	_SignContentSize = *signContentSize
	_ActionLong = *actionLong
	_Batched = *batched
	_WithFileToAppend = *withFileToAppend

	cli = nsbclient.NewNSBClient(*host)
}
