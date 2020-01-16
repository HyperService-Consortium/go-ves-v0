package raw_transaction

import (
	base_raw_transaction "github.com/HyperService-Consortium/go-uip/base-raw-transaction"
	"github.com/HyperService-Consortium/go-uip/uiptypes"
)

func FromRaw(b []byte) uiptypes.RawTransaction {
	return base_raw_transaction.Transaction(b)
}
