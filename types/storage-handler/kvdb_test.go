package storage_handler

import (
	"fmt"
	"github.com/HyperService-Consortium/go-uip/const/value_type"
	"github.com/HyperService-Consortium/go-uip/uiptypes"
	map_index "github.com/HyperService-Consortium/go-ves/lib/database/map-index"
	"github.com/HyperService-Consortium/go-ves/types"
	"math/big"
	"reflect"
	"testing"
)

func TestClone(t *testing.T) {
	fmt.Println(string(clone([]byte{'a', 'b', 'c'})), len(clone([]byte{'a', 'b', 'c'})))
}

func TestCloneWithLen(t *testing.T) {
	fmt.Println(string(cloneWithLen([]byte{'a', 'b', 'c'}, 7)), len(cloneWithLen([]byte{'a', 'b', 'c'}, 7)))
}

func TestDecorate(t *testing.T) {
	fmt.Println(string(decorate([]byte{'a', 'b', 'c'}, []byte{'o', 'r', 'z'})))
}

func TestDatabase_SetStorageOf(t *testing.T) {
	var index = map_index.NewMapIndex()

	type args struct {
		index           types.Index
		chainID         uiptypes.ChainID
		typeID          uiptypes.TypeID
		contractAddress uiptypes.ContractAddress
		pos             []byte
		description     []byte
		variable        uiptypes.Variable
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"test_easy", args{
			index:           index,
			chainID:         1,
			typeID:          value_type.Bool,
			contractAddress: []byte{13},
			pos:             []byte{13},
			description:     []byte{13},
			variable: variable{
				Type:  value_type.Bool,
				Value: true,
			},
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Database{}
			if err := g.SetStorageOf(tt.args.index, tt.args.chainID, tt.args.typeID, tt.args.contractAddress, tt.args.pos, tt.args.description, tt.args.variable); (err != nil) != tt.wantErr {
				t.Errorf("SetStorageOf() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDatabase_GetStorageAt(t *testing.T) {

	var index = map_index.NewMapIndex()

	type args struct {
		index           types.Index
		chainID         uiptypes.ChainID
		typeID          uiptypes.TypeID
		contractAddress uiptypes.ContractAddress
		pos             []byte
		description     []byte
	}
	g := &Database{}
	err := g.SetStorageOf(index, 1, value_type.Bool, []byte{13}, []byte{13}, []byte{13}, variable{
		Type:  value_type.Bool,
		Value: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	err = g.SetStorageOf(index, 1, value_type.Bytes, []byte{13}, []byte{12}, []byte{13}, variable{
		Type:  value_type.Bytes,
		Value: []byte("hello world"),
	})
	if err != nil {
		t.Fatal(err)
	}
	err = g.SetStorageOf(index, 1, value_type.String, []byte{13}, []byte{131}, []byte{13}, variable{
		Type:  value_type.String,
		Value: "hello world!!",
	})
	if err != nil {
		t.Fatal(err)
	}
	err = g.SetStorageOf(index, 1, value_type.Uint256, []byte{13}, []byte{132}, []byte{13}, variable{
		Type:  value_type.Uint256,
		Value: big.NewInt(23333),
	})
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		args    args
		want    uiptypes.Variable
		wantErr bool
	}{
		{"test_easy", args{
			index:           index,
			chainID:         1,
			typeID:          value_type.Bool,
			contractAddress: []byte{13},
			pos:             []byte{13},
			description:     []byte{13},
		}, variable{
			Type:  value_type.Bool,
			Value: true,
		}, false},
		{"test_slice", args{
			index:           index,
			chainID:         1,
			typeID:          value_type.Bytes,
			contractAddress: []byte{13},
			pos:             []byte{12},
			description:     []byte{13},
		}, variable{
			Type:  value_type.Bytes,
			Value: []byte("hello world"),
		}, false},
		{"test_string", args{
			index:           index,
			chainID:         1,
			typeID:          value_type.String,
			contractAddress: []byte{13},
			pos:             []byte{131},
			description:     []byte{13},
		}, variable{
			Type:  value_type.String,
			Value: "hello world!!",
		}, false},
		{"test_big_int", args{
			index:           index,
			chainID:         1,
			typeID:          value_type.Uint256,
			contractAddress: []byte{13},
			pos:             []byte{132},
			description:     []byte{13},
		}, variable{
			Type:  value_type.Uint256,
			Value: big.NewInt(23333),
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Database{}
			got, err := g.GetStorageAt(tt.args.index, tt.args.chainID, tt.args.typeID, tt.args.contractAddress, tt.args.pos, tt.args.description)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetStorageAt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetStorageAt() got = %v, want %v", got, tt.want)
			}
		})
	}
}
