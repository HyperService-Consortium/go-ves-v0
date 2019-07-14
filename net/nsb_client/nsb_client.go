package nsbcli

import (
	"encoding/json"
	"errors"
	"io"
	"strings"

	request "github.com/Myriad-Dreamin/go-ves/net/request"
	"github.com/tidwall/gjson"
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
		return nil, fromGJsonResultError(s)
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
		bufferPool: NewBytesPool(),
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
