package types

type transaction = []byte
type transaction_id = uint64
type transaction_local_id = uint32
type content = []byte
type dependencies = []byte

type OpIntents interface {
	GetContents() []content
	GetDependencies() []dependencies
}

type Attestation interface {
	GetTid() transaction_id
	GetContent() content
	GetSignatures() []Signature
}
