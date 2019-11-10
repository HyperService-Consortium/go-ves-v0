package vesclient

import (
	"bytes"
	"encoding/hex"
	"github.com/HyperService-Consortium/go-uip/signaturer"

	uiptypes "github.com/HyperService-Consortium/go-uip/uiptypes"
	uip_base "github.com/HyperService-Consortium/go-ves/grpc/uiprpc-base"
)

func (vc *VesClient) getRespSigner(acc *uip_base.Account) (uiptypes.Signer, error) {
	if vc.signer != nil {
		return vc.signer, nil
	}
	cid := acc.GetChainId()
	switch cid {
	case 1, 2:
		sadd := hex.EncodeToString(acc.GetAddress())
		for _, acc := range vc.accs.Alias {
			if acc.ChainID == cid && acc.Address == sadd {
				return &acc, nil
			}
		}
	case 3, 4:
		for _, key := range vc.keys.Alias {
			if key.ChainID != cid {
				continue
			}

			signer, err := signaturer.NewTendermintNSBSigner(key.PrivateKey)
			if err != nil {
				return nil, err
			}
			if bytes.Equal(signer.GetPublicKey(), acc.GetAddress()) {
				return signer, nil
			}
		}
	}

	return nil, errNotFound
}
