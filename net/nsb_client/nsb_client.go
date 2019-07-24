package nsbcli

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"strings"
	"time"

	request "github.com/Myriad-Dreamin/go-ves/net/request"
	"github.com/tidwall/gjson"

	appl "github.com/HyperServiceOne/NSB/application"
	cmn "github.com/HyperServiceOne/NSB/common"
	ISC "github.com/HyperServiceOne/NSB/contract/isc"
	tx "github.com/HyperServiceOne/NSB/contract/isc/transaction"
	nmath "github.com/HyperServiceOne/NSB/math"
	mt19937 "github.com/Myriad-Dreamin/go-ves/math/mt19937"
	bytes_pool "github.com/Myriad-Dreamin/go-ves/net/bytes_pool"
	jsonrpc_client "github.com/Myriad-Dreamin/go-ves/net/rpc-client"
)

const (
	httpPrefix   = "http://"
	httpsPrefix  = "https://"
	maxBytesSize = 64 * 1024
)

func decorateHost(host string) string {
	if strings.HasPrefix(host, httpPrefix) || strings.HasPrefix(host, httpsPrefix) {
		return host
	}
	return httpPrefix + host
}

type NSBClient struct {
	handler    *request.RequestClient
	bufferPool *bytes_pool.BytesPool
}

// todo: test invalid json
func (nc *NSBClient) preloadJsonResponse(bb io.ReadCloser) ([]byte, error) {

	var b = nc.bufferPool.Get().([]byte)
	defer nc.bufferPool.Put(b)

	_, err := bb.Read(b)
	if err != nil && err != io.EOF {
		return nil, err
	}
	bb.Close()

	var jm = gjson.ParseBytes(b)
	if s := jm.Get("jsonrpc"); !s.Exists() || s.String() != "2.0" {
		return nil, errors.New("reject ret that is not jsonrpc: 2.0")
	}
	if s := jm.Get("error"); s.Exists() {
		return nil, jsonrpc_client.FromGJsonResultError(s)
	}
	if s := jm.Get("result"); s.Exists() {
		if s.Index > 0 {
			return b[s.Index : s.Index+len(s.Raw)], nil
		}
	}
	return nil, errors.New("bad format of jsonrpc")
}

func NewNSBClient(host string) *NSBClient {
	return &NSBClient{
		handler:    request.NewRequestClient(decorateHost(host)),
		bufferPool: bytes_pool.NewBytesPool(maxBytesSize),
	}
}

func (nc *NSBClient) GetAbciInfo() (*AbciInfoResponse, error) {
	b, err := nc.handler.Group("/abci_info").GetWithParams(request.Param{})
	if err != nil {
		return nil, err
	}
	var bb []byte
	bb, err = nc.preloadJsonResponse(b)
	if err != nil {
		return nil, err
	}
	var a AbciInfo
	err = json.Unmarshal(bb, &a)
	if err != nil {
		return nil, err
	}
	return a.Response, nil
}

