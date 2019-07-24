package nsbi

import (
	uiptypes "github.com/Myriad-Dreamin/go-uip/types"
	nsbcli "github.com/Myriad-Dreamin/go-ves/net/nsb_client"
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
	_, err := nsb.NSBClient.InsuranceClaim(nsb.signer, isc_address, atte.GetTid(), uint64(len(atte.GetSignatures())+2))
	return err
}

func (nsb *NSBClientImpl) SettleContract(isc_address []byte) error {
	_, err := nsb.NSBClient.SettleContract(nsb.signer, isc_address)
	return err
}
