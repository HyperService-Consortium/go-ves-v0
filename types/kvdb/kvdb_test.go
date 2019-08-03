package kvdb

import (
	"fmt"
	"testing"
)

func TestClone(t *testing.T) {
	fmt.Println(string(clone([]byte{'a', 'b', 'c'})), len(clone([]byte{'a', 'b', 'c'})))
}

func TestCloneWithLen(t *testing.T) {
	fmt.Println(string(cloneWithLen([]byte{'a', 'b', 'c'}, 7)), len(cloneWithLen([]byte{'a', 'b', 'c'}, 7)))
}

func TestDecorate(t *testing.T) {
	fmt.Println(string(decorate([]byte{'a', 'b', 'c'}, []byte{'o', 'r', 'z'})))
}
