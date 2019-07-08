package nsbcli

import "time"

/******************************* abci_info ************************************/

type AbciInfo struct {
	Response *AbciInfoResponse `json:"response"`
}

type AbciInfoResponse struct {
	Data       string `json:"data"`
	Version    string `json:"version"`
	AppVersion string `json:"app_version"`
}

/******************************* block_info ***********************************/

type BlockInfo struct {
	BlockMeta *BlockMeta `json:"block_meta"`
	Block     *Block     `json:"block"`
}

type BlockMeta struct {
	BlockID *BlockID `json:"block_id"`
	Header  Header   `json:"header"`
}

type Block struct {
	Header     *Header     `json:"header"`
	Data       *Data       `json:"data"`
	Evidence   *Evidence   `json:"evidence"`
	LastCommit *LastCommit `json:"last_commit"`
}

type BlockID struct {
	Hash  string `json:"hash"`
	Parts *Parts `json:"parts"`
}

type Parts struct {
	Total string `json:"total"`
	Hash  string `json:"hash"`
}

type Header struct {
	Version            *Version     `json:"version"`
	ChainID            string       `json:"chain_id"`
	Height             string       `json:"height"`
	Time               time.Time    `json:"time"`
	NumTxs             string       `json:"num_txs"`
	TotalTxs           string       `json:"total_txs"`
	LastBlockID        *LastBlockID `json:"last_block_id"`
	LastCommitHash     string       `json:"last_commit_hash"`
	DataHash           string       `json:"data_hash"`
	ValidatorsHash     string       `json:"validators_hash"`
	NextValidatorsHash string       `json:"next_validators_hash"`
	ConsensusHash      string       `json:"consensus_hash"`
	AppHash            string       `json:"app_hash"`
	LastResultsHash    string       `json:"last_results_hash"`
	EvidenceHash       string       `json:"evidence_hash"`
	ProposerAddress    string       `json:"proposer_address"`
}

type Version struct {
	Block string `json:"block"`
	App   string `json:"app"`
}

type LastBlockID struct {
	Hash  string `json:"hash"`
	Parts *Parts `json:"parts"`
}

type Data struct {
	Txs []string `json:"txs"`
}

// temporarily unknown
type Evidence struct {
	MaxAge string `json:"max_age"`
}

type LastCommit struct {
	BlockID    *BlockID      `json:"block_id"`
	Precommits []*Precommits `json:"precommits"`
}

type Precommits struct {
	Type             int       `json:"type"`
	Height           string    `json:"height"`
	Round            string    `json:"round"`
	BlockID          *BlockID  `json:"block_id"`
	Timestamp        time.Time `json:"timestamp"`
	ValidatorAddress string    `json:"validator_address"`
	ValidatorIndex   string    `json:"validator_index"`
	Signature        string    `json:"signature"`
}

/**************************** block_results_info ******************************/

type BlockResultsInfo struct {
	Height  string   `json:"height"`
	Results *Results `json:"results"`
}

type Results struct {
	DeliverTxInfo  []*DeliverTxInfo `json:"DeliverTx"`
	EndBlockInfo   *EndBlockInfo    `json:"EndBlock"`
	BeginBlockInfo *BeginBlockInfo  `json:"BeginBlock"`
}

type DeliverTxInfo struct {
	Info string  `json:"info"`
	Tags []*Tags `json:"tags"`
}

type Tags struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// temporarily unknown
type EndBlockInfo struct {
	ValidatorUpdates interface{} `json:"validator_updates"`
}

type BeginBlockInfo struct {
}

/****************************** blocks_info ***********************************/

type BlocksInfo struct {
	LastHeight string       `json:"last_height"`
	BlockMetas []*BlockMeta `json:"block_metas"`
}

/****************************** commit_info ***********************************/

type CommitInfo struct {
	SignedHeader *SignedHeader `json:"signed_header"`
	Canonical    bool          `json:"canonical"`
}

type SignedHeader struct {
	Header *Header `json:"header"`
	Commit *Commit `json:"commit"`
}

type Commit struct {
	BlockID    *BlockID      `json:"block_id"`
	Precommits []*Precommits `json:"precommits"`
}

/************************* consensus_params_info ******************************/

type ConsensusParamsInfo struct {
	BlockHeight     string           `json:"block_height"`
	ConsensusParams *ConsensusParams `json:"consensus_params"`
}

type ConsensusParams struct {
	BlockConsensus *BlockConsensus `json:"block"`
	Evidence       *Evidence       `json:"evidence"`
	Validator      *Validator      `json:"validator"`
}

type BlockConsensus struct {
	MaxBytes   string `json:"max_bytes"`
	MaxGas     string `json:"max_gas"`
	TimeIotaMs string `json:"time_iota_ms"`
}

type Validator struct {
	PubKeyTypes []string `json:"pub_key_types"`
}
