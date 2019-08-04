package service

import (
	"golang.org/x/net/context"

	uiprpc "github.com/Myriad-Dreamin/go-ves/grpc/uiprpc"
	types "github.com/Myriad-Dreamin/go-ves/types"
	// bni "github.com/Myriad-Dreamin/go-ves/types/bn-interface"
)

type ShrotenMerkleProofReceiveService struct {
	Host string
	types.VESDB
	context.Context
	*uiprpc.ShortenMerkleProofReceiveRequest
}

func (s *ShrotenMerkleProofReceiveService) Serve() (*uiprpc.ShortenMerkleProofReceiveReply, error) {
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

		return &uiprpc.ShortenMerkleProofReceiveReply{
			Ok: true,
		}, nil

	} else {
		s.InactivateSession(s.GetSessionId())
		return nil, err
	}
}
