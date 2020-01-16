package config

import (
	ChainType "github.com/HyperService-Consortium/go-uip/const/chain_type"
	"github.com/HyperService-Consortium/go-uip/uiptypes"
	"github.com/HyperService-Consortium/go-ves/types"
	chain_dns "github.com/HyperService-Consortium/go-ves/types/chain-dns"
)

var HostMap = chain_dns.HostMap{
	1: chain_dns.ChainInfo{
		Host:      "127.0.0.1:8545",
		ChainType: ChainType.Ethereum,
	},
	2: chain_dns.ChainInfo{
		Host:      "127.0.0.1:8545",
		ChainType: ChainType.Ethereum,
	},
	3: chain_dns.ChainInfo{
		Host:      "127.0.0.1:26657",
		ChainType: ChainType.TendermintNSB,
	},
	4: chain_dns.ChainInfo{
		Host:      "127.0.0.1:26657",
		ChainType: ChainType.TendermintNSB,
	},
	5: chain_dns.ChainInfo{
		Host:      "127.0.0.1:26657",
		ChainType: ChainType.TendermintNSB,
	},
	6: chain_dns.ChainInfo{
		Host:      "39.100.145.91:8545",
		ChainType: ChainType.Ethereum,
	},
}

func GetHostMap() chain_dns.HostMap {
	return HostMap
}

var ChainDNS types.ChainDNSInterface

type IdleChainDNS struct {
	dns types.ChainDNS
}

func (i IdleChainDNS) GetChainInfo(chainId uiptypes.ChainID) (types.ChainInfo, error) {
	return i.dns.GetChainInfo(nil, chainId)
}

func init() {
	ChainDNS = IdleChainDNS{chain_dns.NewDatabase(HostMap)}
}
