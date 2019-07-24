package service

import (
	"errors"
	"fmt"

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
		ret, err := s.InsuranceClaim(s.Signer, sessionID, tid, Instantiated)
		if err != nil {
			return nil, err
		}
		fmt.Printf("insurance claiming instantiated {\n\tinfo: %v,\n\tdata: %v,\n\tlog: %v, \n\ttags: %v\n}\n", ret.Info, string(ret.Data), ret.Log, ret.Tags)

		return &wsrpc.AttestationReceiveReply{
			Ok: true,
		}, nil
	case Instantiated:
		ret, err := s.InsuranceClaim(s.Signer, sessionID, tid, Open)
		if err != nil {
			return nil, err
		}
		fmt.Printf("insurance claiming open {\n\tinfo: %v,\n\tdata: %v,\n\tlog: %v, \n\ttags: %v\n}\n", ret.Info, string(ret.Data), ret.Log, ret.Tags)

		// type = s.GetAtte().GetContent()
		// content = type.Content
		// s.BroadcastTxCommit(content)

		return &wsrpc.AttestationReceiveReply{
			Ok: true,
		}, nil
	case Open:
		ret, err := s.InsuranceClaim(s.Signer, sessionID, tid, Opened)
		if err != nil {
			return nil, err
		}
		fmt.Printf("insurance claiming opened {\n\tinfo: %v,\n\tdata: %v,\n\tlog: %v, \n\ttags: %v\n}\n", ret.Info, string(ret.Data), ret.Log, ret.Tags)

		return &wsrpc.AttestationReceiveReply{
			Ok: true,
		}, nil
	case Opened:
		ret, err := s.InsuranceClaim(s.Signer, sessionID, tid, Closed)
		if err != nil {
			return nil, err
		}
		fmt.Printf("insurance claiming closed {\n\tinfo: %v,\n\tdata: %v,\n\tlog: %v, \n\ttags: %v\n}\n", ret.Info, string(ret.Data), ret.Log, ret.Tags)

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
