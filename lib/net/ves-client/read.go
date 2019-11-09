package vesclient

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"time"

	"github.com/gogo/protobuf/proto"
	"google.golang.org/grpc"

	TxState "github.com/HyperService-Consortium/go-uip/const/transaction_state_type"

	"github.com/HyperService-Consortium/go-ves/grpc/uiprpc"
	uipbase "github.com/HyperService-Consortium/go-ves/grpc/uiprpc-base"
	"github.com/HyperService-Consortium/go-ves/grpc/wsrpc"

	helper "github.com/HyperService-Consortium/go-ves/lib/net/help-func"
)

func (vc *VesClient) read() {
	for {
		_, message, err := vc.conn.ReadMessage()
		if err != nil {
			vc.logger.Error("VesClient.read.read", "error", err)
			return
		}

		var buf = bytes.NewBuffer(message)
		var messageID uint16

		// todo: BigEndian
		err = binary.Read(buf, binary.BigEndian, &messageID)
		if err != nil {
			vc.logger.Error("decode message failed", "error", err)
			continue
		}

		tag := md5.Sum(message)
		vc.logger.Info("message from server", "tag", hex.EncodeToString(tag[:]), "type", wsrpc.MessageType(messageID))

		switch wsrpc.MessageType(messageID) {
		case wsrpc.CodeMessageReply:
			var messageReply = vc.getShortReplyMessage()

			err = proto.Unmarshal(buf.Bytes(), messageReply)
			if err != nil {
				vc.logger.Error("VesClient.read.MessageReply.proto", "error", err)
				continue
			}

			if bytes.Equal(messageReply.To, vc.getName()) {
				vc.logger.Info("%v is saying: %v\n", "source", string(messageReply.From), "content", messageReply.Contents)
			}
		case wsrpc.CodeClientHelloReply:
			var clientHelloReply = vc.getClientHelloReply()

			err = proto.Unmarshal(buf.Bytes(), clientHelloReply)
			if err != nil {
				vc.logger.Error("VesClient.read.ClientHelloReply.proto", "error", err)
				continue
			}

			vc.ProcessClientHelloReply(clientHelloReply)

		case wsrpc.CodeRequestComingRequest:

			var requestComingRequest = vc.getrequestComingRequest()
			err = proto.Unmarshal(buf.Bytes(), requestComingRequest)
			if err != nil {
				vc.logger.Error("VesClient.read.RequestComingRequest.proto", "error", err)
				continue
			}

			vc.ProcessRequestComingRequest(requestComingRequest)
		case wsrpc.CodeAttestationSendingRequest:
			// attestation sending request has the same format with request
			// coming request
			var attestationSendingRequest = vc.getrequestComingRequest()
			err = proto.Unmarshal(buf.Bytes(), attestationSendingRequest)
			if err != nil {
				vc.logger.Error("VesClient.read.AttestationSendingRequest.proto", "error", err)
				continue
			}

			vc.ProcessAttestationSendingRequest(attestationSendingRequest)

		case wsrpc.CodeUserRegisterReply:

			// todo: ignoring
			vc.cb <- buf
			vc.logger.Info("user request request successfully")

		case wsrpc.CodeAttestationReceiveRequest:
			var s = vc.getReceiveAttestationReceiveRequest()

			err = proto.Unmarshal(buf.Bytes(), s)
			if err != nil {
				vc.logger.Error("VesClient.read.AttestationReceiveRequest.proto", "error", err)
				continue
			}

			atte := s.GetAtte()
			aid := atte.GetAid()

			switch aid {
			case TxState.Unknown:
				vc.logger.Info("transaction is of the status unknown", "tid", atte.Tid)
			case TxState.Initing:
				vc.logger.Info("transaction is of the status initializing", "tid", atte.Tid)
			case TxState.Inited:
				vc.logger.Info("transaction is of the status inited", "tid", atte.Tid)
			case TxState.Closed:
				// skip closed atte (last)
				vc.logger.Info("skip the last attestation of this transaction", "tid", atte.Tid)
			default:
				vc.logger.Info("must send attestation with status", "tid", atte.Tid, "aid", TxState.Description(aid+1))

				signer, err := vc.getNSBSigner()
				if err != nil {
					vc.logger.Error("VesClient.read.AttestationReceiveRequest.getNSBSigner", "error", err)
					continue
				}

				sigs := atte.GetSignatures()
				toSig := sigs[len(sigs)-1].GetContent()

				var sendingAtte = vc.getSendAttestationReceiveRequest()
				sendingAtte.SessionId = s.GetSessionId()
				sendingAtte.GrpcHost = s.GetGrpcHost()

				// todo: iter the atte (copy or refer it? )
				sigg := signer.Sign(toSig)
				sendingAtte.Atte = &uipbase.Attestation{
					Tid: atte.GetTid(),
					// todo: get nx -> more readable
					Aid:     aid + 1,
					Content: atte.GetContent(),
					Signatures: append(sigs, &uipbase.Signature{
						// todo signature
						SignatureType: sigg.GetSignatureType(),
						Content:       sigg.GetContent(),
					}),
				}
				sendingAtte.Src = s.GetDst()
				sendingAtte.Dst = s.GetSrc()

				if aid == TxState.Instantiated {
					acc := s.GetDst()

					vc.logger.Info("the resp is", "address", hex.EncodeToString(acc.GetAddress()), "chain id", acc.GetChainId())

					router := vc.getRouter(acc.ChainId)
					if router == nil {
						vc.logger.Error("VesClient.read.AttestationReceiveRequest.getRouter:", "error", errors.New("get router failed"))
						continue
					}

					if router.MustWithSigner() {
						respSigner, err := vc.getRespSigner(s.GetDst())
						if err != nil {
							vc.logger.Error("VesClient.read.AttestationReceiveRequest.getRespSigner", "error", err)
							continue
						}

						router = router.RouteWithSigner(respSigner)
					}

					receipt, err := router.RouteRawTransaction(acc.ChainId, atte.GetContent())
					if err != nil {
						vc.logger.Error("VesClient.read.AttestationReceiveRequest.router.RouteRaw", "error", err)
						continue
					}
					vc.logger.Info("routing", "receipt", hex.EncodeToString(receipt))

					bid, additional, err := router.WaitForTransact(acc.ChainId, receipt, vc.waitOpt)
					if err != nil {
						vc.logger.Error("VesClient.read.AttestationReceiveRequest.router.WaitForTransact", "error", err)
						continue
					}
					vc.logger.Info("route result", "block id", bid)

					blockStorage := vc.getBlockStorage(acc.ChainId)
					if blockStorage == nil {
						vc.logger.Error("VesClient.read.AttestationReceiveRequest.getBlockStorage:", "error", errors.New("get BlockStorage failed"))
						continue
					}

					proof, err := blockStorage.GetTransactionProof(acc.GetChainId(), bid, additional)
					if err != nil {
						vc.logger.Error("VesClient.read.AttestationReceiveRequest.blockStorage.GetTransactionProof", "error", err)
						continue
					}

					cb, err := vc.nsbClient.AddMerkleProof(signer, nil, proof.GetType(), proof.GetRootHash(), proof.GetProof(), proof.GetKey(), proof.GetValue())
					if err != nil {
						vc.logger.Error("VesClient.read.AttestationReceiveRequest.nsbClient.AddMerkleProof", "error", err)
						continue
					}
					vc.logger.Info("adding merkle proof", "result", cb)

					// todo: add const TransactionsRoot
					cb, err = vc.nsbClient.AddBlockCheck(signer, nil, acc.ChainId, bid, proof.GetRootHash(), 1)
					if err != nil {
						vc.logger.Error("VesClient.read.AttestationReceiveRequest.nsbClient.AddBlockCheck", "error", err)
						continue
					}
					vc.logger.Info("adding block check", "result", cb)
				}

				// sendingAtte.GetAtte()
				ret, err := vc.nsbClient.InsuranceClaim(
					signer,
					s.GetSessionId(),
					atte.GetTid(), aid+1,
				)
				//sessionID, tid, Instantiated)
				if err != nil {
					vc.logger.Error("VesClient.read.AttestationReceiveRequest.InsuranceClaim", "error", err)
					continue
				}

				vc.logger.Info(
					"insurance claiming",
					"tid",atte.GetTid(),
					"aid", TxState.Description(aid+1),
					"info",ret.Info,
					"data", string(ret.Data),
					"log",ret.Log,
					"tags", ret.Tags,
				)

				err = vc.postRawMessage(wsrpc.CodeAttestationReceiveRequest, s.GetSrc(), sendingAtte)
				if err != nil {
					vc.logger.Error("VesClient.read.AttestationReceiveRequest.postRawMessage", "error", err)
					continue
				}

				grpcHost, err := helper.DecodeIP(s.GetGrpcHost())
				if err != nil {
					vc.logger.Error("VesClient.read.AttestationReceiveRequest.decodeIP", "error", err)
					return
				}

				vc.informAttestation(grpcHost, sendingAtte)
			}

		case wsrpc.CodeCloseSessionRequest:

			var closeSessionRequest = vc.getCloseSessionRequest()
			err = proto.Unmarshal(buf.Bytes(), closeSessionRequest)
			if err != nil {
				vc.logger.Error("VesClient.read.CodeCloseSessionRequest.proto", "error", err)
				continue
			}

			vc.cb <- buf
			vc.logger.Info("session closed")
			vc.emitClose(closeSessionRequest.SessionId)
			// case wsrpc.Code
		default:
			// abort
			vc.logger.Warn("aborting message", "id", messageID)
		} // switch end
	} // for end
}

