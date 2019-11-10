package nsbcli

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	transactiontype "github.com/HyperService-Consortium/NSB/application/transaction-type"
	"github.com/HyperService-Consortium/NSB/grpc/nsbrpc"
	"github.com/golang/protobuf/proto"
	"io"
	"strings"
	"sync/atomic"
	"time"

	"github.com/tidwall/gjson"

	"github.com/HyperService-Consortium/go-ves/lib/net/request"
	jsonrpcclient "github.com/HyperService-Consortium/go-ves/lib/net/rpc-client"

	bytespool "github.com/HyperService-Consortium/go-ves/lib/bytes-pool"
)

var SentBytes, ReceivedBytes uint64

type AsyncOption struct {
	Retry   int
	Timeout time.Duration
}

var defaultOption = &AsyncOption{
	Retry:   5,
	Timeout: 10 * time.Second,
}

func NewAsyncOption() *AsyncOption {
	return &AsyncOption{
		Retry:   5,
		Timeout: 10 * time.Second,
	}
}

const (
	mxBytes = 6000
)

func decorateHost(host string) string {
	if strings.HasPrefix(host, httpPrefix) || strings.HasPrefix(host, httpsPrefix) {
		return host
	}
	return httpPrefix + host
}

// NSBClient provides interface to blockchain nsb
type NSBClient struct {
	handler    *request.RequestClient
	bufferPool *bytespool.BytesPool
}

// todo: test invalid json
func (nc *NSBClient) preloadJSONResponse(bb io.ReadCloser) ([]byte, error) {

	var b = nc.bufferPool.Get()
	defer nc.bufferPool.Put(b)

	n, err := bb.Read(b)
	if err != nil && err != io.EOF {
		return nil, err
	}
	bb.Close()
	atomic.AddUint64(&ReceivedBytes, uint64(n))

	var jm = gjson.ParseBytes(b)
	if s := jm.Get("jsonrpc"); !s.Exists() || s.String() != "2.0" {
		return nil, errNotJSON2d0
	}
	if s := jm.Get("error"); s.Exists() {
		return nil, jsonrpcclient.FromGJsonResultError(s)
	}
	if s := jm.Get("result"); s.Exists() {
		if s.Index > 0 {
			return []byte(s.Raw), nil
		}
	}
	return nil, errBadJSON
}

// NewNSBClient return a pointer of nsb client
func NewNSBClient(host string) *NSBClient {
	return &NSBClient{
		handler:    request.NewRequestClient(decorateHost(host)),
		bufferPool: bytespool.NewBytesPool(maxBytesSize),
	}
}

