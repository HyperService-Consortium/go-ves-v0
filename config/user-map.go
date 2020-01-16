package config

import (
	opintent "github.com/HyperService-Consortium/go-uip/op-intent"
	"github.com/HyperService-Consortium/go-uip/uiptypes"
)

type accountProvider struct {
}

func (a accountProvider) AccountBase() opintent.AccountBase {
	return a
}

func (accountProvider) Get(name string, chainId uint64) (uiptypes.Account, error) {
	return searchAccount(name, chainId)
}

func (accountProvider) GetRelay(domain uint64) (uiptypes.Account, error) {
	return getRelay(domain)
}

func (accountProvider) GetTransactionProofType(chainId uint64) (uiptypes.MerkleProofType, error) {
	return getTransactionProofType(chainId)
}

var UserMap = accountProvider{}
