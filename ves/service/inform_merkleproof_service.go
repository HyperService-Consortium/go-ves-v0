package service

import (
	nsbcli "github.com/HyperService-Consortium/go-ves/lib/net/nsb-client"
	"golang.org/x/net/context"

	uiprpc "github.com/HyperService-Consortium/go-ves/grpc/uiprpc"
	types "github.com/HyperService-Consortium/go-ves/types"
	// bni "github.com/HyperService-Consortium/go-ves/types/bn-interface"
)

type InformMerkleProofService struct {
	NsbClient *nsbcli.NSBClient
	types.VESDB
	context.Context
	*uiprpc.MerkleProofReceiveRequest
}

func (s *InformMerkleProofService) Serve() (*uiprpc.MerkleProofReceiveReply, error) {
	s.ActivateSession(s.GetSessionId())
	ses, err := s.FindSessionInfo(s.GetSessionId())
	if err == nil {
		defer func() {
			s.UpdateSessionInfo(ses)
			s.InactivateSession(s.GetSessionId())
		}()

		var merkle = s.GetMerkleproof()

		// todo: verify merkle proof

		err = s.SetKV(
			ses.GetGUID(),
			merkle.GetKey(),
			merkle.GetValue(),
		)

		if err != nil {
			return nil, err
		}

		return &uiprpc.MerkleProofReceiveReply{
			Ok: true,
		}, nil

	} else {
		s.InactivateSession(s.GetSessionId())
		return nil, err
	}
}
