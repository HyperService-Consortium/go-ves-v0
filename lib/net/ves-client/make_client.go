package vesclient

import (
	"github.com/Myriad-Dreamin/minimum-lib/logger"
)

// VanilleMakeClient for humans
func VanilleMakeClient(name, addr string, options ...interface{}) (*VesClient, error) {

	vcClient, err := NewVesClient(append(append(make([]interface{}, 2+len(options)),
		logger.NewZapDevelopmentSugarOption(), CVesHostOption(addr)), options...)...)
	if err != nil {
		return nil, err
	}

	if err = vcClient.Boot(); err != nil {
		return nil, err
	}
	return vcClient, nil
}
