package bni

import (
	"encoding/hex"
	"fmt"
	"github.com/HyperService-Consortium/go-uip/const/trans_type"
	"github.com/HyperService-Consortium/go-uip/uiptypes"
	"github.com/HyperService-Consortium/go-ves/config"
	"github.com/HyperService-Consortium/go-ves/types"
	"github.com/Myriad-Dreamin/minimum-lib/sugar"
	"testing"
)

func TestBN_Translate(t *testing.T) {

//	{
//		"op-intents": [
//	{
//		"name": "Op1",
//		"op_type": "ContractInvocation",
//		"invoker": {
//			"domain": 2,
//			"user_name": "a1"
//		},
//		"contract_addr": "00e1eaa022cc40d4808bfe62b8997540c914d81e",
//		"func": "updateStake",
//		"parameters": [
//		{
//			"type": "uint256",
//			"value": {
//				"constant": "1000"
//			}
//		}
//	],
//		"amount": "0",
//		"unit": "wei"
//	}
//],
//	"dependencies": []
//	}
	type fields struct {
		dns    types.ChainDNSInterface
		signer uiptypes.Signer
	}
	type args struct {
		intent  *uiptypes.TransactionIntent
		storage uiptypes.Storage
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		//want    uiptypes.RawTransaction
		wantErr bool
	}{
		{"test_easy", fields{
			dns:    config.ChainDNS,
			signer: nil,
		}, args{
			intent: &uiptypes.TransactionIntent{
				TransType: trans_type.ContractInvoke,
				Src:       sugar.HandlerError(hex.DecodeString("93334ae4b2d42ebba8cc7c797bfeb02bfb3349d6")).([]byte),
				Dst:       sugar.HandlerError(hex.DecodeString("263fef3fe76fd4075ac16271d5115d01206d3674")).([]byte),
				Meta:      []byte(`{"contract_code":null,"func":"updateStake","parameters":[{"Type":"uint256","Value":{"constant":"1000"}}],"meta":null}`),
				Amt:       "03e8",
				ChainID:   2,
			},
			storage: nil,
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bn := &BN{
				dns:    tt.fields.dns,
				signer: tt.fields.signer,
			}
			got, err := bn.Translate(tt.args.intent, tt.args.storage)
			if (err != nil) != tt.wantErr {
				t.Errorf("Translate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}
			fmt.Println(string(sugar.HandlerError(got.Serialize()).([]byte)))
			//if !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("Translate() got = %v, want %v", got, tt.want)
			//}
		})
	}
}

func Test_decoratePrefix(t *testing.T) {
	fmt.Println(decoratePrefix("041a"))
}