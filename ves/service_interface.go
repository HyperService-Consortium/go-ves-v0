package ves

import (
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/HyperService-Consortium/go-ves/grpc/uiprpc"
	uipbase "github.com/HyperService-Consortium/go-ves/grpc/uiprpc-base"
	log "github.com/HyperService-Consortium/go-ves/lib/log"
	"github.com/HyperService-Consortium/go-ves/ves/service"
	"golang.org/x/net/context"
)

func (server *Server) UserRegister(
	ctx context.Context,
	in *uiprpc.UserRegisterRequest,
) (*uiprpc.UserRegisterReply, error) {
	log.Infof("registering: %v\n", hex.EncodeToString(in.GetAccount().GetAddress()))
	return service.UserRegisterService{
		VESDB:               server.db,
		Context:             ctx,
		UserRegisterRequest: in,
	}.Serve()
}

func (server *Server) SessionStart(
	ctx context.Context,
	in *uiprpc.SessionStartRequest,
) (*uiprpc.SessionStartReply, error) {
	log.Infof("session start requesting\n")
	return (&service.SessionStartService{
		Signer:              server.signer,
		NsbClient:           server.nsbClient,
		CVes:                server.cves,
		VESDB:               server.db,
		Context:             ctx,
		SessionStartRequest: in,
	}).Serve()
}

func (server *Server) SessionAckForInit(
	ctx context.Context,
	in *uiprpc.SessionAckForInitRequest,
) (*uiprpc.SessionAckForInitReply, error) {
	log.Infof("user acknowledged: %v\n", hex.EncodeToString(in.GetUser().GetAddress()))
	return service.SessionAckForInitService{
		CVes:                     server.cves,
		Signer:                   server.signer,
		VESDB:                    server.db,
		Context:                  ctx,
		SessionAckForInitRequest: in,
	}.Serve()
}

func (server *Server) SessionRequireTransact(
	ctx context.Context,
	in *uiprpc.SessionRequireTransactRequest,
) (*uiprpc.SessionRequireTransactReply, error) {
	log.Infof("user request transact\n")
	return service.SessionRequireTransactService{
		VESDB:                         server.db,
		Context:                       ctx,
		SessionRequireTransactRequest: in,
	}.Serve()
}
func (server *Server) SessionRequireRawTransact(
	ctx context.Context,
	in *uiprpc.SessionRequireRawTransactRequest,
) (*uiprpc.SessionRequireRawTransactReply, error) {
	log.Infof("user request transact (computed)\n")
	return service.SessionRequireRawTransactService{
		Resp:                             server.resp,
		VESDB:                            server.db,
		Context:                          ctx,
		SessionRequireRawTransactRequest: in,
	}.Serve()
}

func (server *Server) AttestationReceive(
	ctx context.Context,
	in *uiprpc.AttestationReceiveRequest,
) (*uiprpc.AttestationReceiveReply, error) {
	log.Infof("attestation recevied: %v, %v\n", in.GetAtte().GetTid(), in.GetAtte().GetAid())
	return (&service.AttestationReceiveService{
		Signer:                    server.signer,
		CVes:                      server.cves,
		VESDB:                     server.db,
		NsbClient:                 server.nsbClient,
		Context:                   ctx,
		AttestationReceiveRequest: in,
	}).Serve()
}

func (server *Server) MerkleProofReceive(
	ctx context.Context,
	in *uiprpc.MerkleProofReceiveRequest,
) (*uiprpc.MerkleProofReceiveReply, error) {
	log.Infof("merkleproof recevied: %v, %v\n", in.GetMerkleproof().GetKey(), in.GetMerkleproof().GetValue())
	return (&service.MerkleProofReceiveService{
		VESDB:                     server.db,
		NsbClient:                 server.nsbClient,
		Context:                   ctx,
		MerkleProofReceiveRequest: in,
	}).Serve()
}

func (server *Server) ShrotenMerkleProofReceive(
	ctx context.Context,
	in *uiprpc.ShortenMerkleProofReceiveRequest,
) (*uiprpc.ShortenMerkleProofReceiveReply, error) {
	log.Infof("merkleproof recevied: %v, %v\n", in.GetMerkleproof().GetKey(), in.GetMerkleproof().GetValue())
	return (&service.ShrotenMerkleProofReceiveService{
		VESDB:                            server.db,
		NsbClient:                        server.nsbClient,
		Context:                          ctx,
		ShortenMerkleProofReceiveRequest: in,
	}).Serve()
}

func (server *Server) InformMerkleProof(
	ctx context.Context,
	in *uiprpc.MerkleProofReceiveRequest,
) (*uiprpc.MerkleProofReceiveReply, error) {
	log.Infof("merkleproof recevied: %v, %v\n", in.GetMerkleproof().GetKey(), in.GetMerkleproof().GetValue())
	return (&service.InformMerkleProofService{
		VESDB:                     server.db,
		NsbClient:                 server.nsbClient,
		Context:                   ctx,
		MerkleProofReceiveRequest: in,
	}).Serve()
}

func (server *Server) InformShortenMerkleProof(
	ctx context.Context,
	in *uiprpc.ShortenMerkleProofReceiveRequest,
) (*uiprpc.ShortenMerkleProofReceiveReply, error) {
	log.Infof("merkleproof recevied: %v, %v\n", in.GetMerkleproof().GetKey(), in.GetMerkleproof().GetValue())
	return (&service.InformShortenMerkleProofService{
		VESDB:                            server.db,
		NsbClient:                        server.nsbClient,
		Context:                          ctx,
		ShortenMerkleProofReceiveRequest: in,
	}).Serve()
}

func (server *Server) InformAttestation(
	ctx context.Context,
	in *uiprpc.AttestationReceiveRequest,
) (*uiprpc.AttestationReceiveReply, error) {
	log.Infof("informing attestation: %v, %v\n", in.GetAtte().GetTid(), in.GetAtte().GetAid())
	return (&service.InformAttestationService{
		Signer:                    server.signer,
		CVes:                      server.cves,
		VESDB:                     server.db,
		NsbClient:                 server.nsbClient,
		Context:                   ctx,
		AttestationReceiveRequest: in,
	}).Serve()
}

func (server *Server) requestSendSessionInfo(sessionID []byte, requestingAccount []*uipbase.Account) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	r, err := server.cves.InternalRequestComing(
		ctx,
		&uiprpc.InternalRequestComingRequest{
			SessionId: sessionID,
			Host:      []byte("todo"),
			Accounts: func() []*uipbase.Account {
				return nil
			}(),
		})
	if err != nil {
		return fmt.Errorf("could not request: %v", err)
	}
	if !r.GetOk() {
		return errors.New("internal failed")
	}
	return nil
}
