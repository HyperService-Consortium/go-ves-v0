package bni

import (
	"encoding/hex"
	"fmt"
	"github.com/HyperService-Consortium/go-ethabi"
	"github.com/HyperService-Consortium/go-uip/const/value_type"
	"github.com/HyperService-Consortium/go-uip/uiptypes"
	"github.com/HyperService-Consortium/go-ves/config"
	"github.com/HyperService-Consortium/go-ves/types"
	"github.com/Myriad-Dreamin/minimum-lib/sugar"
	"testing"
)

func TestBN_GetStorageAt(t *testing.T) {
	type fields struct {
		dns    types.ChainDNSInterface
		signer uiptypes.Signer
	}
	type args struct {
		chainID         uiptypes.ChainID
		typeID          uiptypes.TypeID
		contractAddress uiptypes.ContractAddress
		pos             []byte
		description     []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		//want    uiptypes.Variable
		wantErr bool
	}{
		{"test_easy", fields{
			dns:    config.ChainDNS,
			signer: nil,
		}, args{
			chainID:         2,
			typeID:          value_type.Uint256,
			contractAddress: sugar.HandlerError(hex.DecodeString("263fef3fe76fd4075ac16271d5115d01206d3674")).([]byte),
			pos:             []byte{1},
			description:     nil,
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bn := &BN{
				dns:    tt.fields.dns,
				signer: tt.fields.signer,
			}
			got, err := bn.GetStorageAt(tt.args.chainID, tt.args.typeID, tt.args.contractAddress, tt.args.pos, tt.args.description)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetStorageAt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Println(got)
			//if !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("GetStorageAt() got = %v, want %v", got, tt.want)
			//}
		})
	}
}

func TestZ(t *testing.T) {
	var b = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}
	x, err := ethabi.NewDecoder().Decodes([]string{"uint256"}, b)
	fmt.Println(x, err)
}
