package nsbcli

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/big"

	appl "github.com/HyperServiceOne/NSB/application"
	cmn "github.com/HyperServiceOne/NSB/common"
	ISC "github.com/HyperServiceOne/NSB/contract/isc"
	iscTransactionIntent "github.com/HyperServiceOne/NSB/contract/isc/transaction"
	nsbmath "github.com/HyperServiceOne/NSB/math"
	uiptypes "github.com/Myriad-Dreamin/go-uip/types"
)

func (nc *NSBClient) CreateISC(
	user uiptypes.Signer,
	funds []uint32, iscOwners [][]byte,
	bytesTransactionIntents [][]byte,
	vesSig []byte,
) ([]byte, error) {
	var txHeader cmn.TransactionHeader
	var buf = bytes.NewBuffer(make([]byte, mxBytes))
	buf.Reset()
	// fmt.Println(string(buf.Bytes()))
	var transactionIntents []*iscTransactionIntent.TransactionIntent
	var txm map[string]interface{}
	for idx, txb := range bytesTransactionIntents {
		err := json.Unmarshal(txb, &txm)
		if err != nil {
			return nil, err
		}
		var txi = new(iscTransactionIntent.TransactionIntent)
		if txm["src"] == nil && txm["from"] == nil {
			return nil, errNilSrc
		}
		if txm["src"] != nil {
			txi.Fr, err = base64.StdEncoding.DecodeString(txm["src"].(string))
			if err != nil {
				return nil, err
			}
		} else {
			txi.Fr, err = base64.StdEncoding.DecodeString(txm["from"].(string))
			if err != nil {
				return nil, err
			}
		}
		if txm["dst"] != nil {
			txi.To, err = base64.StdEncoding.DecodeString(txm["dst"].(string))
			if err != nil {
				return nil, err
			}
		} else if txm["from"] != nil {
			txi.To, err = base64.StdEncoding.DecodeString(txm["from"].(string))
			if err != nil {
				return nil, err
			}
		}
		if txm["meta"] != nil {
			txi.Meta, err = base64.StdEncoding.DecodeString(txm["meta"].(string))
			if err != nil {
				return nil, err
			}
		}
		txi.Seq = nsbmath.NewUint256FromBigInt(big.NewInt(int64(idx)))
		if txm["amt"] != nil {
			b, _ := hex.DecodeString(txm["amt"].(string))
			txi.Amt = nsbmath.NewUint256FromBytes(b)
		} else {
			txi.Amt = nsbmath.NewUint256FromBytes([]byte{0})
		}
		transactionIntents = append(transactionIntents, txi)
		// fmt.Println("encoding", txm)
	}

	err := nc.createISC(buf, funds, iscOwners, transactionIntents, vesSig)
	if err != nil {
		return nil, err
	}
	txHeader.Data = buf.Bytes()
	txHeader.From = user.GetPublicKey()

	nonce := make([]byte, 32)
	_, err = rand.Read(nonce)
	if err != nil {
		return nil, err
	}
	txHeader.Nonce = nsbmath.NewUint256FromBytes(nonce)
	txHeader.Value = nsbmath.NewUint256FromBytes([]byte{0})
	// bug: buf.Reset()
	buf = bytes.NewBuffer(make([]byte, mxBytes))

	buf.Write(txHeader.From)
	buf.Write(txHeader.ContractAddress)
	buf.Write(txHeader.Data)
	buf.Write(txHeader.Value.Bytes())
	buf.Write(txHeader.Nonce.Bytes())
	txHeader.Signature = user.Sign(buf.Bytes()).Bytes()
	ret, err := nc.sendContractTx([]byte("createContract"), []byte("isc"), &txHeader)
	fmt.Println("create", PretiStruct(ret), err)
	if err != nil {
		return nil, err
	}
	return ret.DeliverTx.Data, nil
}

