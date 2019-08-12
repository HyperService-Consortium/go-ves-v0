package bni

import (
	"encoding/json"

	merkleproof "github.com/Myriad-Dreamin/go-uip/merkle-proof"
	uiptypes "github.com/Myriad-Dreamin/go-uip/types"
	nsbclient "github.com/Myriad-Dreamin/go-ves/lib/net/nsb-client"

	chaininfo "github.com/Myriad-Dreamin/go-uip/temporary-chain-info"
)

type MerkleProofInfo struct {
	Proof [][]byte `json:"proof"`
	Key   []byte   `json:"key"`
	Value []byte   `json:"value"`
}

func (bn *BN) GetTransactionProof(chainID uint64, blockID []byte, additional []byte) (uiptypes.MerkleProof, error) {
	cinfo, err := chaininfo.SearchChainId(chainID)

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

func (bn *BN) WaitForTransact(chainID uint64, receipt []byte, opt *uiptypes.WaitOption) ([]byte, []byte, error) {
	// todo

	var info RTxInfo

	json.Unmarshal(receipt, &info)

	return nil, info.transactionReceipt, nil
	// cinfo, err := SearchChainId(chainID)
	// if err != nil {
	// 	return nil, err
	// }
	// worker := nsbclient.NewNSBClient(cinfo.GetHost())
	// ddl := time.Now().Add(timeout)
	// for time.Now().Before(ddl) {
	// 	tx, err := worker.GetProof(receipt, `"prove_on_tx_trie"`)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	fmt.Println(string(tx))
	// 	if gjson.GetBytes(tx, "blockNumber").Type != gjson.Null {
	// 		b, _ := hex.DecodeString(gjson.GetBytes(tx, "blockHash").String()[2:])
	// 		return b, nil
	// 	}
	// 	time.Sleep(time.Millisecond * 500)
	//
	// }
	// return nil, errors.New("timeout")
}
