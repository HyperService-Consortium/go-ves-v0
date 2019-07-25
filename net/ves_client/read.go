package vesclient

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"time"

	log "github.com/Myriad-Dreamin/go-ves/log"

	"context"

	uiprpc "github.com/Myriad-Dreamin/go-ves/grpc/uiprpc"
	uipbase "github.com/Myriad-Dreamin/go-ves/grpc/uiprpc-base"
	wsrpc "github.com/Myriad-Dreamin/go-ves/grpc/wsrpc"
	"github.com/gogo/protobuf/proto"
	"google.golang.org/grpc"

	TxState "github.com/Myriad-Dreamin/go-uip/const/transaction_state_type"
	helper "github.com/Myriad-Dreamin/go-ves/net/help-func"
	nsbcli "github.com/Myriad-Dreamin/go-ves/net/nsb_client"
	service "github.com/Myriad-Dreamin/go-ves/net/ves_client/service"
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

			signer, err := vc.getSigner()
			if err != nil {
				log.Errorln("VesClient.read.RequestComingRequest.getSigner:", err)
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
				signer.Sign(requestComingRequest.GetSessionId()),
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

			raw, tid, src, dst, err := vc.getRawTransaction(
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

			signer, err := vc.getSigner()
			if err != nil {
				log.Errorln("VesClient.read.AttestationSendingRequest.getSigner:", err)
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
			sendingAtte.Atte = &uipbase.Attestation{
				Tid:     tid,
				Aid:     TxState.Instantiating,
				Content: raw,
				Signatures: append(make([]*uipbase.Signature, 0, 1), &uipbase.Signature{
					// todo use src.signer to sign
					SignatureType: todo,
					Content:       vc.signer.Sign(raw),
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

			signer, err := vc.getSigner()
			if err != nil {
				log.Errorln("VesClient.read.AttestationReceiveRequest.getSigner:", err)
				continue
			}

			if _, err = (&service.AttestationReceiveService{
				Signer:                    signer,
				NSBClient:                 vc.nsbClient,
				AttestationReceiveRequest: s,
			}).Serve(); err != nil {
				log.Errorln("VesClient.read.AttestationReceiveRequest.AttestationReceiveService:", err)
				continue
				// else if err = vc.postMessage(wsrpc.CodeAttestationReceiveReply, msg); err != nil {
				// 	log.Println(err)
				// 	continue
				// }
			} else {
				atte := s.GetAtte()

				// skip closed atte (last)
				if atte.GetAid() == TxState.Closed {
					log.Infoln("skip the last attestation of this transaction")
					continue
				}

				log.Infoln("must send attestation with status:", TxState.Description(s.GetAtte().GetAid()))

				sigs := atte.GetSignatures()
				toSig := sigs[len(sigs)-1].GetContent()

				var sendingAtte = vc.getSendAttestationReceiveRequest()
				sendingAtte.SessionId = s.GetSessionId()
				sendingAtte.GrpcHost = s.GetGrpcHost()

				// todo: iter the atte (copy or refer it? )
				sendingAtte.Atte = &uipbase.Attestation{
					Tid: atte.GetTid(),
					// todo: get nx -> more readable
					Aid:     atte.GetAid() + 1,
					Content: atte.GetContent(),
					Signatures: append(sigs, &uipbase.Signature{
						// todo signature
						SignatureType: todo,
						Content:       signer.Sign(toSig),
					}),
				}
				sendingAtte.Src = s.GetDst()
				sendingAtte.Dst = s.GetSrc()

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
