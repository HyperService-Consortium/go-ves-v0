package nsbcli

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	request "github.com/Myriad-Dreamin/go-ves/net/request"
	"github.com/tidwall/gjson"

	appl "github.com/HyperServiceOne/NSB/application"
	cmn "github.com/HyperServiceOne/NSB/common"
	ISC "github.com/HyperServiceOne/NSB/contract/isc"
	tx "github.com/HyperServiceOne/NSB/contract/isc/transaction"
	nmath "github.com/HyperServiceOne/NSB/math"
	math "github.com/Myriad-Dreamin/go-ves/math"
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
	return "http://" + host
}

type NSBClient struct {
	handler    *request.RequestClient
	bufferPool *BytesPool
}

// todo: test invalid json
func (nc *NSBClient) preloadJsonRPCResponse(responseBuffer io.ReadCloser) ([]byte, error) {

	var b = nc.bufferPool.Get().([]byte)
	defer nc.bufferPool.Put(b)

	_, err := responseBuffer.Read(b)
	if err != nil && err != io.EOF {
		return nil, err
	}
	responseBuffer.Close()

	var jrpcMessage = gjson.ParseBytes(b)
	if jrpcVer := jrpcMessage.Get("jsonrpc"); !jrpcVer.Exists() || jrpcVer.String() != "2.0" {
		return nil, errors.New("reject ret that is not jsonrpc: 2.0")
	}
	if errorResp := jrpcMessage.Get("error"); errorResp.Exists() {
		return nil, fromGJsonResultError(errorResp)
	}
	if resultResp := jrpcMessage.Get("result"); resultResp.Exists() {
		if resultResp.Index > 0 {
			return b[resultResp.Index : resultResp.Index+len(resultResp.Raw)], nil
		} else {
			return nil, nil
		}
	}
	return nil, errors.New("bad format of jsonrpc")
}

func NewNSBClient(host string) *NSBClient {
	return &NSBClient{
		handler:    request.NewRequestClient(decorateHost(host)),
		bufferPool: NewBytesPool(),
	}
}

func (nc *NSBClient) GetAbciInfo() (*AbciInfoResponse, error) {
	b, err := nc.handler.Group("/abci_info").GetWithParams(request.Param{})
	if err != nil {
		return nil, err
	}
	var bb []byte
	bb, err = nc.preloadJsonRPCResponse(b)
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
	bb, err = nc.preloadJsonRPCResponse(b)
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
	bb, err = nc.preloadJsonRPCResponse(b)
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
	bb, err = nc.preloadJsonRPCResponse(b)
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
	bb, err = nc.preloadJsonRPCResponse(b)
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
	bb, err = nc.preloadJsonRPCResponse(b)
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

type type_sig = uint64
type Ed25519SignableAccount interface {
	PublicKey() []byte
	Sign([]byte) []byte
}

func (nc *NSBClient) CreateISC(
	user Ed25519SignableAccount,
	funds []uint32, iscOwners [][]byte,
	transactionIntents []*tx.TransactionIntent,
	vesSig []byte,
) error {
	var txHeader cmn.TransactionHeader
	var fap appl.FAPair
	var b = make([]byte, 65535)
	var buf = bytes.NewBuffer(b)
	err := nc.createISC(buf, funds, iscOwners, transactionIntents, vesSig)
	if err != nil {
		return err
	}
	fap.Args = buf.Bytes()
	txHeader.From = user.PublicKey()
	txHeader.Data, err = json.Marshal(fap)
	if err != nil {
		return err
	}
	var mrand = math.New()
	mrand.Seed(time.Now().UnixNano())
	var n1, n2, n3, n4 = mrand.Uint64(), mrand.Uint64(), mrand.Uint64(), mrand.Uint64()

	txHeader.Nonce = nmath.NewUint256FromBytes([]byte{
		uint8(n1 >> 24), uint8(n1>>16) & 0xff, uint8(n1>>8) & 0xff, uint8(n1>>0) & 0xff,
		uint8(n2 >> 24), uint8(n2>>16) & 0xff, uint8(n2>>8) & 0xff, uint8(n2>>0) & 0xff,
		uint8(n3 >> 24), uint8(n3>>16) & 0xff, uint8(n3>>8) & 0xff, uint8(n3>>0) & 0xff,
		uint8(n4 >> 24), uint8(n4>>16) & 0xff, uint8(n4>>8) & 0xff, uint8(n4>>0) & 0xff,
	})
	buf.Reset()
	err = binary.Write(buf, binary.LittleEndian, &txHeader)
	if err != nil {
		return err
	}
	txHeader.Signature = user.Sign(buf.Bytes())
	return nil
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

	fmt.Println(args, b)
	_, err = w.Write(b)
	return err
}
