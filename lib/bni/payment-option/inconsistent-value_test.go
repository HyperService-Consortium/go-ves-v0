package payment_option

import (
	"fmt"
	base_variable "github.com/HyperService-Consortium/go-uip/base-variable"
	"github.com/HyperService-Consortium/go-uip/const/value_type"
	"github.com/HyperService-Consortium/go-uip/uiptypes"
	"github.com/tidwall/gjson"
	"math/big"
	"testing"
)

type storage0777 struct {
}

func (s storage0777) GetTransactionProof(chainID uiptypes.ChainID, blockID uiptypes.BlockID, color []byte) (uiptypes.MerkleProof, error) {
	panic("implement me")
}

func (s storage0777) GetStorageAt(chainID uiptypes.ChainID, typeID uiptypes.TypeID, contractAddress uiptypes.ContractAddress, pos []byte, description []byte) (uiptypes.Variable, error) {
	fmt.Println(chainID, typeID, contractAddress, pos, description)
	return base_variable.Variable{
		Type:  value_type.Uint256,
		Value: big.NewInt(0x777),
	}, nil
}

func (s storage0777) SetStorageOf(chainID uiptypes.ChainID, typeID uiptypes.TypeID, contractAddress uiptypes.ContractAddress, pos []byte, description []byte, variable uiptypes.Variable) error {
	panic("implement me")
}

func TestParseInconsistentValueOption(t *testing.T) {

	type args struct {
		meta         gjson.Result
		storage      uiptypes.Storage
		defaultValue string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"test_easy", args{
			meta:         gjson.Parse(`{"value-inconsistent":{"type":"uint256","value":{"constant":"123"}}`),
			storage:      nil,
			defaultValue: "1123",
		}, "0123", false},
		{"test_storage", args{
			meta:         gjson.Parse(`{"value-inconsistent":{"type":"uint256","value":{"domain": 2, "contract":"00e1eaa022cc40d4808bfe62b8997540c914d81e","field":"strikePrice","pos":"01"}}}`),
			storage:      storage0777{},
			defaultValue: "1123",
		}, "0777", false},
		{"test_nothing", args{
			meta:         gjson.Parse(``),
			storage:      storage0777{},
			defaultValue: "1123",
		}, "1123", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseInconsistentValueOption(tt.args.meta, tt.args.storage, tt.args.defaultValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseInconsistentValueOption() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseInconsistentValueOption() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNeedInconsistentValueOption(t *testing.T) {
	type args struct {
		meta gjson.Result
	}
	tests := []struct {
		name string
		args args
		//want    *Need
		//want1   bool
		wantErr bool
	}{
		{"test_easy", args{
			meta: gjson.Parse(`{"value-inconsistent":{"type":"uint256","value":{"constant":"123"}}`),
		}, false},
		{"test_storage", args{
			meta: gjson.Parse(`{"value-inconsistent":{"type":"uint256","value":{"domain": 2, "contract":"00e1eaa022cc40d4808bfe62b8997540c914d81e","field":"strikePrice","pos":"01"}}}`),
		}, false},
		{"test_nothing", args{
			meta: gjson.Parse(``),
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := NeedInconsistentValueOption(tt.args.meta)
			if (err != nil) != tt.wantErr {
				t.Errorf("NeedInconsistentValueOption() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}
			fmt.Println(got, got1)
			//if !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("NeedInconsistentValueOption() got = %v, want %v", got, tt.want)
			//}
			//if got1 != tt.want1 {
			//	t.Errorf("NeedInconsistentValueOption() got1 = %v, want %v", got1, tt.want1)
			//}
		})
	}
}
