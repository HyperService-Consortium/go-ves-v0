package bni

import (
	"encoding/json"

	merkleproof "github.com/HyperService-Consortium/go-uip/merkle-proof"
	uiptypes "github.com/HyperService-Consortium/go-uip/uiptypes"
	nsbclient "github.com/HyperService-Consortium/go-ves/lib/net/nsb-client"
)

type MerkleProofInfo struct {
	Proof [][]byte `json:"proof"`
	Key   []byte   `json:"key"`
	Value []byte   `json:"value"`
}

func (bn *BN) GetTransactionProof(chainID uint64, blockID []byte, additional []byte) (uiptypes.MerkleProof, error) {
	cinfo, err := bn.dns.GetChainInfo(chainID)

	if err != nil {
		return nil, err
	}

	resp, err := nsbclient.NewNSBClient(cinfo.GetChainHost()).GetProof(additional, `"prove_on_tx_trie"`)

	if err != nil {
		return nil, err
	}

	var info MerkleProofInfo
	err = json.Unmarshal([]byte(resp.Info), &info)
	if err != nil {
		return nil, err
	}
	return merkleproof.NewMPTUsingKeccak256(info.Proof, info.Key, info.Value), nil
}
