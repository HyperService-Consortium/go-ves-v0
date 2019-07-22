package service

import (
	"errors"
	"fmt"

	"golang.org/x/net/context"

	uiptypes "github.com/Myriad-Dreamin/go-uip/types"
	uiprpc "github.com/Myriad-Dreamin/go-ves/grpc/uiprpc"
	uipbase "github.com/Myriad-Dreamin/go-ves/grpc/uiprpc-base"
	types "github.com/Myriad-Dreamin/go-ves/types"
	bni "github.com/Myriad-Dreamin/go-ves/types/bn-interface"
	nsbi "github.com/Myriad-Dreamin/go-ves/types/nsb-interface"
)

type AttestationReceiveService struct {
	Host string
	uiptypes.Signer
	types.VESDB
	context.Context
	*uiprpc.AttestationReceiveRequest
}

type AtteAdapdator struct {
	*uipbase.Attestation
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

	ses, err := s.FindSessionInfo(s.GetSessionId())
	if err == nil {
		ses.SetSigner(s.Signer)

		var success bool
		var helpInfo string
		success, helpInfo, err = ses.ProcessAttestation(
			nsbi.NSBInterfaceImpl(s.Host, s.Signer),
			&bni.BN{},
			&AtteAdapdator{s.GetAtte()},
		)
		if err != nil {
			// todo, log
			return nil, fmt.Errorf("internal error: %v", err)
		} else if !success {
			return nil, errors.New(helpInfo)
		} else {

			return &uiprpc.AttestationReceiveReply{
				Ok: true,
			}, nil
		}

	} else {
		return nil, err
	}
}
