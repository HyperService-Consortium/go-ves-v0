package EthProcessor

type Processor struct {
}

func (p *Processor) CheckAddress(addr []byte) bool {
	return true
}

var p = new(Processor)

func GetProcessor() *Processor {
	return p
}
