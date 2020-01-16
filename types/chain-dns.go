package types

type ChainDNS interface {
	GetChainInfo(Index, chain_id) (ChainInfo, error)
}

type ChainDNSInterface interface {
	GetChainInfo(chain_id) (ChainInfo, error)
}
