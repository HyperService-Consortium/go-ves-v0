package session

import (
	"fmt"
	"testing"
)

func TestComparator(t *testing.T) {
	c := makeComparator()
	fmt.Println(c.Insert(33, []byte{1}))
	fmt.Println(c.Insert(32, []byte{1}))
	fmt.Println(c.Insert(33, []byte{2}))
	fmt.Println(c.Insert(32, []byte{2}))
	fmt.Println(c.Insert(33, []byte{1}))
	fmt.Println(c.Insert(32, []byte{1}))
	fmt.Println(c.Insert(33, []byte{2}))
	fmt.Println(c.Insert(32, []byte{2}))
}
