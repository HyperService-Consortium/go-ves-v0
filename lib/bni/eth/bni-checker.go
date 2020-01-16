package bni

func (bn *BN) CheckAddress(addr []byte) bool {
	return len(addr) == 20
}
