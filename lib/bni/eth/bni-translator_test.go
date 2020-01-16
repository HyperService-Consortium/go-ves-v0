package bni

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/HyperService-Consortium/go-uip/const/trans_type"
	"github.com/HyperService-Consortium/go-uip/const/value_type"
	"github.com/HyperService-Consortium/go-uip/uiptypes"
	"github.com/HyperService-Consortium/go-ves/config"
	"github.com/HyperService-Consortium/go-ves/types"
	"github.com/Myriad-Dreamin/minimum-lib/sugar"
	"math/big"
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
		name   string
		fields fields
		args   args
		//want    uiptypes.RawTransaction
		wantErr bool
		assert gJSONAssertion
	}{
		//{"id":1,"jsonrpc":"2.0","method":"eth_sendTransaction","params":[{"data":"0x7c1f751f00000000000000000000000000000000000000000000000000000000000003e9","from":"0x93334ae4b2d42ebba8cc7c797bfeb02bfb3349d6","to":"0x263fef3fe76fd4075ac16271d5115d01206d3674"}]}
		{"test_easy", fields{
			dns:    config.ChainDNS,
			signer: nil,
		}, args{
			intent: &uiptypes.TransactionIntent{
				TransType: trans_type.ContractInvoke,
				Src:       sugar.HandlerError(hex.DecodeString("93334ae4b2d42ebba8cc7c797bfeb02bfb3349d6")).([]byte),
				Dst:       sugar.HandlerError(hex.DecodeString("263fef3fe76fd4075ac16271d5115d01206d3674")).([]byte),
				Meta: sugar.HandlerError(
					serializer.Meta.Contract.Marshal(
						&uiptypes.ContractInvokeMeta{
							FuncName: "updateStake",
							Params: []uiptypes.RawParams{
								{
									Type: "uint256",
									Value: marshal(h{
										"constant": 1001,
									}),
								},
							},
						})).([]byte),
				Amt:     "00",
				ChainID: 2,
			},
			storage: nil,
		}, false, gJSONWant(
				kv{"method", "eth_sendTransaction"},
			kv{"params.0.data", "0x7c1f751f00000000000000000000000000000000000000000000000000000000000003e9"},
			kv{"params.0.from", "0x93334ae4b2d42ebba8cc7c797bfeb02bfb3349d6"},
			kv{"params.0.to", "0x263fef3fe76fd4075ac16271d5115d01206d3674"},
			kv{"params.0.value", nil},
				)},
		{"test_with_storage_var", fields{
			dns:    config.ChainDNS,
			signer: nil,
		}, args{
			intent: &uiptypes.TransactionIntent{
				TransType: trans_type.ContractInvoke,
				Src:       sugar.HandlerError(hex.DecodeString("93334ae4b2d42ebba8cc7c797bfeb02bfb3349d6")).([]byte),
				Dst:       sugar.HandlerError(hex.DecodeString("263fef3fe76fd4075ac16271d5115d01206d3674")).([]byte),
				Meta: sugar.HandlerError(
					serializer.Meta.Contract.Marshal(
						&uiptypes.ContractInvokeMeta{
							FuncName: "updateStake",
							Params: []uiptypes.RawParams{
								{

									//if result.Get("contract").Exists() &&
									//	result.Get("pos").Exists() &&
									//	result.Get("field").Exists() {
									//	proposal = append(proposal, &MerkleProofProposal{
									//		DescriptionType:   merkleproof_proposal_type.DataProof,
									//		MerkleProofType:   merkleproofType,
									//		ValueType:         intDesc,
									//		SourceDescription: param.Value,
									//	})
									Type: "uint256",
									Value: marshal(h{
										"contract": "0000000000000000000000000000000000000000",
										"pos": "00",
										"field": "staking",
									}),
								},
							},
						})).([]byte),
				Amt:     "00",
				ChainID: 2,
			},
			storage: mockBNIStorage{data: []mockData{
				{
					//
					//"contract": make([]byte, 32),
					//"pos": []byte("00"),
					//"field": "staking",
					chainID:         2,
					typeID:          value_type.Uint256,
					contractAddress: make([]byte, 20),
					pos:             []byte("00"),
					description:     []byte("staking"),
					v:               mockValue{value_type.Uint256, big.NewInt(0x0233)},
				},
			}},
		}, false, gJSONWant(
			kv{"method", "eth_sendTransaction"},
			kv{"params.0.data", "0x7c1f751f0000000000000000000000000000000000000000000000000000000000000233"},
			kv{"params.0.from", "0x93334ae4b2d42ebba8cc7c797bfeb02bfb3349d6"},
			kv{"params.0.to", "0x263fef3fe76fd4075ac16271d5115d01206d3674"},
			kv{"params.0.value", nil},
		)},
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
			if err = tt.assert.AssertBytes(sugar.HandlerError(got.Serialize()).([]byte)); err != nil {
				t.Error(err)
			}
		})
	}
}

type h map[string]interface{}

func marshal(x interface{}) []byte {
	return sugar.HandlerError(json.Marshal(x)).([]byte)
}

func Test_decoratePrefix(t *testing.T) {
	fmt.Println(decoratePrefix("041a"))
}