func (nc *NSBClient) GetBlock(id int64) (*BlockInfo, error) {
	b, err := nc.handler.Group("/block").GetWithParams(request.Param{
		"height": id,
	})
	if err != nil {
		return nil, err
	}
	var bb []byte
	bb, err = nc.preloadJsonResponse(b)
	if err != nil {
		return nil, err
	}
	var a BlockInfo
	err = json.Unmarshal(bb, &a)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (nc *NSBClient) GetBlocks(rangeL, rangeR int64) (*BlocksInfo, error) {
	b, err := nc.handler.Group("/blockchain").GetWithParams(request.Param{
		"minHeight": rangeL,
		"maxHeight": rangeR,
	})
	if err != nil {
		return nil, err
	}
	var bb []byte
	bb, err = nc.preloadJsonResponse(b)
	if err != nil {
		return nil, err
	}
	var a BlocksInfo
	err = json.Unmarshal(bb, &a)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (nc *NSBClient) GetBlockResults(id int64) (*BlockResultsInfo, error) {
	b, err := nc.handler.Group("/block_results").GetWithParams(request.Param{
		"height": id,
	})
	if err != nil {
		return nil, err
	}
	var bb []byte
	bb, err = nc.preloadJsonResponse(b)
	if err != nil {
		return nil, err
	}
	var a BlockResultsInfo
	err = json.Unmarshal(bb, &a)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (nc *NSBClient) GetCommitInfo(id int64) (*CommitInfo, error) {
	b, err := nc.handler.Group("/commit").GetWithParams(request.Param{
		"height": id,
	})
	if err != nil {
		return nil, err
	}
	var bb []byte
	bb, err = nc.preloadJsonResponse(b)
	if err != nil {
		return nil, err
	}
	var a CommitInfo
	err = json.Unmarshal(bb, &a)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (nc *NSBClient) GetConsensusParamsInfo(id int64) (*ConsensusParamsInfo, error) {
	b, err := nc.handler.Group("/consensus_params").GetWithParams(request.Param{
		"height": id,
	})
	if err != nil {
		return nil, err
	}
	var bb []byte
	bb, err = nc.preloadJsonResponse(b)
	if err != nil {
		return nil, err
	}
	var a ConsensusParamsInfo
	err = json.Unmarshal(bb, &a)
	if err != nil {
		return nil, err
	}
	return &a, nil
}
func (nc *NSBClient) BroadcastTxCommit(body []byte) (*ResultInfo, error) {
	b, err := nc.handler.Group("/broadcast_tx_commit").GetWithParams(request.Param{
		"tx": "0x" + hex.EncodeToString(body),
	})
	if err != nil {
		return nil, err
	}
	var bb []byte
	bb, err = nc.preloadJsonResponse(b)
	if err != nil {
		return nil, err
	}
	var a ResultInfo
	err = json.Unmarshal(bb, &a)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (nc *NSBClient) GetConsensusState() (*ConsensusStateInfo, error) {
	b, err := nc.handler.Group("/consensus_state").Get()
	if err != nil {
		return nil, err
	}
	var bb []byte
	bb, err = nc.preloadJsonResponse(b)
	if err != nil {
		return nil, err
	}
	var a ConsensusStateInfo
	err = json.Unmarshal(bb, &a)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (nc *NSBClient) GetGenesis() (*GenesisInfo, error) {
	b, err := nc.handler.Group("/genesis").Get()
	if err != nil {
		return nil, err
	}
	var bb []byte
	bb, err = nc.preloadJsonResponse(b)
	if err != nil {
		return nil, err
	}
	var a GenesisInfo
	err = json.Unmarshal(bb, &a)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

//NOT DONE
func (nc *NSBClient) GetHealth() (interface{}, error) {
	b, err := nc.handler.Group("/health").Get()
	if err != nil {
		return nil, err
	}
	var bb []byte
	bb, err = nc.preloadJsonResponse(b)
	if err != nil {
		return nil, err
	}
	var a interface{}
	err = json.Unmarshal(bb, &a)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (nc *NSBClient) GetNetInfo() (*NetInfo, error) {
	b, err := nc.handler.Group("/net_info").Get()
	if err != nil {
		return nil, err
	}
	var bb []byte
	bb, err = nc.preloadJsonResponse(b)
	if err != nil {
		return nil, err
	}
	var a NetInfo
	err = json.Unmarshal(bb, &a)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (nc *NSBClient) GetNumUnconfirmedTxs() (*NumUnconfirmedTxsInfo, error) {
	b, err := nc.handler.Group("/net_info").Get()
	if err != nil {
		return nil, err
	}
	var bb []byte
	bb, err = nc.preloadJsonResponse(b)
	if err != nil {
		return nil, err
	}
	var a NumUnconfirmedTxsInfo
	err = json.Unmarshal(bb, &a)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (nc *NSBClient) GetStatus() (*StatusInfo, error) {
	b, err := nc.handler.Group("/status").Get()
	if err != nil {
		return nil, err
	}
	var bb []byte
	bb, err = nc.preloadJsonResponse(b)
	if err != nil {
		return nil, err
	}
	var a StatusInfo
	err = json.Unmarshal(bb, &a)
	if err != nil {
		return nil, err
	}
	return &a, nil
}
func (nc *NSBClient) GetUnconfirmedTxs(limit int64) (*NumUnconfirmedTxsInfo, error) {
	b, err := nc.handler.Group("/unconfirmed_txs").GetWithParams(request.Param{
		"limit": limit,
	})
	if err != nil {
		return nil, err
	}
	var bb []byte
	bb, err = nc.preloadJsonResponse(b)
	if err != nil {
		return nil, err
	}
	var a NumUnconfirmedTxsInfo
	err = json.Unmarshal(bb, &a)
	if err != nil {
		return nil, err
	}
	return &a, nil
}
func (nc *NSBClient) GetValidators(id int64) (*ValidatorsInfo, error) {
	b, err := nc.handler.Group("/validators").GetWithParams(request.Param{
		"height": id,
	})
	if err != nil {
		return nil, err
	}
	var bb []byte
	bb, err = nc.preloadJsonResponse(b)
	if err != nil {
		return nil, err
	}
	var a ValidatorsInfo
	err = json.Unmarshal(bb, &a)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

type type_sig = uint64
type Ed25519SignableAccount interface {
	GetPublicKey() []byte
	Sign([]byte) []byte
}

func (nc *NSBClient) sendContractTx(
	transType, contractName []byte,
	txContent *cmn.TransactionHeader,
) (*ResultInfo, error) {
	var b = make([]byte, 0, 65535)
	var buf = bytes.NewBuffer(b)
	buf.Write(transType)
	buf.WriteByte(0x19)
	buf.Write(contractName)
	buf.WriteByte(0x18)
	c, err := json.Marshal(txContent)
	if err != nil {
		return nil, err
	}
	buf.Write(c)
	// fmt.Println(string(c))
	json.Unmarshal(c, txContent)

	return nc.BroadcastTxCommit(buf.Bytes())
}

func (nc *NSBClient) CreateISC(
	user Ed25519SignableAccount,
	funds []uint32, iscOwners [][]byte,
	bytesTransactionIntents [][]byte,
	vesSig []byte,
) ([]byte, error) {
	var txHeader cmn.TransactionHeader
	var buf = bytes.NewBuffer(make([]byte, 65535))
	buf.Reset()
	// fmt.Println(string(buf.Bytes()))
	var transactionIntents []*tx.TransactionIntent
	var txm map[string]interface{}
	for idx, txb := range bytesTransactionIntents {
		err := json.Unmarshal(txb, &txm)
		if err != nil {
			return nil, err
		}
		var txi = new(tx.TransactionIntent)
		if txm["src"] == nil && txm["from"] == nil {
			return nil, errors.New("nil src")
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
		txi.Seq = nmath.NewUint256FromBigInt(big.NewInt(int64(idx)))
		if txm["amt"] != nil {
			b, _ := hex.DecodeString(txm["amt"].(string))
			txi.Amt = nmath.NewUint256FromBytes(b)
		} else {
			txi.Amt = nmath.NewUint256FromBytes([]byte{0})
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
	var mrand = mt19937.New()
	mrand.Seed(time.Now().UnixNano())
	var n1, n2, n3, n4 = mrand.Uint64(), mrand.Uint64(), mrand.Uint64(), mrand.Uint64()

	txHeader.Nonce = nmath.NewUint256FromBytes([]byte{
		uint8(n1 >> 56), uint8(n1>>48) & 0xff, uint8(n1>>40) & 0xff, uint8(n1>>32) & 0xff,
		uint8(n1>>24) & 0xff, uint8(n1>>16) & 0xff, uint8(n1>>8) & 0xff, uint8(n1>>0) & 0xff,
		uint8(n2 >> 56), uint8(n2>>48) & 0xff, uint8(n2>>40) & 0xff, uint8(n2>>32) & 0xff,
		uint8(n2>>24) & 0xff, uint8(n2>>16) & 0xff, uint8(n2>>8) & 0xff, uint8(n2>>0) & 0xff,
		uint8(n3 >> 56), uint8(n3>>48) & 0xff, uint8(n3>>40) & 0xff, uint8(n3>>32) & 0xff,
		uint8(n3>>24) & 0xff, uint8(n3>>16) & 0xff, uint8(n3>>8) & 0xff, uint8(n3>>0) & 0xff,
		uint8(n4 >> 56), uint8(n4>>48) & 0xff, uint8(n4>>40) & 0xff, uint8(n4>>32) & 0xff,
		uint8(n4>>24) & 0xff, uint8(n4>>16) & 0xff, uint8(n4>>8) & 0xff, uint8(n4>>0) & 0xff,
	})
	txHeader.Value = nmath.NewUint256FromBytes([]byte{0})
	// bug: buf.Reset()
	buf = bytes.NewBuffer(make([]byte, 65535))

	buf.Write(txHeader.From)
	buf.Write(txHeader.ContractAddress)
	buf.Write(txHeader.Data)
	buf.Write(txHeader.Value.Bytes())
	buf.Write(txHeader.Nonce.Bytes())
	txHeader.Signature = user.Sign(buf.Bytes())
	ret, err := nc.sendContractTx([]byte("createContract"), []byte("isc"), &txHeader)
	fmt.Println(PretiJson(ret), err)
	if err != nil {
		return nil, err
	}
	return ret.DeliverTx.Data, nil
}

func (nc *NSBClient) createISC(
	w io.Writer,
	funds []uint32, iscOwners [][]byte,
	transactionIntents []*tx.TransactionIntent,
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

	// fmt.Println(PretiJson(args), b)
	_, err = w.Write(b)
	return err
}

func (nc *NSBClient) UpdateTxInfo(
	user Ed25519SignableAccount, contractAddress []byte,
	tid uint64, transactionIntent *tx.TransactionIntent,
) ([]byte, error) {
	var txHeader cmn.TransactionHeader
	var buf = bytes.NewBuffer(make([]byte, 65535))
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
	var mrand = mt19937.New()
	mrand.Seed(time.Now().UnixNano())
	var n1, n2, n3, n4 = mrand.Uint64(), mrand.Uint64(), mrand.Uint64(), mrand.Uint64()

	txHeader.Nonce = nmath.NewUint256FromBytes([]byte{
		uint8(n1 >> 56), uint8(n1>>48) & 0xff, uint8(n1>>40) & 0xff, uint8(n1>>32) & 0xff,
		uint8(n1>>24) & 0xff, uint8(n1>>16) & 0xff, uint8(n1>>8) & 0xff, uint8(n1>>0) & 0xff,
		uint8(n2 >> 56), uint8(n2>>48) & 0xff, uint8(n2>>40) & 0xff, uint8(n2>>32) & 0xff,
		uint8(n2>>24) & 0xff, uint8(n2>>16) & 0xff, uint8(n2>>8) & 0xff, uint8(n2>>0) & 0xff,
		uint8(n3 >> 56), uint8(n3>>48) & 0xff, uint8(n3>>40) & 0xff, uint8(n3>>32) & 0xff,
		uint8(n3>>24) & 0xff, uint8(n3>>16) & 0xff, uint8(n3>>8) & 0xff, uint8(n3>>0) & 0xff,
		uint8(n4 >> 56), uint8(n4>>48) & 0xff, uint8(n4>>40) & 0xff, uint8(n4>>32) & 0xff,
		uint8(n4>>24) & 0xff, uint8(n4>>16) & 0xff, uint8(n4>>8) & 0xff, uint8(n4>>0) & 0xff,
	})
	txHeader.Value = nmath.NewUint256FromBytes([]byte{0})
	// bug: buf.Reset()
	buf = bytes.NewBuffer(make([]byte, 65535))

	buf.Write(txHeader.From)
	buf.Write(txHeader.ContractAddress)
	buf.Write(txHeader.Data)
	buf.Write(txHeader.Value.Bytes())
	buf.Write(txHeader.Nonce.Bytes())
	txHeader.Signature = user.Sign(buf.Bytes())
	_, err = nc.sendContractTx([]byte("sendTransaction"), []byte("isc"), &txHeader)
	// fmt.Println(PretiJson(ret), err)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (nc *NSBClient) updateTxInfo(
	w io.Writer,
	tid uint64, transactionIntent *tx.TransactionIntent,
) error {
	var args ISC.ArgsUpdateTxInfo
	args.Tid = tid
	args.TransactionIntent = transactionIntent
	b, err := json.Marshal(args)
	if err != nil {
		return err
	}

	// fmt.Println(PretiJson(args), b)
	_, err = w.Write(b)
	return err
}

func (nc *NSBClient) FreezeInfo(
	user Ed25519SignableAccount, contractAddress []byte,
	tid uint64,
) ([]byte, error) {
	var txHeader cmn.TransactionHeader
	var buf = bytes.NewBuffer(make([]byte, 65535))
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
	var mrand = mt19937.New()
	mrand.Seed(time.Now().UnixNano())
	var n1, n2, n3, n4 = mrand.Uint64(), mrand.Uint64(), mrand.Uint64(), mrand.Uint64()

	txHeader.Nonce = nmath.NewUint256FromBytes([]byte{
		uint8(n1 >> 56), uint8(n1>>48) & 0xff, uint8(n1>>40) & 0xff, uint8(n1>>32) & 0xff,
		uint8(n1>>24) & 0xff, uint8(n1>>16) & 0xff, uint8(n1>>8) & 0xff, uint8(n1>>0) & 0xff,
		uint8(n2 >> 56), uint8(n2>>48) & 0xff, uint8(n2>>40) & 0xff, uint8(n2>>32) & 0xff,
		uint8(n2>>24) & 0xff, uint8(n2>>16) & 0xff, uint8(n2>>8) & 0xff, uint8(n2>>0) & 0xff,
		uint8(n3 >> 56), uint8(n3>>48) & 0xff, uint8(n3>>40) & 0xff, uint8(n3>>32) & 0xff,
		uint8(n3>>24) & 0xff, uint8(n3>>16) & 0xff, uint8(n3>>8) & 0xff, uint8(n3>>0) & 0xff,
		uint8(n4 >> 56), uint8(n4>>48) & 0xff, uint8(n4>>40) & 0xff, uint8(n4>>32) & 0xff,
		uint8(n4>>24) & 0xff, uint8(n4>>16) & 0xff, uint8(n4>>8) & 0xff, uint8(n4>>0) & 0xff,
	})
	txHeader.Value = nmath.NewUint256FromBytes([]byte{0})
	// bug: buf.Reset()
	buf = bytes.NewBuffer(make([]byte, 65535))

	buf.Write(txHeader.From)
	buf.Write(txHeader.ContractAddress)
	buf.Write(txHeader.Data)
	buf.Write(txHeader.Value.Bytes())
	buf.Write(txHeader.Nonce.Bytes())
	txHeader.Signature = user.Sign(buf.Bytes())
	_, err = nc.sendContractTx([]byte("sendTransaction"), []byte("isc"), &txHeader)
	// fmt.Println(PretiJson(ret), err)
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

	// fmt.Println(PretiJson(args), b)
	_, err = w.Write(b)
	return err
}

func (nc *NSBClient) UserAck(
	user Ed25519SignableAccount, contractAddress []byte,
	address, signature []byte,
) (*DeliverTx, error) {
	var txHeader cmn.TransactionHeader
	var buf = bytes.NewBuffer(make([]byte, 65535))
	buf.Reset()
	// fmt.Println(string(buf.Bytes()))
	err := nc.userAck(buf, address, signature)
	if err != nil {
		return nil, err
	}

	var fap appl.FAPair
	fap.FuncName = "UserAck"
	fap.Args = buf.Bytes()
	txHeader.Data, err = json.Marshal(fap)
	if err != nil {
		return nil, err
	}
	txHeader.ContractAddress = contractAddress
	txHeader.From = user.GetPublicKey()
	var mrand = mt19937.New()
	mrand.Seed(time.Now().UnixNano())
	var n1, n2, n3, n4 = mrand.Uint64(), mrand.Uint64(), mrand.Uint64(), mrand.Uint64()

	txHeader.Nonce = nmath.NewUint256FromBytes([]byte{
		uint8(n1 >> 56), uint8(n1>>48) & 0xff, uint8(n1>>40) & 0xff, uint8(n1>>32) & 0xff,
		uint8(n1>>24) & 0xff, uint8(n1>>16) & 0xff, uint8(n1>>8) & 0xff, uint8(n1>>0) & 0xff,
		uint8(n2 >> 56), uint8(n2>>48) & 0xff, uint8(n2>>40) & 0xff, uint8(n2>>32) & 0xff,
		uint8(n2>>24) & 0xff, uint8(n2>>16) & 0xff, uint8(n2>>8) & 0xff, uint8(n2>>0) & 0xff,
		uint8(n3 >> 56), uint8(n3>>48) & 0xff, uint8(n3>>40) & 0xff, uint8(n3>>32) & 0xff,
		uint8(n3>>24) & 0xff, uint8(n3>>16) & 0xff, uint8(n3>>8) & 0xff, uint8(n3>>0) & 0xff,
		uint8(n4 >> 56), uint8(n4>>48) & 0xff, uint8(n4>>40) & 0xff, uint8(n4>>32) & 0xff,
		uint8(n4>>24) & 0xff, uint8(n4>>16) & 0xff, uint8(n4>>8) & 0xff, uint8(n4>>0) & 0xff,
	})
	txHeader.Value = nmath.NewUint256FromBytes([]byte{0})
	// bug: buf.Reset()
	buf = bytes.NewBuffer(make([]byte, 65535))

	buf.Write(txHeader.From)
	buf.Write(txHeader.ContractAddress)
	buf.Write(txHeader.Data)
	buf.Write(txHeader.Value.Bytes())
	buf.Write(txHeader.Nonce.Bytes())
	txHeader.Signature = user.Sign(buf.Bytes())
	ret, err := nc.sendContractTx([]byte("sendTransaction"), []byte("isc"), &txHeader)
	// fmt.Println(PretiJson(ret), err)
	if err != nil {
		return nil, err
	}
	return &ret.DeliverTx, nil
}

func (nc *NSBClient) userAck(
	w io.Writer,
	address, signature []byte,
) error {
	var args ISC.ArgsUserAck
	args.Address = address
	args.Signature = signature
	b, err := json.Marshal(args)
	if err != nil {
		return err
	}

	// fmt.Println(PretiJson(args), b)
	_, err = w.Write(b)
	return err
}

func (nc *NSBClient) InsuranceClaim(
	user Ed25519SignableAccount, contractAddress []byte,
	tid, aid uint64,
) (*DeliverTx, error) {
	var txHeader cmn.TransactionHeader
	var buf = bytes.NewBuffer(make([]byte, 65535))
	buf.Reset()
	// fmt.Println(string(buf.Bytes()))
	err := nc.insuranceClaim(buf, tid, aid)
	if err != nil {
		return nil, err
	}

	var fap appl.FAPair
	fap.FuncName = "InsuranceClaim"
	fap.Args = buf.Bytes()
	txHeader.Data, err = json.Marshal(fap)
	if err != nil {
		return nil, err
	}
	txHeader.ContractAddress = contractAddress
	txHeader.From = user.GetPublicKey()
	var mrand = mt19937.New()
	mrand.Seed(time.Now().UnixNano())
	var n1, n2, n3, n4 = mrand.Uint64(), mrand.Uint64(), mrand.Uint64(), mrand.Uint64()

	txHeader.Nonce = nmath.NewUint256FromBytes([]byte{
		uint8(n1 >> 56), uint8(n1>>48) & 0xff, uint8(n1>>40) & 0xff, uint8(n1>>32) & 0xff,
		uint8(n1>>24) & 0xff, uint8(n1>>16) & 0xff, uint8(n1>>8) & 0xff, uint8(n1>>0) & 0xff,
		uint8(n2 >> 56), uint8(n2>>48) & 0xff, uint8(n2>>40) & 0xff, uint8(n2>>32) & 0xff,
		uint8(n2>>24) & 0xff, uint8(n2>>16) & 0xff, uint8(n2>>8) & 0xff, uint8(n2>>0) & 0xff,
		uint8(n3 >> 56), uint8(n3>>48) & 0xff, uint8(n3>>40) & 0xff, uint8(n3>>32) & 0xff,
		uint8(n3>>24) & 0xff, uint8(n3>>16) & 0xff, uint8(n3>>8) & 0xff, uint8(n3>>0) & 0xff,
		uint8(n4 >> 56), uint8(n4>>48) & 0xff, uint8(n4>>40) & 0xff, uint8(n4>>32) & 0xff,
		uint8(n4>>24) & 0xff, uint8(n4>>16) & 0xff, uint8(n4>>8) & 0xff, uint8(n4>>0) & 0xff,
	})
	txHeader.Value = nmath.NewUint256FromBytes([]byte{0})
	// bug: buf.Reset()
	buf = bytes.NewBuffer(make([]byte, 65535))

	buf.Write(txHeader.From)
	buf.Write(txHeader.ContractAddress)
	buf.Write(txHeader.Data)
	buf.Write(txHeader.Value.Bytes())
	buf.Write(txHeader.Nonce.Bytes())
	txHeader.Signature = user.Sign(buf.Bytes())
	ret, err := nc.sendContractTx([]byte("sendTransaction"), []byte("isc"), &txHeader)
	// fmt.Println(PretiJson(ret), err)
	if err != nil {
		return nil, err
	}
	return &ret.DeliverTx, nil
}

func (nc *NSBClient) insuranceClaim(
	w io.Writer,
	tid, aid uint64,
) error {
	// var args ISC.ArgsInsuranceClaim
	// args.Tid = tid
	// args.Aid = aid
	// b, err := json.Marshal(args)
	err := binary.Write(w, binary.BigEndian, tid)
	if err != nil {
		return err
	}
	err = binary.Write(w, binary.BigEndian, aid)

	// fmt.Println(PretiJson(args), b)
	// _, err = w.Write(b)
	return err
}

func (nc *NSBClient) SettleContract(
	user Ed25519SignableAccount, contractAddress []byte,
) ([]byte, error) {
	var txHeader cmn.TransactionHeader
	var err error
	var fap appl.FAPair
	fap.FuncName = "SettleContract"
	txHeader.Data, err = json.Marshal(fap)
	if err != nil {
		return nil, err
	}
	txHeader.ContractAddress = contractAddress
	txHeader.From = user.GetPublicKey()
	var mrand = mt19937.New()
	mrand.Seed(time.Now().UnixNano())
	var n1, n2, n3, n4 = mrand.Uint64(), mrand.Uint64(), mrand.Uint64(), mrand.Uint64()

	txHeader.Nonce = nmath.NewUint256FromBytes([]byte{
		uint8(n1 >> 56), uint8(n1>>48) & 0xff, uint8(n1>>40) & 0xff, uint8(n1>>32) & 0xff,
		uint8(n1>>24) & 0xff, uint8(n1>>16) & 0xff, uint8(n1>>8) & 0xff, uint8(n1>>0) & 0xff,
		uint8(n2 >> 56), uint8(n2>>48) & 0xff, uint8(n2>>40) & 0xff, uint8(n2>>32) & 0xff,
		uint8(n2>>24) & 0xff, uint8(n2>>16) & 0xff, uint8(n2>>8) & 0xff, uint8(n2>>0) & 0xff,
		uint8(n3 >> 56), uint8(n3>>48) & 0xff, uint8(n3>>40) & 0xff, uint8(n3>>32) & 0xff,
		uint8(n3>>24) & 0xff, uint8(n3>>16) & 0xff, uint8(n3>>8) & 0xff, uint8(n3>>0) & 0xff,
		uint8(n4 >> 56), uint8(n4>>48) & 0xff, uint8(n4>>40) & 0xff, uint8(n4>>32) & 0xff,
		uint8(n4>>24) & 0xff, uint8(n4>>16) & 0xff, uint8(n4>>8) & 0xff, uint8(n4>>0) & 0xff,
	})
	txHeader.Value = nmath.NewUint256FromBytes([]byte{0})
	// bug: buf.Reset()
	buf := bytes.NewBuffer(make([]byte, 65535))

	buf.Write(txHeader.From)
	buf.Write(txHeader.ContractAddress)
	buf.Write(txHeader.Data)
	buf.Write(txHeader.Value.Bytes())
	buf.Write(txHeader.Nonce.Bytes())
	txHeader.Signature = user.Sign(buf.Bytes())
	_, err = nc.sendContractTx([]byte("sendTransaction"), []byte("isc"), &txHeader)
	// fmt.Println(PretiJson(ret), err)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
