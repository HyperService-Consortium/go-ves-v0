package vesclient

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	log "github.com/HyperService-Consortium/go-ves/lib/log"

	"github.com/gogo/protobuf/proto"
	"google.golang.org/grpc"

	TxState "github.com/HyperService-Consortium/go-uip/const/transaction_state_type"

	uiprpc "github.com/HyperService-Consortium/go-ves/grpc/uiprpc"
	uipbase "github.com/HyperService-Consortium/go-ves/grpc/uiprpc-base"
	wsrpc "github.com/HyperService-Consortium/go-ves/grpc/wsrpc"

	helper "github.com/HyperService-Consortium/go-ves/lib/net/help-func"
	nsbcli "github.com/HyperService-Consortium/go-ves/lib/net/nsb-client"
)

func (vc *VesClient) read() {
	for {
		_, message, err := vc.conn.ReadMessage()
		if err != nil {
			log.Errorln("VesClient.read.read:", err)
			return
		}

		tag := md5.Sum(message)
		log.Infoln("message tag", hex.EncodeToString(tag[:]))

		var buf = bytes.NewBuffer(message)
		var messageID uint16

		// todo: BigEndian
		binary.Read(buf, binary.BigEndian, &messageID)
		switch messageID {
		case wsrpc.CodeMessageReply:
			var messageReply = vc.getShortReplyMessage()

			err = proto.Unmarshal(buf.Bytes(), messageReply)
			if err != nil {
				log.Errorln("VesClient.read.MessageReply.proto:", err)
				continue
			}

			if bytes.Equal(messageReply.To, vc.getName()) {
				fmt.Printf("%v is saying: %v\n", string(messageReply.From), messageReply.Contents)
			}
		case wsrpc.CodeClientHelloReply:
			var clientHelloReply = vc.getClientHelloReply()

			err = proto.Unmarshal(buf.Bytes(), clientHelloReply)
			if err != nil {
				log.Errorln("VesClient.read.ClientHelloReply.proto:", err)
				continue
			}

			vc.grpcip, err = helper.DecodeIP(clientHelloReply.GetGrpcHost())
			if err != nil {
				log.Errorln("VesClient.read.ClientHelloReply.decodeGRPCHost:", err)
			} else {
				log.Infoln("adding default grpcip ", vc.grpcip)
			}

			vc.nsbip, err = helper.DecodeIP(clientHelloReply.GetNsbHost())
			if err != nil {
				log.Errorln("VesClient.read.ClientHelloReply.decodeNSBHost:", err)
			} else {
				log.Infoln("adding default nsbip ", vc.nsbip)
			}

		case wsrpc.CodeRequestComingRequest:
			var requestComingRequest = vc.getrequestComingRequest()

			err = proto.Unmarshal(buf.Bytes(), requestComingRequest)
			if err != nil {
				log.Errorln("VesClient.read.RequestComingRequest.proto:", err)
				continue
			}

			fmt.Println(
				"new session request comming:",
				"session id:", hex.EncodeToString(requestComingRequest.GetSessionId()),
				"resposible address:", hex.EncodeToString(requestComingRequest.GetAccount().GetAddress()),
			)

			signer, err := vc.getNSBSigner()
			if err != nil {
				log.Errorln("VesClient.read.RequestComingRequest.getNSBSigner:", err)
				continue
			}

			hs, err := helper.DecodeIP(requestComingRequest.GetNsbHost())
			if err != nil {
				log.Errorln("VesClient.read.RequestComingRequest.DecodeIP:", err)
				continue
			}

			fmt.Println("send ack to nsb", hs)

			// todo: new nsbclient
			if ret, err := nsbcli.NewNSBClient(hs).UserAck(
				signer,
				requestComingRequest.GetSessionId(),
				requestComingRequest.GetAccount().GetAddress(),
				// todo: signature
				[]byte("123"),
			); err != nil {
				log.Errorln("VesClient.read.RequestComingRequest.UserAck:", err)
				continue
			} else {
				fmt.Printf(
					"user ack {\n\tinfo: %v,\n\tdata: %v,\n\tlog: %v, \n\ttags: %v\n}\n",
					ret.Info, string(ret.Data), ret.Log, ret.Tags,
				)
			}

			if err = vc.sendAck(
				requestComingRequest.GetAccount(),
				requestComingRequest.GetSessionId(),
				requestComingRequest.GetGrpcHost(),
				signer.Sign(requestComingRequest.GetSessionId()).Bytes(),
			); err != nil {
				log.Errorln("VesClient.read.RequestComingRequest.sendAck:", err)
				continue
			}

		case wsrpc.CodeAttestationSendingRequest:
			// attestation sending request has the same format with request
			// coming request
			var attestationSendingRequest = vc.getrequestComingRequest()
			err = proto.Unmarshal(buf.Bytes(), attestationSendingRequest)
			if err != nil {
				log.Errorln("VesClient.read.AttestationSendingRequest.proto:", err)
				continue
			}

			fmt.Println(
				"new transaction's attestation must be created",
				hex.EncodeToString(attestationSendingRequest.GetSessionId()),
				hex.EncodeToString(attestationSendingRequest.GetAccount().GetAddress()),
			)

			raw, tid, src, dst, err := vc.GetRawTransaction(
				attestationSendingRequest.GetSessionId(),
				attestationSendingRequest.GetGrpcHost(),
			)
			if err != nil {
				log.Errorln("VesClient.read.AttestationSendingRequest.getRawTransaction:", err)
				continue
			}

			fmt.Printf(
				"the instance of the %vth transaction intent is: %v\n", tid,
				string(raw),
			)

			signer, err := vc.getNSBSigner()
			if err != nil {
				log.Errorln("VesClient.read.AttestationSendingRequest.getNSBSigner:", err)
				continue
			}

			hs, err := helper.DecodeIP(attestationSendingRequest.GetNsbHost())
			if err != nil {
				log.Errorln("VesClient.read.AttestationSendingRequest.DecodeIP:", err)
				continue
			}

			// packet attestation
			var sendingAtte = vc.getReceiveAttestationReceiveRequest()
			sendingAtte.SessionId = attestationSendingRequest.GetSessionId()
			sendingAtte.GrpcHost = attestationSendingRequest.GetGrpcHost()

			sigg := signer.Sign(raw)
			sendingAtte.Atte = &uipbase.Attestation{
				Tid:     tid,
				Aid:     TxState.Instantiating,
				Content: raw,
				Signatures: append(make([]*uipbase.Signature, 0, 1), &uipbase.Signature{
					// todo use src.signer to sign
					SignatureType: sigg.GetSignatureType(),
					Content:       sigg.GetContent(),
				}),
			}
			sendingAtte.Src = src
			sendingAtte.Dst = dst

			fmt.Println("send ack to nsb", hs)
			if ret, err := nsbcli.NewNSBClient(hs).InsuranceClaim(
				signer,
				sendingAtte.SessionId,
				sendingAtte.Atte.Tid,
				TxState.Instantiating,
			); err != nil {
				log.Errorln("VesClient.read.AttestationSendingRequest.InsuranceClaim:", err)
				continue
			} else {
				fmt.Printf(
					formatInsuranceClaim,
					"instantiating", ret.Info, string(ret.Data), ret.Log, ret.Tags,
				)
			}

			err = vc.postRawMessage(wsrpc.CodeAttestationReceiveRequest, dst, sendingAtte)
			if err != nil {
				log.Errorln("VesClient.read.AttestationSendingRequest.postRawMessage:", err)
				continue
			}

		case wsrpc.CodeUserRegisterReply:
			// todo: ignoring
			vc.cb <- buf

		case wsrpc.CodeAttestationReceiveRequest:
			var s = vc.getReceiveAttestationReceiveRequest()

			err = proto.Unmarshal(buf.Bytes(), s)
			if err != nil {
				log.Errorln("VesClient.read.AttestationReceiveRequest.proto:", err)
				continue
			}

			atte := s.GetAtte()
			aid := atte.GetAid()

			switch aid {
			case TxState.Unknown:
				log.Infoln("transaction is of the status unknown")
			case TxState.Initing:
				log.Infoln("transaction is of the status initing")
			case TxState.Inited:
				log.Infoln("transaction is of the status inited")
			case TxState.Closed:
				// skip closed atte (last)
				log.Infoln("skip the last attestation of this transaction")
			default:
				log.Infoln("must send attestation with status:", TxState.Description(aid+1))

				signer, err := vc.getNSBSigner()
				if err != nil {
					log.Errorln("VesClient.read.AttestationReceiveRequest.getNSBSigner:", err)
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

					log.Infoln("the resp is", hex.EncodeToString(acc.GetAddress()), acc.GetChainId())

					router := vc.getRouter(acc.ChainId)
					if router == nil {
						log.Errorln("VesClient.read.AttestationReceiveRequest.getRouter:", errors.New("get router failed"))
						continue
					}

					if router.MustWithSigner() {
						respSigner, err := vc.getRespSigner(s.GetDst())
						if err != nil {
							log.Errorln("VesClient.read.AttestationReceiveRequest.getRespSigner:", err)
							continue
						}

						router = router.RouteWithSigner(respSigner)
					}

					receipt, err := router.RouteRawTransaction(acc.ChainId, atte.GetContent())
					if err != nil {
						log.Errorln("VesClient.read.AttestationReceiveRequest.router.RouteRaw:", err)
						continue
					}
					fmt.Println("receipt:", hex.EncodeToString(receipt), string(receipt))

					bid, additional, err := router.WaitForTransact(acc.ChainId, receipt, vc.waitOpt)
					if err != nil {
						log.Errorln("VesClient.read.AttestationReceiveRequest.router.WaitForTransact:", err)
						continue
					}
					fmt.Println("route result:", bid)

					blockStorage := vc.getBlockStorage(acc.ChainId)
					if blockStorage == nil {
						log.Errorln("VesClient.read.AttestationReceiveRequest.getBlockStorage:", errors.New("get BlockStorage failed"))
						continue
					}

					proof, err := blockStorage.GetTransactionProof(acc.GetChainId(), bid, additional)
					if err != nil {
						log.Errorln("VesClient.read.AttestationReceiveRequest.blockStorage.GetTransactionProof:", err)
						continue
					}

					cb, err := vc.nsbClient.AddMerkleProof(signer, nil, proof.GetType(), proof.GetRootHash(), proof.GetProof(), proof.GetKey(), proof.GetValue())
					if err != nil {
						log.Errorln("VesClient.read.AttestationReceiveRequest.nsbClient.AddMerkleProof:", err)
						continue
					}
					fmt.Println("adding merkle proof", cb)

					// todo: add const TransactionsRoot
					cb, err = vc.nsbClient.AddBlockCheck(signer, nil, acc.ChainId, bid, proof.GetRootHash(), 1)
					if err != nil {
						log.Errorln("VesClient.read.AttestationReceiveRequest.nsbClient.AddBlockCheck:", err)
						continue
					}
					fmt.Println("adding block check", cb)
				}

				// sendingAtte.GetAtte()
				ret, err := vc.nsbClient.InsuranceClaim(
					signer,
					s.GetSessionId(),
					atte.GetTid(), aid+1,
				)
				//sessionID, tid, Instantiated)
				if err != nil {
					log.Errorln("VesClient.read.AttestationReceiveRequest.InsuranceClaim:", err)
					continue
				}

				fmt.Printf(
					"insurance claiming %v, %v {\n\tinfo: %v,\n\tdata: %v,\n\tlog: %v, \n\ttags: %v\n}\n",
					atte.GetTid(), TxState.Description(aid+1), ret.Info, string(ret.Data), ret.Log, ret.Tags,
				)

				err = vc.postRawMessage(wsrpc.CodeAttestationReceiveRequest, s.GetSrc(), sendingAtte)
				if err != nil {
					log.Errorln("VesClient.read.AttestationReceiveRequest.postRawMessage:", err)
					continue
				}

				grpcHost, err := helper.DecodeIP(s.GetGrpcHost())
				if err != nil {
					log.Println(err)
					return
				}

				vc.informAttestation(grpcHost, sendingAtte)
			}

		case wsrpc.CodeCloseSessionRequest:
			vc.cb <- buf
			log.Infoln("session closed")
			// case wsrpc.Code
		default:
			// abort
			log.Warnln("aborting message that has id:", messageID)
		} // switch end
	} // for end
}

