package vesclient

import (
	"encoding/hex"
	signaturetype "github.com/HyperService-Consortium/go-uip/const/signature_type"
	"github.com/HyperService-Consortium/go-uip/signaturer"
	uiptypes "github.com/HyperService-Consortium/go-uip/uiptypes"
)

// ECCKey is the private key object in memory
type ECCKey struct {
	PrivateKey []byte `json:"private_key"`
	ChainID    uint64 `json:"chain_id"`
}

// ECCKeyAlias is the private key object in json
type ECCKeyAlias struct {
	PrivateKey string `json:"private_key"`
	ChainID    uint64 `json:"chain_id"`
	Alias      string `json:"alias"`
}

// EthAccount is the account object in memory
type EthAccount struct {
	Address    string `json:"address"`
	ChainID    uint64 `json:"chain_id"`
	PassPhrase string `json:"pass_phrase"`
}

// EthAccountAlias is the account object in json
type EthAccountAlias struct {
	EthAccount
	Alias string `json:"alias"`
}

// ECCKeys is the object saved in files
type ECCKeys struct {
	Keys  []*ECCKey
	Alias map[string]ECCKey
}

// EthAccounts is the object saved in files
type EthAccounts struct {
	Accs  []*EthAccount
	Alias map[string]EthAccount
}


func (signer *EthAccount) GetPublicKey() []byte {
	b, _ := hex.DecodeString(signer.Address)
	return b
}

func (signer *EthAccount) Sign(b []byte) uiptypes.Signature {
	// todo: sign b
	return signaturer.FromRaw(b, uiptypes.SignatureUnderlyingType(signaturetype.Secp256k1))
}