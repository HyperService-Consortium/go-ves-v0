package payment_option

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/HyperService-Consortium/go-uip/const/value_type"
	"github.com/HyperService-Consortium/go-uip/uiptypes"
	"github.com/tidwall/gjson"
)

func ParseInconsistentValueOption(meta gjson.Result, storage uiptypes.Storage, defaultValue string) (string, error) {
	optionIc := meta.Get("value-inconsistent")
	if optionIc.Exists() {
		t := optionIc.Get("type").String()
		optionValue := optionIc.Get("value")
		constant := optionValue.Get("constant")
		if !constant.Exists() {
			domain := optionValue.Get("domain")
			if !domain.Exists() {
				return defaultValue, errors.New("domain not found")
			}
			field := optionValue.Get("field")
			if !field.Exists() {
				return defaultValue, errors.New("field not found")
			}
			pos := optionValue.Get("pos")
			if !pos.Exists() {
				return defaultValue, errors.New("pos not found")
			}
			contract := optionValue.Get("contract")
			if !contract.Exists() {
				return defaultValue, errors.New("contract not found")
			}
			var contractAddress, err = hex.DecodeString(contract.String())
			if err != nil {
				return defaultValue, err
			}
			bPos, err := toPos(pos)
			if err != nil {
				return defaultValue, err
			}
			v, err := storage.GetStorageAt(domain.Uint(), value_type.FromString(t), contractAddress, bPos, []byte{})
			if err != nil {
				return defaultValue, err
			}
			defaultValue = fmt.Sprintf("%x", v.GetValue())
		} else {
			defaultValue = constant.String()
		}
		if (len(defaultValue) & 1) != 0 {
			defaultValue = "0" + defaultValue
		}
	}
	return defaultValue, nil
}

type Need struct {
	ChainID         uiptypes.ChainID
	TypeID          uiptypes.TypeID
	ContractAddress uiptypes.ContractAddress
	Pos             []byte
	Description     []byte
}

func NeedInconsistentValueOption(meta gjson.Result) (*Need, bool, error) {
	optionIc := meta.Get("value-inconsistent")
	if optionIc.Exists() {
		t := optionIc.Get("type").String()
		optionValue := optionIc.Get("value")
		constant := optionValue.Get("constant")
		if !constant.Exists() {
			domain := optionValue.Get("domain")
			if !domain.Exists() {
				return nil, false, errors.New("domain not found")
			}
			field := optionValue.Get("field")
			if !field.Exists() {
				return nil, false, errors.New("field not found")
			}
			pos := optionValue.Get("pos")
			if !pos.Exists() {
				return nil, false, errors.New("pos not found")
			}
			contract := optionValue.Get("contract")
			if !contract.Exists() {
				return nil, false, errors.New("contract not found")
			}
			var contractAddress, err = hex.DecodeString(contract.String())
			if err != nil {
				return nil, false, err
			}
			bPos, err := toPos(pos)
			if err != nil {
				return nil, false, err
			}
			return &Need{
				domain.Uint(),
				value_type.FromString(t),
				contractAddress,
				bPos,
				[]byte{},
			}, true, nil
		}
	}
	return nil, false, nil
}

func toPos(result gjson.Result) ([]byte, error) {
	i := result.String()
	y, err := hex.DecodeString(i)
	if err != nil {
		return nil, err
	}
	if len(y) > 32 {
		return nil, errors.New("overflowed")
	}
	return append(make([]byte, 32-len(y)), y...), nil
}
