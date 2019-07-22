package bni

import (
	"errors"
	"net/url"

	chain_type "github.com/Myriad-Dreamin/go-uip/const/chain_type"
	ethclient "github.com/Myriad-Dreamin/go-ves/net/eth_client"
)

type BN struct {
}

type ChainInfo struct {
	Host      string
	ChainType uint64
}

func (c *ChainInfo) GetHost() string {
	return c.Host
}

func (c *ChainInfo) GetChainType() uint64 {
	return c.ChainType
}

type ChainInfoInterface interface {
	GetHost() string
	GetChainType() uint64
}

func SearchChainId(domain uint64) (ChainInfoInterface, error) {
	switch domain {
	case 0:
		return nil, errors.New("nil domain is not allowed")
	case 1: // ethereum chain 1
		return &ChainInfo{
			Host:      "127.0.0.1",
			ChainType: chain_type.Ethereum,
		}, nil
	case 2: // ethereum chain 2
		return &ChainInfo{
			Host:      "127.0.0.1",
			ChainType: chain_type.Ethereum,
		}, nil
	case 3: // tendermint chain 1
		return &ChainInfo{
			Host:      "47.254.66.11",
			ChainType: chain_type.TendermintNSB,
		}, nil
	case 4: // ethereum chain 1
		return &ChainInfo{
			Host:      "47.251.2.73",
			ChainType: chain_type.TendermintNSB,
		}, nil
	default:
		return nil, errors.New("not found")
	}
}

func (bn *BN) RouteRaw(destination uint64, payload []byte) ([]byte, error) {
	ci, err := SearchChainId(destination)
	if err != nil {
		return nil, err
	}
	return ethclient.Do((&url.URL{Scheme: "http", Host: ci.GetHost(), Path: "/"}).String(), payload)
}

func (bn *BN) Route(destination uint64, on_chain_transaction []byte) ([]byte, error) {
	// todo
	return bn.RouteRaw(destination, on_chain_transaction)
}
