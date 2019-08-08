package vesclient

import (
	"bufio"
	"flag"
	"fmt"
	"net/url"
	"os"

	"github.com/gorilla/websocket"

	uiptypes "github.com/Myriad-Dreamin/go-uip/types"
	log "github.com/Myriad-Dreamin/go-ves/lib/log"
)

// Main is the origin main of ves client
func Main() {
	var (
		dialer        *websocket.Dialer
		addr          = flag.String("addr", "localhost:23452", "http service address")
		u             = url.URL{Scheme: "ws", Host: *addr, Path: "/"}
		vcClient, err = NewVesClient()
	)

	vcClient.waitOpt = uiptypes.NewWaitOption()

	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println("input your name:")

	vcClient.name, _, err = bufio.NewReader(os.Stdin).ReadLine()

	if err != nil {
		log.Println(err)
		return
	}

	if err = vcClient.load(dataPrefix + "/" + string(vcClient.name)); err != nil {
		log.Println(err)
		return
	}
	phandler.register(vcClient.save)

	vcClient.conn, _, err = dialer.Dial(u.String(), nil)
	if err != nil {
		log.Println(err)
		return
	}
	go phandler.atExit()
	go vcClient.read()

	vcClient.sayClientHello(vcClient.name)

	go vcClient.write()

	phandler.register(func() { vcClient.quit <- true })
	// close
	select {
	case <-vcClient.quit:
		return
	}
}
