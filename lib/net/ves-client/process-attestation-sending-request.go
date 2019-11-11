package vesclient

import (
	"encoding/hex"
	"fmt"
	TxState "github.com/HyperService-Consortium/go-uip/const/transaction_state_type"
	"github.com/HyperService-Consortium/go-uip/uiptypes"
	"github.com/HyperService-Consortium/go-ves/grpc/uiprpc-base"
	"github.com/HyperService-Consortium/go-ves/grpc/wsrpc"
	helper "github.com/HyperService-Consortium/go-ves/lib/net/help-func"
	nsbcli "github.com/HyperService-Consortium/go-ves/lib/net/nsb-client"
)

func (vc *VesClient) ProcessAttestationSendingRequest(attestationSendingRequest *wsrpc.RequestComingRequest) {
	vc.logger.Info(
		"new transaction's attestation must be created",
		"session id", hex.EncodeToString(attestationSendingRequest.GetSessionId()),
		"address", hex.EncodeToString(attestationSendingRequest.GetAccount().GetAddress()),
	)

	transactionReply, err := vc.GetRawTransaction(
		attestationSendingRequest.GetSessionId(),
		attestationSendingRequest.GetGrpcHost(),
	)
	if err != nil {
		vc.logger.Error("VesClient.read.AttestationSendingRequest.getRawTransaction", "error", err)
		return
	}

	vc.logger.Info(
		"the instance of the transaction intent is",
		"tid", transactionReply.Tid,
		"tx", hex.EncodeToString(transactionReply.RawTransaction),
	)

	fmt.Printf("client: %v: new transaction's attestation must be created, session id: %v, tid: %v\n",
		hex.EncodeToString(attestationSendingRequest.GetSessionId()),
		hex.EncodeToString(attestationSendingRequest.GetAccount().GetAddress()),
		transactionReply.Tid,
		)

	signer, err := vc.getNSBSigner()
	if err != nil {
		vc.logger.Error("VesClient.read.AttestationSendingRequest.getNSBSigner", "error", err)
		return
	}

	hs, err := helper.DecodeIP(attestationSendingRequest.GetNsbHost())
	if err != nil {
		vc.logger.Error("VesClient.read.AttestationSendingRequest.DecodeIP", "error", err)
		return
	}

	// packet attestation
	var sendingAtte = vc.getReceiveAttestationReceiveRequest()
	sendingAtte.SessionId = attestationSendingRequest.GetSessionId()
	sendingAtte.GrpcHost = attestationSendingRequest.GetGrpcHost()

	sigg := signer.Sign(transactionReply.RawTransaction)
	sendingAtte.Atte = &uiprpc_base.Attestation{
		Tid:     transactionReply.Tid,
		Aid:     TxState.Instantiating,
		Content: transactionReply.RawTransaction,
		Signatures: append(make([]*uiprpc_base.Signature, 0, 1), &uiprpc_base.Signature{
			// todo use src.signer to sign
			SignatureType: uiptypes.SignatureTypeUnderlyingType(sigg.GetSignatureType()),
			Content:       sigg.GetContent(),
		}),
	}
	sendingAtte.Src = transactionReply.Src
	sendingAtte.Dst = transactionReply.Dst

	if ret, err := nsbcli.NewNSBClient(hs).InsuranceClaim(
		signer,
		sendingAtte.SessionId,
		sendingAtte.Atte.Tid,
		TxState.Instantiating,
	); err != nil {
		vc.logger.Error("VesClient.read.AttestationSendingRequest.InsuranceClaim", "ip", hs, "error", err)
		return
	} else {
		vc.logger.Info(
			"insurance claiming",
			"info", ret.Info,
			"data", string(ret.Data),
			"log", ret.Log,
			"tags", ret.Tags,
		)
	}

	err = vc.postRawMessage(wsrpc.CodeAttestationReceiveRequest, transactionReply.Dst, sendingAtte)
	if err != nil {
		vc.logger.Error("VesClient.read.AttestationSendingRequest.postRawMessage", "error", err)
		return
	}
	vc.logger.Info("post next attestation request successfully")
}
