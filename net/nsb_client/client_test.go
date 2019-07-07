package nsbcli

import (
	"fmt"
	"testing"

	mtest "github.com/Myriad-Dreamin/mydrest"
)

var s mtest.TestHelper

const host = "47.251.2.73:26657"

func TestGetAbciInfo(t *testing.T) {
	b, err := NewNSBClient(host).GetAbciInfo()
	s.AssertNoErr(t, err)
	fmt.Println(PretiJson(b))
}

func TestGetBlock(t *testing.T) {
	b, err := NewNSBClient(host).GetBlock(1)
	s.AssertNoErr(t, err)
	fmt.Println(b)
}
