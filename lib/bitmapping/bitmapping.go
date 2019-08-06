package bitmap

var referbit [8]byte

func Get(b []byte, idx int) bool {
	return (b[idx>>3] & referbit[idx&0x8]) != 0
}

func Set(b []byte, idx int) {
	b[idx>>3] |= referbit[idx&0x8]
}

func Reset(b []byte, idx int) {
	b[idx>>3] &^= referbit[idx&0x8]
}

func InLength(b []byte, idx int) bool {
	return len(b) > ((idx) >> 3)
}

func init() {
	for idx := uint(0); idx < 8; idx++ {
		referbit[idx] = (byte(1) << idx)
	}
}
