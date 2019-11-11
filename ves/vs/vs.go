package vs

import (
	"github.com/HyperService-Consortium/go-uip/signaturer"
	"github.com/HyperService-Consortium/go-ves/grpc/uiprpc"
	uiprpc_base "github.com/HyperService-Consortium/go-ves/grpc/uiprpc-base"
	nsbcli "github.com/HyperService-Consortium/go-ves/lib/net/nsb-client"
	"github.com/HyperService-Consortium/go-ves/types"
	"github.com/Myriad-Dreamin/minimum-lib/logger"
)

type VServer struct {
	Logger logger.Logger

	DB        types.VESDB
	Resp      *uiprpc_base.Account
	Signer    *signaturer.TendermintNSBSigner
	CVes      uiprpc.CenteredVESClient
	NsbClient *nsbcli.NSBClient
	Host      []byte
	NsbHost   []byte
}
