package nsbi

import (
	uiptypes "github.com/HyperService-Consortium/go-uip/types"
	nsbcli "github.com/HyperService-Consortium/go-ves/lib/net/nsb-client"
)

type NSBClientImpl struct {
	*nsbcli.NSBClient
	signer uiptypes.Signer
}

func NSBInterfaceImpl(host string, signer uiptypes.Signer) *NSBClientImpl {
	return &NSBClientImpl{nsbcli.NewNSBClient(host), signer}
}

func (nsb *NSBClientImpl) SaveAttestation(isc_address []byte, atte uiptypes.Attestation) error {
	// todo
	return nil
}
func (nsb *NSBClientImpl) InsuranceClaim(isc_address []byte, atte uiptypes.Attestation) error {
	_, err := nsb.NSBClient.InsuranceClaim(nsb.signer, isc_address, atte.GetTid(), atte.GetAid())
	return err
}

func (nsb *NSBClientImpl) SettleContract(isc_address []byte) error {
	_, err := nsb.NSBClient.SettleContract(nsb.signer, isc_address)
	return err
}
