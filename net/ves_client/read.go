package vesclient

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"time"

	uiprpc "github.com/Myriad-Dreamin/go-ves/grpc/uiprpc"
	uipbase "github.com/Myriad-Dreamin/go-ves/grpc/uiprpc-base"
	wsrpc "github.com/Myriad-Dreamin/go-ves/grpc/wsrpc"
	"github.com/gogo/protobuf/proto"
	"google.golang.org/grpc"

	log "github.com/Myriad-Dreamin/go-ves/log"
	helper "github.com/Myriad-Dreamin/go-ves/net/help-func"
	nsbcli "github.com/Myriad-Dreamin/go-ves/net/nsb_client"
	service "github.com/Myriad-Dreamin/go-ves/net/ves_client/service"
)

const todo = uint32(1)

func (vc *VesClient) read() {
	for {
		_, message, err := vc.conn.ReadMessage()
		if err != nil {
			fmt.Println("read:", err)
			return
		}

		var buf = bytes.NewBuffer(message)
		var messageID uint16
		binary.Read(buf, binary.BigEndian, &messageID)
		switch messageID {
		case wsrpc.CodeMessageReply:

			var s = vc.getShortReplyMessage()
			err = proto.Unmarshal(buf.Bytes(), s)
			if err != nil {
				// ignoring
				// todo: add hidden log
				continue
			}

			if bytes.Equal(s.To, vc.getName()) {
				fmt.Printf("%v is saying: %v\n", string(s.From), s.Contents)
			}
		case wsrpc.CodeClientHelloReply:
			var s = vc.getClientHelloReply()
			err = proto.Unmarshal(buf.Bytes(), s)
			if err != nil {
				// ignoring
				// todo: add hidden log
				continue
			}

			vc.grpcip, err = helper.DecodeIP(s.GetGrpcHost())
			if err != nil {
				// ignoring
				// todo: add hidden log
			}
			vc.nsbip, err = helper.DecodeIP(s.GetNsbHost())
			if err != nil {
				// ignoring
				// todo: add hidden log
			}
		case wsrpc.CodeRequestComingRequest:
			var s = vc.getrequestComingRequest()
			err = proto.Unmarshal(buf.Bytes(), s)
			if err != nil {
				// ignoring
				// todo: add hidden log
				log.Println(err)
				continue
			}
			fmt.Println("comming", hex.EncodeToString(s.GetSessionId()), hex.EncodeToString(s.GetAccount().GetAddress()))
			signer, err := vc.getSigner()
			if err != nil {
				log.Println(err)
				continue
			}

			// hs, err := helper.DecodeIP(s.GetNsbHost())
			// if err != nil {
			// 	log.Println(err)
			// 	continue
			// }
			hs := "http://47.251.2.73:26657"
			fmt.Println("send ack to nsb", hs)
			nc := nsbcli.NewNSBClient(hs)
			if ret, err := nc.UserAck(vc.signer, s.GetSessionId(), s.GetAccount().GetAddress(), []byte("123")); err != nil {

				log.Println(err)
				continue
			} else {
				fmt.Printf("user ack {\n\tinfo: %v,\n\tdata: %v,\n\tlog: %v, \n\ttags: %v\n}\n", ret.Info, string(ret.Data), ret.Log, ret.Tags)
			}
			if err = vc.sendAck(s.GetAccount(), s.GetSessionId(), s.GetGrpcHost(), signer.Sign(s.GetSessionId())); err != nil {

				log.Println(err)
				continue
			}
		case wsrpc.CodeAttestationSendingRequest:
			var s = vc.getrequestComingRequest()
			err = proto.Unmarshal(buf.Bytes(), s)
			if err != nil {
				// ignoring
				// todo: add hidden log
				log.Println(err)
				continue
			}
			fmt.Println(
				"comming attestation",
				hex.EncodeToString(s.GetSessionId()),
				hex.EncodeToString(s.GetAccount().GetAddress()))
			raw, tid, src, dst, err := vc.getRawTransaction(s.GetSessionId(), s.GetGrpcHost())
			if err != nil {
				// ignoring
				// todo: add hidden log
				log.Println(err)
				continue
			}
			fmt.Println(string(raw))
			var t = vc.getReceiveAttestationReceiveRequest()
			t.SessionId = s.GetSessionId()
			t.GrpcHost = s.GetGrpcHost()

			// todo
			type Type = uint64

			const (
				Unknown Type = 0 + iota
				Initing
				Inited
				Instantiating
				Instantiated
				Open
				Opened
				Closed
			)

			t.Atte = &uipbase.Attestation{
				Tid:     tid,
				Aid:     Instantiating,
				Content: raw,
				Signatures: append(make([]*uipbase.Signature, 0, 1), &uipbase.Signature{
					SignatureType: 123456,
					// todo use src.signer to sign
					Content: vc.signer.Sign(raw),
				}),
			}

			// hs, err := helper.DecodeIP(s.GetNsbHost())
			// if err != nil {
			// 	log.Println(err)
			// 	continue
			// }
			hs := "http://47.251.2.73:26657"
			fmt.Println("send ack to nsb", hs)
			nc := nsbcli.NewNSBClient(hs)
			if ret, err := nc.InsuranceClaim(vc.signer, s.GetSessionId(), tid, Instantiating); err != nil {

				log.Println(err)
				continue
			} else {
				fmt.Printf("insurance claim instantiating {\n\tinfo: %v,\n\tdata: %v,\n\tlog: %v, \n\ttags: %v\n}\n", ret.Info, string(ret.Data), ret.Log, ret.Tags)
			}

			t.Src = src
			t.Dst = dst
			err = vc.postRawMessage(wsrpc.CodeAttestationReceiveRequest, dst, t)
			if err != nil {
				log.Println(err)
				continue
			}

		case wsrpc.CodeUserRegisterReply:
			buf.Reset()
			vc.cb <- buf

		case wsrpc.CodeAttestationReceiveRequest:
			fmt.Println("receive attestation...")
			var s = vc.getReceiveAttestationReceiveRequest()
			err = proto.Unmarshal(buf.Bytes(), s)
			if err != nil {
				// ignoring
				// todo: add hidden log
				log.Println(err)
				continue
			}

			signer, err := vc.getSigner()
			if err != nil {
				log.Println(err)
				continue
			}

			if _, err = (&service.AttestationReceiveService{
				Signer:                    signer,
				NSBClient:                 vc.nsbClient,
				AttestationReceiveRequest: s,
			}).Serve(); err != nil {
				log.Println(err)
				continue
				// else if err = vc.postMessage(wsrpc.CodeAttestationReceiveReply, msg); err != nil {
				// 	log.Println(err)
				// 	continue
				// 	// Closed = 7
				// }
			} else if s.GetAtte().GetAid() != 7 {

				log.Infof("must send attestation")

				func() {

					var t = vc.getSendAttestationReceiveRequest()
					t.SessionId = s.GetSessionId()
					t.GrpcHost = s.GetGrpcHost()
					atte := s.GetAtte()
					sigs := atte.GetSignatures()
					t.Atte = &uipbase.Attestation{
						Tid:     atte.GetTid(),
						Aid:     atte.GetAid() + 1,
						Content: sigs[len(sigs)-1].GetContent(),
						Signatures: append(sigs, &uipbase.Signature{
							SignatureType: todo,
							Content:       signer.Sign(sigs[len(sigs)-1].GetContent()),
						}),
					}
					t.Src = s.GetDst()
					err = vc.postRawMessage(wsrpc.CodeAttestationReceiveRequest, s.GetSrc(), t)
					if err != nil {
						log.Println(err)
						return
					}
					return
				}()

				func() {

					mAddress, err := helper.DecodeIP(s.GetGrpcHost())
					if err != nil {
						log.Println(err)
						return
					}

					conn, err := grpc.Dial(mAddress, grpc.WithInsecure())
					if err != nil {
						log.Fatalf("did not connect: %v", err)
					}
					defer conn.Close()
					c := uiprpc.NewVESClient(conn)

					ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
					defer cancel()
					atte := s.GetAtte()
					sigs := atte.GetSignatures()
					r, err := c.AttestationReceive(
						ctx,
						&uiprpc.AttestationReceiveRequest{
							SessionId: s.GetSessionId(),
							Atte: &uipbase.Attestation{
								Tid:     atte.GetTid(),
								Aid:     atte.GetAid() + 1,
								Content: sigs[len(sigs)-1].GetContent(),
								Signatures: append(sigs, &uipbase.Signature{
									SignatureType: todo,
									Content:       signer.Sign(sigs[len(sigs)-1].GetContent()),
								}),
							},
						})
					if err != nil {
						log.Printf("could not greet: %v\n", err)
						return
					}
					if !r.GetOk() {
						log.Println("atte to grpc failed")
					}
				}()
			}
		case wsrpc.CodeCloseSessionRequest:
			vc.cb <- buf
			log.Println("session closed")
		// case wsrpc.Code
		default:
			fmt.Println("aborting message", string(message))
			// abort
		}

	}
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

// func (vc *VesClient) attestationReceive() {
//
// }
