package types

type Bytable interface {
	FromBytes() error
	Bytes() []byte
}

type Stringable interface {
	String() string
}