func (vc *VesClient) sendAck(acc *uipbase.Account, sessionID, address, signature []byte) error {
	// Set up a connection to the server.
	sss, err := helper.DecodeIP(address)
	if err != nil {
		return fmt.Errorf("did not resolve: %v", err)
	}
	conn, err := grpc.Dial(sss, grpc.WithInsecure())
	if err != nil {
		return fmt.Errorf("did not connect: %v", err)
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
		return fmt.Errorf("could not greet: %v", err)
	}
	fmt.Printf("Session ack: %v\n", r.GetOk())
	return nil
}

func (vc *VesClient) informAttestation(grpcHost string, sendingAtte *wsrpc.AttestationReceiveRequest) {
	conn, err := grpc.Dial(grpcHost, grpc.WithInsecure())
	if err != nil {
		log.Errorf("VesClient.informAttestation(%v, %v).grpc.Dial: %v\n", sendingAtte.GetAtte().Tid, sendingAtte.GetAtte().Aid, err)
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
		log.Errorf("VesClient.informAttestation(%v, %v).InformAttestation: %v\n", sendingAtte.GetAtte().Tid, sendingAtte.GetAtte().Aid, err)
		return
	}

	if !r.GetOk() {
		log.Errorf("VesClient.informAttestation(%v, %v) failed\n", sendingAtte.GetAtte().Tid, sendingAtte.GetAtte().Aid)
	}
	return
}
