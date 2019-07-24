package ethtl

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"strings"

	TransType "github.com/Myriad-Dreamin/go-uip/const/trans_type"
	opintent "github.com/Myriad-Dreamin/go-uip/op-intent"
)

type Translator struct {
}

func decoratePrefix(hexs string) string {
	if strings.HasPrefix(hexs, "0x") {
		return hexs
	} else {
		return "0x" + hexs
	}
}

func (cl *Translator) Translate(
	intent *opintent.TransactionIntent,
	// kvs KeyValues,
) ([]byte, error) {
	switch intent.TransType {
	case TransType.Payment:
		return json.Marshal(map[string]interface{}{
			"from":  decoratePrefix(hex.EncodeToString(intent.Src)),
			"to":    decoratePrefix(hex.EncodeToString(intent.Dst)),
			"value": decoratePrefix(intent.Amt),
		})
	case TransType.ContractInvoke:
		return nil, errors.New("todo")
	default:
		return nil, errors.New("cant translate")
	}
}