func (nc *NSBClient) CreateISCAsync(
	user uiptypes.Signer,
	funds []uint32, iscOwners [][]byte,
	bytesTransactionIntents [][]byte,
	vesSig []byte,
	option *AsyncOption,
) ([]byte, error) {
	var txHeader cmn.TransactionHeader
	var buf = bytes.NewBuffer(make([]byte, mxBytes))
	buf.Reset()
	// fmt.Println(string(buf.Bytes()))
	var transactionIntents []*iscTransactionIntent.TransactionIntent
	var txm map[string]interface{}
	for idx, txb := range bytesTransactionIntents {
		err := json.Unmarshal(txb, &txm)
		if err != nil {
			return nil, err
		}
		var txi = new(iscTransactionIntent.TransactionIntent)
		if txm["src"] == nil && txm["from"] == nil {
			return nil, errNilSrc
		}
		if txm["src"] != nil {
			txi.Fr, err = base64.StdEncoding.DecodeString(txm["src"].(string))
			if err != nil {
				return nil, err
			}
		} else {
			txi.Fr, err = base64.StdEncoding.DecodeString(txm["from"].(string))
			if err != nil {
				return nil, err
			}
		}
		if txm["dst"] != nil {
			txi.To, err = base64.StdEncoding.DecodeString(txm["dst"].(string))
			if err != nil {
				return nil, err
			}
		} else if txm["from"] != nil {
			txi.To, err = base64.StdEncoding.DecodeString(txm["from"].(string))
			if err != nil {
				return nil, err
			}
		}
		if txm["meta"] != nil {
			txi.Meta, err = base64.StdEncoding.DecodeString(txm["meta"].(string))
			if err != nil {
				return nil, err
			}
		}
		txi.Seq = nsbmath.NewUint256FromBigInt(big.NewInt(int64(idx)))
		if txm["amt"] != nil {
			b, _ := hex.DecodeString(txm["amt"].(string))
			txi.Amt = nsbmath.NewUint256FromBytes(b)
		} else {
			txi.Amt = nsbmath.NewUint256FromBytes([]byte{0})
		}
		transactionIntents = append(transactionIntents, txi)
		// fmt.Println("encoding", txm)
	}

	err := nc.createISC(buf, funds, iscOwners, transactionIntents, vesSig)
	if err != nil {
		return nil, err
	}
	txHeader.Data = buf.Bytes()
	txHeader.From = user.GetPublicKey()

	nonce := make([]byte, 32)
	_, err = rand.Read(nonce)
	if err != nil {
		return nil, err
	}
	txHeader.Nonce = nsbmath.NewUint256FromBytes(nonce)
	txHeader.Value = nsbmath.NewUint256FromBytes([]byte{0})
	// bug: buf.Reset()
	buf = bytes.NewBuffer(make([]byte, mxBytes))

	buf.Write(txHeader.From)
	buf.Write(txHeader.ContractAddress)
	buf.Write(txHeader.Data)
	buf.Write(txHeader.Value.Bytes())
	buf.Write(txHeader.Nonce.Bytes())
	txHeader.Signature = user.Sign(buf.Bytes()).Bytes()
	ret, err := nc.sendContractTxAsync([]byte("createContract"), []byte("isc"), &txHeader, option)
	fmt.Println("create", PretiStruct(ret), err)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (nc *NSBClient) createISC(
	w io.Writer,
	funds []uint32, iscOwners [][]byte,
	transactionIntents []*iscTransactionIntent.TransactionIntent,
	vesSig []byte,
) error {
	var args ISC.ArgsCreateNewContract
	args.IscOwners = iscOwners
	args.Funds = funds
	args.TransactionIntents = transactionIntents
	args.VesSig = vesSig
	b, err := json.Marshal(args)
	if err != nil {
		return err
	}

	// fmt.Println(PretiStruct(args), b)
	_, err = w.Write(b)
	return err
}

type AddActionsBatcher struct {
	nc        *NSBClient
	user      uiptypes.Signer
	toAddress []byte
	argss     []appl.ArgsAddAction
}

func (batcher *AddActionsBatcher) Insert(
	iscAddress []byte, tid uint64, aid uint64, stype uint8,
	content []byte, signature []byte,
) *AddActionsBatcher {
	batcher.argss = append(batcher.argss, appl.ArgsAddAction{
		ISCAddress: iscAddress,
		Tid:        tid,
		Aid:        aid,
		Type:       stype,
		Content:    content,
		Signature:  signature,
	})
	return batcher
}

func (batcher *AddActionsBatcher) Commit() ([]byte, error) {

	var txHeader cmn.TransactionHeader
	var buf = bytes.NewBuffer(make([]byte, mxBytes))
	buf.Reset()

	var args appl.ArgsAddActions
	args.Args = batcher.argss

	var fap appl.FAPair
	var err error
	fap.FuncName = "addActions"
	fap.Args, err = json.Marshal(args)
	if err != nil {
		return nil, err
	}

	txHeader.Data, err = json.Marshal(fap)
	if err != nil {
		return nil, err
	}
	txHeader.ContractAddress = batcher.toAddress
	txHeader.From = batcher.user.GetPublicKey()

	nonce := make([]byte, 32)
	_, err = rand.Read(nonce)
	if err != nil {
		return nil, err
	}
	txHeader.Nonce = nsbmath.NewUint256FromBytes(nonce)
	txHeader.Value = nsbmath.NewUint256FromBytes([]byte{0})
	// bug: buf.Reset()
	buf = bytes.NewBuffer(make([]byte, mxBytes))

	buf.Write(txHeader.From)
	buf.Write(txHeader.ContractAddress)
	buf.Write(txHeader.Data)
	buf.Write(txHeader.Value.Bytes())
	buf.Write(txHeader.Nonce.Bytes())
	txHeader.Signature = batcher.user.Sign(buf.Bytes()).Bytes()
	_, err = batcher.nc.sendContractTx([]byte("systemCall"), []byte("system.action"), &txHeader)
	// fmt.Println(PretiStruct(ret), err)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (nc *NSBClient) AddActions(
	user uiptypes.Signer, toAddress []byte, predictNumbers int,
) *AddActionsBatcher {
	return &AddActionsBatcher{
		nc:        nc,
		user:      user,
		toAddress: toAddress,
		argss:     make([]appl.ArgsAddAction, 0, predictNumbers),
	}
}

func (nc *NSBClient) AddAction(
	user uiptypes.Signer, toAddress []byte,
	iscAddress []byte, tid uint64, aid uint64, stype uint8, content []byte, signature []byte,
) ([]byte, error) {
	var txHeader cmn.TransactionHeader
	var buf = bytes.NewBuffer(make([]byte, mxBytes))
	buf.Reset()
	// fmt.Println(string(buf.Bytes()))
	err := nc.addAction(buf, iscAddress, tid, aid, stype, content, signature)
	if err != nil {
		return nil, err
	}

	var fap appl.FAPair
	fap.FuncName = "addAction"
	fap.Args = buf.Bytes()
	txHeader.Data, err = json.Marshal(fap)
	if err != nil {
		return nil, err
	}
	txHeader.ContractAddress = toAddress
	txHeader.From = user.GetPublicKey()

	nonce := make([]byte, 32)
	_, err = rand.Read(nonce)
	if err != nil {
		return nil, err
	}
	txHeader.Nonce = nsbmath.NewUint256FromBytes(nonce)
	txHeader.Value = nsbmath.NewUint256FromBytes([]byte{0})
	// bug: buf.Reset()
	buf = bytes.NewBuffer(make([]byte, mxBytes))

	buf.Write(txHeader.From)
	buf.Write(txHeader.ContractAddress)
	buf.Write(txHeader.Data)
	buf.Write(txHeader.Value.Bytes())
	buf.Write(txHeader.Nonce.Bytes())
	txHeader.Signature = user.Sign(buf.Bytes()).Bytes()
	_, err = nc.sendContractTx([]byte("systemCall"), []byte("system.action"), &txHeader)
	// fmt.Println(PretiStruct(ret), err)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (nc *NSBClient) addAction(
	w io.Writer,
	iscAddress []byte, tid uint64, aid uint64, stype uint8, content []byte, signature []byte,
) error {
	var args appl.ArgsAddAction
	args.ISCAddress = iscAddress
	args.Tid = tid
	args.Aid = aid
	args.Type = stype
	args.Content = content
	args.Signature = signature
	b, err := json.Marshal(args)
	if err != nil {
		return err
	}

	// fmt.Println(PretiStruct(args), b)
	_, err = w.Write(b)
	return err
}

func (nc *NSBClient) GetAction(
	user uiptypes.Signer, toAddress []byte,
	iscAddress []byte, tid uint64, aid uint64,
) ([]byte, error) {
	var txHeader cmn.TransactionHeader
	var buf = bytes.NewBuffer(make([]byte, mxBytes))
	buf.Reset()
	// fmt.Println(string(buf.Bytes()))
	err := nc.getAction(buf, iscAddress, tid, aid)
	if err != nil {
		return nil, err
	}

	var fap appl.FAPair
	fap.FuncName = "getAction"
	fap.Args = buf.Bytes()
	txHeader.Data, err = json.Marshal(fap)
	if err != nil {
		return nil, err
	}
	txHeader.ContractAddress = toAddress
	txHeader.From = user.GetPublicKey()

	nonce := make([]byte, 32)
	_, err = rand.Read(nonce)
	if err != nil {
		return nil, err
	}
	txHeader.Nonce = nsbmath.NewUint256FromBytes(nonce)
	txHeader.Value = nsbmath.NewUint256FromBytes([]byte{0})
	// bug: buf.Reset()
	buf = bytes.NewBuffer(make([]byte, mxBytes))

	buf.Write(txHeader.From)
	buf.Write(txHeader.ContractAddress)
	buf.Write(txHeader.Data)
	buf.Write(txHeader.Value.Bytes())
	buf.Write(txHeader.Nonce.Bytes())
	txHeader.Signature = user.Sign(buf.Bytes()).Bytes()
	_, err = nc.sendContractTx([]byte("systemCall"), []byte("system.action"), &txHeader)
	// fmt.Println(PretiStruct(ret), err)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (nc *NSBClient) getAction(
	w io.Writer,
	iscAddress []byte, tid uint64, aid uint64,
) error {
	var args appl.ArgsAddAction
	args.ISCAddress = iscAddress
	args.Tid = tid
	args.Aid = aid
	b, err := json.Marshal(args)
	if err != nil {
		return err
	}

	// fmt.Println(PretiStruct(args), b)
	_, err = w.Write(b)
	return err
}

func (nc *NSBClient) AddMerkleProof(
	user uiptypes.Signer, toAddress []byte,
	merkletype uint16, rootHash, proof, key, value []byte,
) (*ResultInfo, error) {
	var txHeader cmn.TransactionHeader
	var buf = bytes.NewBuffer(make([]byte, mxBytes))
	buf.Reset()
	// fmt.Println(string(buf.Bytes()))
	err := nc.addMerkleProof(buf, merkletype, rootHash, proof, key, value)
	if err != nil {
		return nil, err
	}

	var fap appl.FAPair
	fap.FuncName = "validateMerkleProof"
	fap.Args = buf.Bytes()
	txHeader.Data, err = json.Marshal(fap)
	if err != nil {
		return nil, err
	}
	txHeader.ContractAddress = toAddress
	txHeader.From = user.GetPublicKey()

	nonce := make([]byte, 32)
	_, err = rand.Read(nonce)
	if err != nil {
		return nil, err
	}
	txHeader.Nonce = nsbmath.NewUint256FromBytes(nonce)
	txHeader.Value = nsbmath.NewUint256FromBytes([]byte{0})
	// bug: buf.Reset()
	buf = bytes.NewBuffer(make([]byte, mxBytes))

	buf.Write(txHeader.From)
	buf.Write(txHeader.ContractAddress)
	buf.Write(txHeader.Data)
	buf.Write(txHeader.Value.Bytes())
	buf.Write(txHeader.Nonce.Bytes())
	txHeader.Signature = user.Sign(buf.Bytes()).Bytes()
	ret, err := nc.sendContractTx([]byte("systemCall"), []byte("system.merkleproof"), &txHeader)
	// fmt.Println(PretiStruct(ret), err)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (nc *NSBClient) addMerkleProof(
	w io.Writer,
	merkletype uint16, rootHash []byte, proof []byte, key []byte, value []byte,
) error {
	var args appl.ArgsValidateMerkleProof

	args.Type = merkletype
	args.RootHash = rootHash

	args.Proof = proof
	args.Key = key
	args.Value = value
	b, err := json.Marshal(args)
	if err != nil {
		return err
	}

	// fmt.Println(PretiStruct(args), b)
	_, err = w.Write(b)
	return err
}

/*
iscAddress []byte, cid uint64, bid uint64,
rootHash []byte,
*/
func (nc *NSBClient) AddBlockCheck(
	user uiptypes.Signer, toAddress []byte,
	chainID uint64, blockID, rootHash []byte, rcType uint8,
) (*ResultInfo, error) {
	var txHeader cmn.TransactionHeader
	var buf = bytes.NewBuffer(make([]byte, mxBytes))
	buf.Reset()
	// fmt.Println(string(buf.Bytes()))
	err := nc.addBlockCheck(buf, chainID, blockID, rootHash, rcType)
	if err != nil {
		return nil, err
	}

	var fap appl.FAPair
	fap.FuncName = "addBlockCheck"
	fap.Args = buf.Bytes()
	txHeader.Data, err = json.Marshal(fap)
	if err != nil {
		return nil, err
	}
	txHeader.ContractAddress = toAddress
	txHeader.From = user.GetPublicKey()

	nonce := make([]byte, 32)
	_, err = rand.Read(nonce)
	if err != nil {
		return nil, err
	}
	txHeader.Nonce = nsbmath.NewUint256FromBytes(nonce)
	txHeader.Value = nsbmath.NewUint256FromBytes([]byte{0})
	// bug: buf.Reset()
	buf = bytes.NewBuffer(make([]byte, mxBytes))

	buf.Write(txHeader.From)
	buf.Write(txHeader.ContractAddress)
	buf.Write(txHeader.Data)
	buf.Write(txHeader.Value.Bytes())
	buf.Write(txHeader.Nonce.Bytes())
	txHeader.Signature = user.Sign(buf.Bytes()).Bytes()
	ret, err := nc.sendContractTx([]byte("systemCall"), []byte("system.merkleproof"), &txHeader)
	// fmt.Println(PretiStruct(ret), err)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (nc *NSBClient) addBlockCheck(
	w io.Writer,
	chainID uint64, blockID, rootHash []byte, rtType uint8,
) error {
	var args appl.ArgsAddBlockCheck
	args.ChainID = chainID
	args.BlockID = blockID
	args.RootHash = rootHash
	args.RtType = rtType
	b, err := json.Marshal(args)
	if err != nil {
		return err
	}

	// fmt.Println(PretiStruct(args), b)
	_, err = w.Write(b)
	return err
}

func (nc *NSBClient) GetMerkleProof(
	user uiptypes.Signer, toAddress []byte,
	merkleProofType uint16, rootHash, key []byte, chainID uint64, blockID []byte, rcType uint8,
) ([]byte, error) {
	var txHeader cmn.TransactionHeader
	var buf = bytes.NewBuffer(make([]byte, mxBytes))
	buf.Reset()
	// fmt.Println(string(buf.Bytes()))
	err := nc.getMerkleProof(buf, merkleProofType, rootHash, key, chainID, blockID, rcType)
	if err != nil {
		return nil, err
	}

	var fap appl.FAPair
	fap.FuncName = "getMerkleProof"
	fap.Args = buf.Bytes()
	txHeader.Data, err = json.Marshal(fap)
	if err != nil {
		return nil, err
	}
	txHeader.ContractAddress = toAddress
	txHeader.From = user.GetPublicKey()

	nonce := make([]byte, 32)
	_, err = rand.Read(nonce)
	if err != nil {
		return nil, err
	}
	txHeader.Nonce = nsbmath.NewUint256FromBytes(nonce)
	txHeader.Value = nsbmath.NewUint256FromBytes([]byte{0})
	// bug: buf.Reset()
	buf = bytes.NewBuffer(make([]byte, mxBytes))

	buf.Write(txHeader.From)
	buf.Write(txHeader.ContractAddress)
	buf.Write(txHeader.Data)
	buf.Write(txHeader.Value.Bytes())
	buf.Write(txHeader.Nonce.Bytes())
	txHeader.Signature = user.Sign(buf.Bytes()).Bytes()
	_, err = nc.sendContractTx([]byte("systemCall"), []byte("system.merkleproof"), &txHeader)
	// fmt.Println(PretiStruct(ret), err)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (nc *NSBClient) getMerkleProof(
	w io.Writer,
	merkleProofType uint16, rootHash, key []byte, chainID uint64, blockID []byte, rtType uint8,
) error {
	var args appl.ArgsGetMerkleProof
	args.Type = merkleProofType
	args.RootHash = rootHash
	args.Key = key
	args.ChainID = chainID
	args.BlockID = blockID
	args.RtType = rtType
	b, err := json.Marshal(args)
	if err != nil {
		return err
	}

	// fmt.Println(PretiStruct(args), b)
	_, err = w.Write(b)
	return err
}

//
// func (nc *NSBClient) AddMerkleProof(
// 	user uiptypes.Signer, toAddress []byte,
// 	iscAddress []byte, cid uint64, bid uint64,
// 	rootHash []byte, key []byte, value []byte, proof []byte,
// ) ([]byte, error) {
// 	var txHeader cmn.TransactionHeader
// 	var buf = bytes.NewBuffer(make([]byte, mxBytes))
// 	buf.Reset()
// 	// fmt.Println(string(buf.Bytes()))
// 	err := nc.addMerkleProof(buf, iscAddress, cid, bid, rootHash, key, value, proof)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	var fap appl.FAPair
// 	fap.FuncName = "validateMerkleProof"
// 	fap.Args = buf.Bytes()
// 	txHeader.Data, err = json.Marshal(fap)
// 	if err != nil {
// 		return nil, err
// 	}
// 	txHeader.ContractAddress = toAddress
// 	txHeader.From = user.GetPublicKey()
//
// 	nonce := make([]byte, 32)
// 	_, err = rand.Read(nonce)
// 	if err != nil {
// 		return nil, err
// 	}
// 	txHeader.Nonce = nmath.NewUint256FromBytes(nonce)
// 	txHeader.Value = nmath.NewUint256FromBytes([]byte{0})
// 	// bug: buf.Reset()
// 	buf = bytes.NewBuffer(make([]byte, mxBytes))
//
// 	buf.Write(txHeader.From)
// 	buf.Write(txHeader.ContractAddress)
// 	buf.Write(txHeader.Data)
// 	buf.Write(txHeader.Value.Bytes())
// 	buf.Write(txHeader.Nonce.Bytes())
// 	txHeader.Signature = user.Sign(buf.Bytes()).Bytes()
// 	_, err = nc.sendContractTx([]byte("systemCall"), []byte("system.merkleproof"), &txHeader)
// 	// fmt.Println(PretiJson(ret), err)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return nil, nil
// }
//
// type ArgsAddMerkleProof struct {
// 	ISCAddress []byte `json:"1"`
// 	Cid        uint64 `json:"2"`
// 	Bid        uint64 `json:"3"`
// 	RootHash   []byte `json:"4"`
// 	Key        []byte `json:"5"`
// 	Value      []byte `json:"6"`
// 	Proof      []byte `json:"7"`
// }
//
// func (nc *NSBClient) addMerkleProof(
// 	w io.Writer,
// 	iscAddress []byte, cid uint64, bid uint64,
// 	rootHash []byte, key []byte, value []byte, proof []byte,
// ) error {
// 	var args ArgsAddMerkleProof
// 	args.ISCAddress = iscAddress
// 	args.Cid = cid
// 	args.Bid = bid
// 	args.RootHash = rootHash
// 	args.Key = key
// 	args.Value = value
// 	args.Proof = proof
// 	b, err := json.Marshal(args)
// 	if err != nil {
// 		return err
// 	}
//
// 	// fmt.Println(PretiJson(args), b)
// 	_, err = w.Write(b)
// 	return err
// }
//
// type ArgsGetMerkleProof struct {
// 	ISCAddress []byte `json:"1"`
// 	Cid        uint64 `json:"2"`
// 	Bid        uint64 `json:"3"`
// 	RootHash   []byte `json:"4"`
// 	Key        []byte `json:"5"`
// }
//
// func (nc *NSBClient) GetMerkleProof(
// 	user uiptypes.Signer, toAddress []byte,
// 	iscAddress []byte, cid uint64, bid uint64,
// ) ([]byte, error) {
// 	var txHeader cmn.TransactionHeader
// 	var buf = bytes.NewBuffer(make([]byte, mxBytes))
// 	buf.Reset()
// 	// fmt.Println(string(buf.Bytes()))
// 	err := nc.getMerkleProof(buf, iscAddress, cid, bid)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	var fap appl.FAPair
// 	fap.FuncName = "getMerkleProof"
// 	fap.Args = buf.Bytes()
// 	txHeader.Data, err = json.Marshal(fap)
// 	if err != nil {
// 		return nil, err
// 	}
// 	txHeader.ContractAddress = toAddress
// 	txHeader.From = user.GetPublicKey()
//
// 	nonce := make([]byte, 32)
// 	_, err = rand.Read(nonce)
// 	if err != nil {
// 		return nil, err
// 	}
// 	txHeader.Nonce = nmath.NewUint256FromBytes(nonce)
// 	txHeader.Value = nmath.NewUint256FromBytes([]byte{0})
// 	// bug: buf.Reset()
// 	buf = bytes.NewBuffer(make([]byte, mxBytes))
//
// 	buf.Write(txHeader.From)
// 	buf.Write(txHeader.ContractAddress)
// 	buf.Write(txHeader.Data)
// 	buf.Write(txHeader.Value.Bytes())
// 	buf.Write(txHeader.Nonce.Bytes())
// 	txHeader.Signature = user.Sign(buf.Bytes()).Bytes()
// 	_, err = nc.sendContractTx([]byte("systemCall"), []byte("system.merkleproof"), &txHeader)
// 	// fmt.Println(PretiJson(ret), err)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return nil, nil
// }
//
// func (nc *NSBClient) getMerkleProof(
// 	w io.Writer,
// 	iscAddress []byte, cid uint64, bid uint64,
// ) error {
// 	var args ArgsGetMerkleProof
// 	args.ISCAddress = iscAddress
// 	args.Cid = cid
// 	args.Bid = bid
// 	b, err := json.Marshal(args)
// 	if err != nil {
// 		return err
// 	}
//
// 	// fmt.Println(PretiJson(args), b)
// 	_, err = w.Write(b)
// 	return err
// }

func (nc *NSBClient) UpdateTxInfo(
	user uiptypes.Signer, contractAddress []byte,
	tid uint64, transactionIntent *iscTransactionIntent.TransactionIntent,
) ([]byte, error) {
	var txHeader cmn.TransactionHeader
	var buf = bytes.NewBuffer(make([]byte, mxBytes))
	buf.Reset()
	// fmt.Println(string(buf.Bytes()))
	err := nc.updateTxInfo(buf, tid, transactionIntent)
	if err != nil {
		return nil, err
	}

	var fap appl.FAPair
	fap.FuncName = "UpdateTxInfo"
	fap.Args = buf.Bytes()
	txHeader.Data, err = json.Marshal(fap)
	if err != nil {
		return nil, err
	}
	txHeader.ContractAddress = contractAddress
	txHeader.From = user.GetPublicKey()

	nonce := make([]byte, 32)
	_, err = rand.Read(nonce)
	if err != nil {
		return nil, err
	}
	txHeader.Nonce = nsbmath.NewUint256FromBytes(nonce)
	txHeader.Value = nsbmath.NewUint256FromBytes([]byte{0})
	// bug: buf.Reset()
	buf = bytes.NewBuffer(make([]byte, mxBytes))

	buf.Write(txHeader.From)
	buf.Write(txHeader.ContractAddress)
	buf.Write(txHeader.Data)
	buf.Write(txHeader.Value.Bytes())
	buf.Write(txHeader.Nonce.Bytes())
	txHeader.Signature = user.Sign(buf.Bytes()).Bytes()
	_, err = nc.sendContractTx([]byte("sendTransaction"), []byte("isc"), &txHeader)
	// fmt.Println(PretiStruct(ret), err)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (nc *NSBClient) updateTxInfo(
	w io.Writer,
	tid uint64, transactionIntent *iscTransactionIntent.TransactionIntent,
) error {
	var args ISC.ArgsUpdateTxInfo
	args.Tid = tid
	args.TransactionIntent = transactionIntent
	b, err := json.Marshal(args)
	if err != nil {
		return err
	}

	// fmt.Println(PretiStruct(args), b)
	_, err = w.Write(b)
	return err
}

func (nc *NSBClient) FreezeInfo(
	user uiptypes.Signer, contractAddress []byte,
	tid uint64,
) ([]byte, error) {
	var txHeader cmn.TransactionHeader
	var buf = bytes.NewBuffer(make([]byte, mxBytes))
	buf.Reset()
	// fmt.Println(string(buf.Bytes()))
	err := nc.freezeInfo(buf, tid)
	if err != nil {
		return nil, err
	}

	var fap appl.FAPair
	fap.FuncName = "FreezeInfo"
	fap.Args = buf.Bytes()
	txHeader.Data, err = json.Marshal(fap)
	if err != nil {
		return nil, err
	}
	txHeader.ContractAddress = contractAddress
	txHeader.From = user.GetPublicKey()

	nonce := make([]byte, 32)
	_, err = rand.Read(nonce)
	if err != nil {
		return nil, err
	}
	txHeader.Nonce = nsbmath.NewUint256FromBytes(nonce)
	txHeader.Value = nsbmath.NewUint256FromBytes([]byte{0})
	// bug: buf.Reset()
	buf = bytes.NewBuffer(make([]byte, mxBytes))

	buf.Write(txHeader.From)
	buf.Write(txHeader.ContractAddress)
	buf.Write(txHeader.Data)
	buf.Write(txHeader.Value.Bytes())
	buf.Write(txHeader.Nonce.Bytes())
	txHeader.Signature = user.Sign(buf.Bytes()).Bytes()
	_, err = nc.sendContractTx([]byte("sendTransaction"), []byte("isc"), &txHeader)
	// fmt.Println(PretiStruct(ret), err)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (nc *NSBClient) freezeInfo(
	w io.Writer,
	tid uint64,
) error {
	var args ISC.ArgsFreezeInfo
	args.Tid = tid
	b, err := json.Marshal(args)
	if err != nil {
		return err
	}

	// fmt.Println(PretiStruct(args), b)
	_, err = w.Write(b)
	return err
}