func (vc *VesClient) sendAck(acc *uipbase.Account, sessionID, address, signature []byte) error {
	// Set up a connection to the server.
	sss, err := helper.DecodeIP(address)
	if err != nil {
		vc.logger.Error("did not resolve", "error", err)
		return err
	}
	conn, err := grpc.Dial(sss, grpc.WithInsecure())
	if err != nil {
		vc.logger.Error("did not connect", "error", err)
		return err
	}
	defer conn.Close()
	c := uiprpc.NewVESClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	r, err := c.SessionAckForInit(
		ctx,
		&uiprpc.SessionAckForInitRequest{
			SessionId: sessionID,
			User:      acc,
			UserSignature: &uipbase.Signature{
				SignatureType: 123456,
				Content:       signature,
			},
		})
	if err != nil {
		vc.logger.Error("could not send ack", "error", err)
		return err
	}
	vc.logger.Info("Session ack", "ok", r.GetOk(), "session id", sessionID)
	return nil
}

func (vc *VesClient) informAttestation(grpcHost string, sendingAtte *wsrpc.AttestationReceiveRequest) {
	conn, err := grpc.Dial(grpcHost, grpc.WithInsecure())
	if err != nil {
		vc.logger.Error("VesClient.informAttestation.grpc.Dial.Failed\n",
			"tid", sendingAtte.GetAtte().Tid, "aid", sendingAtte.GetAtte().Aid, "error", err)
		return
	}
	defer conn.Close()

	c := uiprpc.NewVESClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	r, err := c.InformAttestation(
		ctx,
		&uiprpc.AttestationReceiveRequest{
			SessionId: sendingAtte.SessionId,
			Atte:      sendingAtte.Atte,
		},
	)
	if err != nil {
		vc.logger.Error("VesClient.informAttestation.grpc.Send.Failed\n",
			"tid", sendingAtte.GetAtte().Tid, "aid", sendingAtte.GetAtte().Aid, "error", err)
		return
	}

	if !r.GetOk() {
		vc.logger.Error("VesClient.informAttestation.grpc.Result.Failed\n",
			"tid", sendingAtte.GetAtte().Tid, "aid", sendingAtte.GetAtte().Aid)
	}
	return
}
