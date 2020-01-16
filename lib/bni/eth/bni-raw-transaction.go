package bni

import (
	"encoding/hex"
	"encoding/json"
	"github.com/HyperService-Consortium/go-uip/uiptypes"
	ethclient "github.com/HyperService-Consortium/go-ves/lib/net/eth-client"
)

type RawTransaction struct {
	B           []byte `json:"b" form:"b"`
	Responsible []byte `json:"r" form:"r"`
	IsSigned    bool   `json:"i" form:"i"`
	// todo... remove targetHost
	TargetHost string `json:"t" form:"t"`
}

func NewRawTransaction(b []byte, responsible []byte, isSigned bool, targetHost string) *RawTransaction {
	return &RawTransaction{B: b, Responsible: responsible, IsSigned: isSigned, TargetHost: targetHost}
}

func (t RawTransaction) Serialize() ([]byte, error) {
	return json.Marshal(&t)
}

func (t RawTransaction) Bytes() ([]byte, error) {
	return t.B, nil
}

func (t RawTransaction) Signed() bool {
	return t.IsSigned
}

type PasswordSigner interface {
	uiptypes.Signer
	GetEthPassword() string
}

type passwordSigner struct {
	pb []byte
	ps string
}

func (p passwordSigner) GetPublicKey() uiptypes.PublicKey {
	return p.pb
}

func (p passwordSigner) Sign(uiptypes.SignatureContent) uiptypes.Signature {
	panic("implement me")
}

func (p passwordSigner) GetEthPassword() string {
	return p.ps
}

// todo change raw transaction signature = sign(signer, context)
func (t RawTransaction) Sign(signer uiptypes.Signer) (uiptypes.RawTransaction, error) {
	if s, ok := signer.(PasswordSigner); ok {
		unlock, err := ethclient.NewEthClient(t.TargetHost).PersonalUnlockAccout(hex.EncodeToString(s.GetPublicKey()), s.GetEthPassword(), 100)
		t.IsSigned = unlock
		return t, err
	}
	return t, nil
}