// GetAbciInfo return the abci information of this rpc service
func (nc *NSBClient) GetAbciInfo() (*AbciInfoResponse, error) {
	b, err := nc.handler.Group("/abci_info").GetWithParams(request.Param{})
	if err != nil {
		return nil, err
	}
	var bb []byte
	bb, err = nc.preloadJSONResponse(b)
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

// GetBlock return the the block's information requested of this blockchain
func (nc *NSBClient) GetBlock(id int64) (*BlockInfo, error) {
	b, err := nc.handler.Group("/block").GetWithParams(request.Param{
		"height": id,
	})
	if err != nil {
		return nil, err
	}
	var bb []byte
	bb, err = nc.preloadJSONResponse(b)
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

// GetBlocks return the the blocks's information requested from L to R of this blockchain
func (nc *NSBClient) GetBlocks(rangeL, rangeR int64) (*BlocksInfo, error) {
	b, err := nc.handler.Group("/blockchain").GetWithParams(request.Param{
		"minHeight": rangeL,
		"maxHeight": rangeR,
	})
	if err != nil {
		return nil, err
	}
	var bb []byte
	bb, err = nc.preloadJSONResponse(b)
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

// GetBlockResults return the the blocks's results requested of this blockchain
func (nc *NSBClient) GetBlockResults(id int64) (*BlockResultsInfo, error) {
	b, err := nc.handler.Group("/block_results").GetWithParams(request.Param{
		"height": id,
	})
	if err != nil {
		return nil, err
	}
	var bb []byte
	bb, err = nc.preloadJSONResponse(b)
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

// GetCommitInfo return the the commit information whose blockid is id
func (nc *NSBClient) GetCommitInfo(id int64) (*CommitInfo, error) {
	b, err := nc.handler.Group("/commit").GetWithParams(request.Param{
		"height": id,
	})
	if err != nil {
		return nil, err
	}
	var bb []byte
	bb, err = nc.preloadJSONResponse(b)
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
	bb, err = nc.preloadJSONResponse(b)
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
	atomic.AddUint64(&SentBytes, uint64(len(body)*2))
	b, err := nc.handler.Group("/broadcast_tx_commit").GetWithParams(request.Param{
		"tx": "0x" + hex.EncodeToString(body),
	})
	if err != nil {
		return nil, err
	}
	var bb []byte
	bb, err = nc.preloadJSONResponse(b)
	if err != nil {
		return nil, err
	}
	var a ResultInfo
	err = json.Unmarshal(bb, &a)
	if err != nil {
		return nil, err
	}
	fmt.Println("res", a)
	return &a, nil
}

func (nc *NSBClient) BroadcastTxAsync(body []byte) ([]byte, error) {
	atomic.AddUint64(&SentBytes, uint64(len(body)*2))
	b, err := nc.handler.Group("/broadcast_tx_async").GetWithParams(request.Param{
		"tx": "0x" + hex.EncodeToString(body),
	})
	if err != nil {
		return nil, err
	}
	var bb []byte
	bb, err = nc.preloadJSONResponse(b)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(bb))
	// var a ResultInfo
	// err = json.Unmarshal(bb, &a)
	// if err != nil {
	// 	return nil, err
	// }
	return bb, nil
}

func (nc *NSBClient) BroadcastTxCommitReturnBytes(body []byte) ([]byte, error) {
	b, err := nc.handler.Group("/broadcast_tx_commit").GetWithParams(request.Param{
		"tx": "0x" + hex.EncodeToString(body),
	})
	if err != nil {
		return nil, err
	}
	var bb []byte
	bb, err = nc.preloadJSONResponse(b)
	if err != nil {
		return nil, err
	}
	return bb, nil
}

func (nc *NSBClient) GetConsensusState() (*ConsensusStateInfo, error) {
	b, err := nc.handler.Group("/consensus_state").Get()
	if err != nil {
		return nil, err
	}
	var bb []byte
	bb, err = nc.preloadJSONResponse(b)
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
	bb, err = nc.preloadJSONResponse(b)
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
	bb, err = nc.preloadJSONResponse(b)
	if err != nil {
		return nil, err
		fmt.Println(string(bb))
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
	bb, err = nc.preloadJSONResponse(b)
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

func (nc *NSBClient) GetProof(txHeader []byte, subQuery string) (*ProofResponse, error) {
	b, err := nc.handler.Group("/abci_query").GetWithParams(request.Param{
		//todo: reduce cost of 0x
		"data": "0x" + hex.EncodeToString(txHeader),
		"path": subQuery,
	})
	if err != nil {
		return nil, err
	}
	var bb []byte
	bb, err = nc.preloadJSONResponse(b)
	if err != nil {
		return nil, err
	}
	var a ProofInfo
	err = json.Unmarshal(bb, &a)
	if err != nil {
		return nil, err
	}
	return &a.Response, nil
}

func (nc *NSBClient) GetNumUnconfirmedTxs() (*NumUnconfirmedTxsInfo, error) {
	b, err := nc.handler.Group("/net_info").Get()
	if err != nil {
		return nil, err
	}
	var bb []byte
	bb, err = nc.preloadJSONResponse(b)
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
	bb, err = nc.preloadJSONResponse(b)
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
	bb, err = nc.preloadJSONResponse(b)
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
	bb, err = nc.preloadJSONResponse(b)
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

func (nc *NSBClient) GetTransaction(hash string) ([]byte, error) {
	b, err := nc.handler.Group("/tx").GetWithParams(request.Param{
		"hash": hash,
		//"prove":false,
	})
	if err != nil {
		return nil, err
	}
	var bb []byte
	bb, err = nc.preloadJSONResponse(b)
	if err != nil {
		return nil, err
	}
	// var a NumUnconfirmedTxsInfo
	// err = json.Unmarshal(bb, &a)
	// if err != nil {
	// 	return nil, err
	// }
	return bb, nil
}

// func (nc *NSBClient) sendContractTx(
// 	transType, contractName []byte,
// 	txContent *cmn.TransactionHeader,
// ) (*ResultInfo, error) {
// 	var b = make([]byte, 0, mxBytes)
// 	var buf = bytes.NewBuffer(b)
// 	buf.Write(transType)
// 	buf.WriteByte(0x19)
// 	buf.Write(contractName)
// 	buf.WriteByte(0x18)
// 	c, err := json.Marshal(txContent)
// 	if err != nil {
// 		return nil, err
// 	}
// 	buf.Write(c)
// 	// fmt.Println(string(c))
// 	json.Unmarshal(c, txContent)

// 	return nc.BroadcastTxCommit(buf.Bytes())
// }

func (nc *NSBClient) Serialize(transType transactiontype.Type, txContent *nsbrpc.TransactionHeader) ([]byte, error) {
	x, err := proto.Marshal(txContent)
	if err != nil {
		return nil, err
	}
	var b = make([]byte, 0, len(x)+1)
	var buf = bytes.NewBuffer(b)
	buf.WriteByte(transType)
	buf.Write(x)
	// todo
	return buf.Bytes(), nil
}

func (nc *NSBClient) sendContractTx(
	transType transactiontype.Type, txContent *nsbrpc.TransactionHeader,
) (*ResultInfo, error) {
	b, err := nc.Serialize(transType, txContent)
	if err != nil {
		return nil, err
	}
	return nc.BroadcastTxCommit(b)
}

func (nc *NSBClient) sendContractTxAsync(
	transType uint8,
	txContent *nsbrpc.TransactionHeader,
	option *AsyncOption,
) ([]byte, error) {
	b, err := nc.Serialize(transType, txContent)
	if err != nil {
		return nil, err
	}
	bb, err := nc.BroadcastTxAsync(b)
	if err != nil {
		return nil, err
	}
	b = nil

	var receipt TransactionReceipt
	err = json.Unmarshal(bb, &receipt)
	if err != nil {
		return nil, err
	}

	if receipt.Code != 0 {
		return nil, fmt.Errorf("errorcode: %d, data:%s, log:%s", receipt.Code, receipt.Data, receipt.Log)
	}

	if option == nil {
		option = NewAsyncOption()
	}

	var hash = "0x" + receipt.Hash
	for t := option.Retry; t != 0; t-- {
		bb, err = nc.GetTransaction(hash)
		fmt.Println(string(bb), "...")
		if err != nil {
			time.Sleep(time.Second)
			continue
		}

		return bb, nil
	}

	return bb, nil
}
