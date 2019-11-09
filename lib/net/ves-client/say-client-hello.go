package vesclient

import "github.com/HyperService-Consortium/go-ves/grpc/wsrpc"

func (vc *VesClient) SayClientHello(name []byte) error {
	clientHello := vc.getClientHello()
	clientHello.Name = name

	err := vc.postMessage(wsrpc.CodeClientHelloRequest, clientHello)
	if err != nil {
		vc.logger.Error("say client hello", "name", name, "error", err)
		return err
	}
	vc.logger.Info("say client hello to server successfully")
	return nil
}
