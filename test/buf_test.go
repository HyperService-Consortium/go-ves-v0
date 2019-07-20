package b

import (
	"fmt"
	"testing"
)

func TestBuf(t *testing.T) {
	s := make([]byte, 0)
	fmt.Println(cap(s))
	s = append(s, 2)
	fmt.Println(cap(s))
	s = append(s, []byte{0, 0, 0, 0, 0, 0, 0, 0}...)
	fmt.Println(cap(s))
	s = append(s, []byte{0, 0, 0, 0, 0, 0, 0, 0}...)
	fmt.Println(cap(s))
}
