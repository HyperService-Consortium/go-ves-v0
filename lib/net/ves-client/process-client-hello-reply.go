package vesclient

import(
	"github.com/HyperService-Consortium/go-ves/grpc/wsrpc"
	helper "github.com/HyperService-Consortium/go-ves/lib/net/help-func"
	nsbcli "github.com/HyperService-Consortium/go-ves/lib/net/nsb-client"
)

func (vc *VesClient) ProcessClientHelloReply(clientHelloReply *wsrpc.ClientHelloReply) {
	var err error
	vc.grpcip, err = helper.DecodeIP(clientHelloReply.GetGrpcHost())
	if err != nil {
		vc.logger.Error("VesClient.read.ClientHelloReply.decodeGRPCHost", "error", err)
	} else {
		vc.logger.Info("adding default grpc ip ", "ip", vc.grpcip)
	}

	vc.nsbip, err = helper.DecodeIP(clientHelloReply.GetNsbHost())
	if err != nil {
		vc.logger.Error("VesClient.read.ClientHelloReply.decodeNSBHost", "error", err)
	} else {
		vc.logger.Info("adding default nsb ip ", "ip", vc.nsbip)
	}

	// todo: restrict scope
	vc.nsbClient = nsbcli.NewNSBClient(vc.nsbip)
}
