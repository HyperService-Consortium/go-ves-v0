package chain_dns

import (
	"errors"
	"github.com/HyperService-Consortium/go-uip/uiptypes"
	"github.com/HyperService-Consortium/go-ves/types"
)

type NSBHostOption string

type ChainInfo struct {
	ChainType uiptypes.ChainType
	Host string
}

func (c ChainInfo) GetChainType() uiptypes.ChainType {
	return c.ChainType
}

func (c ChainInfo) GetChainHost() string {
	return c.Host
}

type HostMap map[uiptypes.ChainID]ChainInfo

type ServerOptions struct {
	HostMap HostMap
}

func defaultServerOptions() ServerOptions {
	return ServerOptions{
		HostMap: make(HostMap),
	}
}

func parseOptions(rOptions []interface{}) ServerOptions {
	var options = defaultServerOptions()
	for i := range rOptions {
		switch option := rOptions[i].(type) {
		case HostMap:
			options.HostMap = option
		}
	}
	return options
}


type Database struct {
	HostMap HostMap
}

var ErrNotFound = errors.New("not found")

func (d Database) GetChainInfo(_ types.Index, chainID uint64) (types.ChainInfo, error) {
	if res, ok := d.HostMap[chainID]; ok {
		return res, nil
	}

	return nil, ErrNotFound
}

func NewDatabase(rOptions ...interface{}) *Database {
	options := parseOptions(rOptions)
	return &Database{
		HostMap: options.HostMap,
	}
}







