package vesclient

import (
	"encoding/hex"
	"fmt"
	"github.com/HyperService-Consortium/go-ves/grpc/wsrpc"
	helper "github.com/HyperService-Consortium/go-ves/lib/net/help-func"
	nsbcli "github.com/HyperService-Consortium/go-ves/lib/net/nsb-client"
)


func (vc *VesClient) ProcessRequestComingRequest(requestComingRequest *wsrpc.RequestComingRequest) {

	vc.logger.Info(
		"new session request coming",
		"session id", hex.EncodeToString(requestComingRequest.GetSessionId()),
		"responsible address", hex.EncodeToString(requestComingRequest.GetAccount().GetAddress()),
	)

	fmt.Printf(
		"client: new session request coming, session id: %v, responsible address: %v\n",
		hex.EncodeToString(requestComingRequest.GetSessionId()),
		hex.EncodeToString(requestComingRequest.GetAccount().GetAddress()),
	)


	signer, err := vc.getNSBSigner()
	if err != nil {
		vc.logger.Error("VesClient.read.RequestComingRequest.getNSBSigner", "error", err)
		return
	}

	hs, err := helper.DecodeIP(requestComingRequest.GetNsbHost())
	if err != nil {
		vc.logger.Error("VesClient.read.RequestComingRequest.DecodeIP", "error", err)
		return
	}

	// todo: new nsbclient
	if ret, err := nsbcli.NewNSBClient(hs).UserAck(
		signer,
		requestComingRequest.GetSessionId(),
		requestComingRequest.GetAccount().GetAddress(),
		// todo: signature
		[]byte("123"),
	); err != nil {
		vc.logger.Error("VesClient.read.RequestComingRequest.UserAck", "ip", hs, "error", err)
		return
	} else {
		vc.logger.Info(
			"user ack to nsb",
			"info", ret.Info, "data", string(ret.Data), "log", ret.Log, "tags", ret.Tags,
		)
		fmt.Printf(
			"client: user ack to nsb, address: %v\n",
			hex.EncodeToString(requestComingRequest.GetAccount().GetAddress()),
		)
	}

	if err = vc.sendAck(
		requestComingRequest.GetAccount(),
		requestComingRequest.GetSessionId(),
		requestComingRequest.GetGrpcHost(),
		signer.Sign(requestComingRequest.GetSessionId()).Bytes(),
	); err != nil {
		vc.logger.Error("VesClient.read.RequestComingRequest.sendAck", "error", err)
		return
	} else {
		vc.logger.Info(
			"user ack to server",
		)
		fmt.Printf(
			"client: user ack to server, address: %v\n",
			hex.EncodeToString(requestComingRequest.GetAccount().GetAddress()),
		)
	}

}
