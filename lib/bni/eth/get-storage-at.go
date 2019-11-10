package bni

import (
	"encoding/hex"
	"github.com/HyperService-Consortium/go-ethabi"
	"github.com/HyperService-Consortium/go-uip/const/value_type"
	"github.com/HyperService-Consortium/go-uip/uiptypes"
	ethclient "github.com/HyperService-Consortium/go-ves/lib/net/eth-client"
	"net/url"
)

type variable struct {
	Type uiptypes.TypeID
	Value interface{}
}

func (v variable) GetType() uiptypes.TypeID {
	return v.Type
}

func (v variable) GetValue() interface{} {
	return v.Value
}

func (bn *BN) GetStorageAt(chainID uiptypes.ChainID, typeID uiptypes.TypeID, contractAddress uiptypes.ContractAddress, pos []byte, description []byte) (uiptypes.Variable, error) {
	// todo
	ci, err := bn.dns.GetChainInfo(chainID)
	if err != nil {
		return nil, err
	}

	switch typeID {
	case value_type.Bool:
		s, err := ethclient.NewEthClient((&url.URL{Scheme: "http", Host: ci.GetChainHost(), Path: "/"}).String()).GetStorageAt(contractAddress, pos, "latest")
		if err != nil {
			return nil, err
		}

		b, err := hex.DecodeString(s[2:])
		if err != nil {
			return nil, err
		}
		bs, err := ethabi.NewDecoder().Decodes([]string{"bool"}, b)
		return variable{
			Type:  typeID,
			Value: bs[0],
		}, nil
	case value_type.Uint256:
		s, err := ethclient.NewEthClient(ci.GetChainHost()).GetStorageAt(contractAddress, pos, "latest")
		if err != nil {
			return nil, err
		}

		b, err := hex.DecodeString(s[2:])
		if err != nil {
			return nil, err
		}
		bs, err := ethabi.NewDecoder().Decodes([]string{"uint256"}, b)
		return variable{
			Type:  typeID,
			Value: bs[0],
		}, nil
	}

	return nil, nil
}
