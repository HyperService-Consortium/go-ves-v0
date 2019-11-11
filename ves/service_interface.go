package ves

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/HyperService-Consortium/go-ves/ves/vs"
	"time"

	"github.com/HyperService-Consortium/go-ves/grpc/uiprpc"
	uipbase "github.com/HyperService-Consortium/go-ves/grpc/uiprpc-base"
	"github.com/HyperService-Consortium/go-ves/ves/service"
	"golang.org/x/net/context"
)

func (server *Server) UserRegister(
	ctx context.Context,
	in *uiprpc.UserRegisterRequest,
) (*uiprpc.UserRegisterReply, error) {
	return service.NewUserRegisterService((*vs.VServer)(server), ctx, in).Serve()
}

func (server *Server) SessionStart(
	ctx context.Context,
	in *uiprpc.SessionStartRequest,
) (*uiprpc.SessionStartReply, error) {
	fmt.Printf("ves server: session start: intents count: %v, dependecies count: %v\n", len(in.Opintents.Contents), len(in.Opintents.Dependencies))
	return service.NewSessionStartService((*vs.VServer)(server), ctx, in).Serve()
}

func (server *Server) SessionAckForInit(
	ctx context.Context,
	in *uiprpc.SessionAckForInitRequest,
) (*uiprpc.SessionAckForInitReply, error) {
	fmt.Printf("ves server: session acknowledging: address: %v\n", hex.EncodeToString(in.GetUser().GetAddress()))
	return service.NewSessionAckForInitService((*vs.VServer)(server), ctx, in).Serve()
}

func (server *Server) SessionRequireTransact(
	ctx context.Context,
	in *uiprpc.SessionRequireTransactRequest,
) (*uiprpc.SessionRequireTransactReply, error) {
	return service.SessionRequireTransactService{
		VESDB:                         server.DB,
		Context:                       ctx,
		SessionRequireTransactRequest: in,
	}.Serve()
}

func (server *Server) SessionRequireRawTransact(
	ctx context.Context,
	in *uiprpc.SessionRequireRawTransactRequest,
) (*uiprpc.SessionRequireRawTransactReply, error) {
	return service.NewSessionRequireRawTransactService((*vs.VServer)(server), ctx, in).Serve()
}

func (server *Server) AttestationReceive(
	ctx context.Context,
	in *uiprpc.AttestationReceiveRequest,
) (*uiprpc.AttestationReceiveReply, error) {
	return service.NewAttestationReceiveService((*vs.VServer)(server), ctx, in).Serve()
}

func (server *Server) MerkleProofReceive(
	ctx context.Context,
	in *uiprpc.MerkleProofReceiveRequest,
) (*uiprpc.MerkleProofReceiveReply, error) {
	return service.NewMerkleProofReceiveService((*vs.VServer)(server), ctx, in).Serve()
}

func (server *Server) ShrotenMerkleProofReceive(
	ctx context.Context,
	in *uiprpc.ShortenMerkleProofReceiveRequest,
) (*uiprpc.ShortenMerkleProofReceiveReply, error) {
	return service.NewShortenMerkleProofReceiveService((*vs.VServer)(server), ctx, in).Serve()
}

func (server *Server) InformMerkleProof(
	ctx context.Context,
	in *uiprpc.MerkleProofReceiveRequest,
) (*uiprpc.MerkleProofReceiveReply, error) {
	return service.NewInformMerkleProofService((*vs.VServer)(server), ctx, in).Serve()
}

func (server *Server) InformShortenMerkleProof(
	ctx context.Context,
	in *uiprpc.ShortenMerkleProofReceiveRequest,
) (*uiprpc.ShortenMerkleProofReceiveReply, error) {
	return service.NewInformShortenMerkleProofService((*vs.VServer)(server), ctx, in).Serve()
}

func (server *Server) InformAttestation(
	ctx context.Context,
	in *uiprpc.AttestationReceiveRequest,
) (*uiprpc.AttestationReceiveReply, error) {
	fmt.Printf("ves server: attestation from client: tid:%v, aid: %v\n", in.Atte.Tid, in.Atte.Aid)
	return service.NewInformAttestationService((*vs.VServer)(server), ctx, in).Serve()
}

func (server *Server) requestSendSessionInfo(sessionID []byte, requestingAccount []*uipbase.Account) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	r, err := server.CVes.InternalRequestComing(
		ctx,
		&uiprpc.InternalRequestComingRequest{
			SessionId: sessionID,
			Host:      server.Host,
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
