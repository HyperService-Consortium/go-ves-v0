package vesclient

import (
	"net/url"

	uiptypes "github.com/HyperService-Consortium/go-uip/types"
	"github.com/gorilla/websocket"
)

// VanilleMakeClient for humans
func VanilleMakeClient(name, addr string) (*VesClient, error) {
	var (
		dialer        *websocket.Dialer
		vcClient, err = NewVesClient()
	)
	if err != nil {
		return nil, err
	}

	vcClient.waitOpt = uiptypes.NewWaitOption()

	vcClient.name = []byte(name)

	if err = vcClient.load(dataPrefix + "/" + string(vcClient.name)); err != nil {
		return nil, err
	}
	phandler.register(vcClient.save)

	if vcClient.conn, _, err = dialer.Dial(
		(&url.URL{Scheme: "ws", Host: addr, Path: "/"}).String(), nil,
	); err != nil {
		return nil, err
	}

	if err = vcClient.sayClientHello(vcClient.name); err != nil {
		return nil, err
	}

	return vcClient, nil
}
