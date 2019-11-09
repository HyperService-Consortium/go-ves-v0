package vesclient

import (
	"encoding/hex"
	"fmt"
	"github.com/HyperService-Consortium/go-uip/signaturer"
)

func (vc *VesClient) ListKeys() {
	fmt.Println("privatekeys -> publickeys:")
	for alias, key := range vc.keys.Alias {
		fmt.Println(
			"alias:", alias,
			"public key:", hex.EncodeToString(signaturer.NewTendermintNSBSigner(key.PrivateKey).GetPublicKey()),
			"chain id:", key.ChainID,
		)
	}
	fmt.Println("ethAccounts:")
	for alias, acc := range vc.accs.Alias {
		fmt.Println(
			"alias:", alias,
			"public address:", acc.Address,
			"chain id:", acc.ChainID,
		)
	}
}
