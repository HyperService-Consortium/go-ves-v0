package bni

import (
	"github.com/HyperService-Consortium/go-ves/types"
	"time"

	"github.com/HyperService-Consortium/go-uip/uiptypes"
)


type BN struct {
	dns types.ChainDNSInterface
	signer uiptypes.Signer
}


type options struct {
	timeout time.Duration
}

func parseOptions(rOption []interface{}) options {
	var parsedOptions options
	for i := range rOption {
		switch option := rOption[i].(type) {
		case time.Duration:
			parsedOptions.timeout = option
		case uiptypes.RouteOptionTimeout:
			parsedOptions.timeout = time.Duration(option)
		}
	}
	return parsedOptions
}

func NewBN(dns types.ChainDNSInterface) *BN {
	return &BN{dns: dns}
}