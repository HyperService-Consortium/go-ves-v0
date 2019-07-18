package service

import (
	"errors"
	"fmt"

	"golang.org/x/net/context"

	uiptypes "github.com/Myriad-Dreamin/go-uip/types"
	uiprpc "github.com/Myriad-Dreamin/go-ves/grpc/uip-rpc"
	types "github.com/Myriad-Dreamin/go-ves/types"
)

type AttestationReceiveService struct {
	types.VESDB
	context.Context
	*uiprpc.AttestationReceiveRequest
}

type AtteAdapdator struct {
	*uiprpc.Attestation
}

func (atte *AtteAdapdator) GetSignatures() []uiptypes.Signature {
	var ss = atte.Attestation.GetSignatures()
	ret := make([]uiptypes.Signature, len(ss))
	for _, s := range ss {
		ret = append(ret, uiptypes.Signature(s))
	}
	return ret
}

func (s *AttestationReceiveService) Serve() (*uiprpc.AttestationReceiveReply, error) {
	ses, err := s.FindSessionInfo(s.GetSessionId())
	if err == nil {
		var success bool
		var help_info string
		success, help_info, err = ses.ProcessAttestation(&AtteAdapdator{s.GetAtte()})
		if err != nil {
			// todo, log
			return nil, fmt.Errorf("internal error: %v", err)
		} else if !success {
			return nil, errors.New(help_info)
		} else {
			return &uiprpc.AttestationReceiveReply{
				Ok: true,
			}, nil
		}

	} else {
		return nil, err
	}
}
