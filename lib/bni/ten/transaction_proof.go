package bni

import (
	"encoding/json"

	merkleproof "github.com/Myriad-Dreamin/go-uip/merkle-proof"
	uiptypes "github.com/Myriad-Dreamin/go-uip/types"
	nsbclient "github.com/Myriad-Dreamin/go-ves/lib/net/nsb-client"
)

type MerkleProofInfo struct {
	Proof [][]byte `json:"proof"`
	Key   []byte   `json:"key"`
	Value []byte   `json:"value"`
}

func (bn *BN) GetTransactionProof(chainID uint64, blockID []byte, additional []byte) (uiptypes.MerkleProof, error) {
	cinfo, err := SearchChainId(chainID)

	if err != nil {
		return nil, err
	}

	resp, err := nsbclient.NewNSBClient(cinfo.GetHost()).GetProof(additional, `"prove_on_tx_trie"`)

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
