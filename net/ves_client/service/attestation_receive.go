package service

import (
	"errors"

	uiptypes "github.com/Myriad-Dreamin/go-uip/types"
	wsrpc "github.com/Myriad-Dreamin/go-ves/grpc/wsrpc"
	nsbclient "github.com/Myriad-Dreamin/go-ves/net/nsb_client"
)

type AttestationReceiveService struct {
	uiptypes.Signer
	*nsbclient.NSBClient
	*wsrpc.AttestationReceiveRequest
}

func (s *AttestationReceiveService) Serve() (*wsrpc.AttestationReceiveReply, error) {

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
	atte, sessionID := s.GetAtte(), s.GetSessionId()
	tid, aid := atte.GetTid(), atte.GetAid()
	switch aid {
	case Unknown:
		return nil, errors.New("transaction is of the status unknown")
	case Initing:
		return nil, errors.New("transaction is of the status initing")
	case Inited:
		return nil, errors.New("transaction is of the status inited")
	case Instantiating:
		s.InsuranceClaim(s.Signer, sessionID, tid, Instantiated)
		return &wsrpc.AttestationReceiveReply{
			Ok: true,
		}, nil
	case Instantiated:
		s.InsuranceClaim(s.Signer, sessionID, tid, Open)

		// type = s.GetAtte().GetContent()
		// content = type.Content
		// s.BroadcastTxCommit(content)

		return &wsrpc.AttestationReceiveReply{
			Ok: true,
		}, nil
	case Open:
		s.InsuranceClaim(s.Signer, sessionID, tid, Opened)
		return &wsrpc.AttestationReceiveReply{
			Ok: true,
		}, nil
	case Opened:
		s.InsuranceClaim(s.Signer, sessionID, tid, Closed)
		return &wsrpc.AttestationReceiveReply{
			Ok: true,
		}, nil
	case Closed:
		return &wsrpc.AttestationReceiveReply{
			Ok: true,
		}, nil
	default:
		return nil, errors.New("unknown aid types")
	}
}
