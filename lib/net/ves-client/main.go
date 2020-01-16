package vesclient

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/Myriad-Dreamin/minimum-lib/logger"
	"os"
)

var (
	addr = flag.String("addr", "localhost:23452", "http service address")
)

func init() {
	go phandler.atExit()
	if !flag.Parsed() {
		flag.Parse()
	}
}

// Main is the origin main of ves client
func Main() {
	vcClient, err := NewVesClient(logger.NewZapDevelopmentSugarOption(), CVesHostOption(*addr))
	if err != nil {
		globalLogger.Fatal("get ves client error", "error", err)
	}

	fmt.Println("input your name:")
	vcClient.name, _, err = bufio.NewReader(os.Stdin).ReadLine()
	if err != nil {
		vcClient.logger.Fatal("set name", "error", err)
	}

	if err = vcClient.Boot(); err != nil {
		os.Exit(1)
	}

	phandler.register(func() { vcClient.quit <- true })
	go vcClient.write()
	select {
	case <-vcClient.quit:
		return
	}
}
