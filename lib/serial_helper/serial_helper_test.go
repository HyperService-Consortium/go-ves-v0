package SerialHelper

import (
	"bytes"
	"fmt"
	"io"
	"testing"

	types "github.com/HyperService-Consortium/go-uip/uiptypes"

	mtest "github.com/Myriad-Dreamin/mydrest"
)

var s mtest.TestHelper

type mA struct {
	x uint64
	y []byte
}

func (a *mA) GetChainId() uint64 { return a.x }
func (a *mA) GetAddress() []byte { return a.y }

func TestSerializeAccountInterface(t *testing.T) {
	a := &mA{1, []byte{6, 6}}

	qwq, err := SerializeAccountInterface(a)
	s.AssertNoErr(t, err)

	var b = new(mA)
	var n int64
	n, b.x, b.y, err = UnserializeAccountInterface(qwq)
	s.AssertNoErr(t, err)
	fmt.Println(n, len(qwq), a, b)

	s.AssertEqual(t, a.x, b.x)

	s.AssertTrue(t, bytes.Equal(a.y, b.y))
}

func TestInsufficientBytesError(t *testing.T) {
	a := &mA{1, []byte{6, 6}}

	qwq, err := SerializeAccountInterface(a)
	s.AssertNoErr(t, err)

	qwq = qwq[:len(qwq)-1]
	var b = new(mA)
	_, b.x, b.y, err = UnserializeAccountInterface(qwq)
	s.AssertEqual(t, err, errInsufficientBytes)
}

func TestTwoAccounts(t *testing.T) {
	a, b := &mA{1, []byte{6, 6}}, &mA{2, []byte{6, 2, 2, 2, 3, 4, 6}}

	var buf = new(bytes.Buffer)
	qwq, err := SerializeAccountInterface(a)
	s.AssertNoErr(t, err)
	_, err = buf.Write(qwq)
	s.AssertNoErr(t, err)
	qwq, err = SerializeAccountInterface(b)
	s.AssertNoErr(t, err)
	_, err = buf.Write(qwq)
	s.AssertNoErr(t, err)
	qwq = buf.Bytes()

	var c = new(mA)
	var n int64
	n, c.x, c.y, err = UnserializeAccountInterface(qwq)
	s.AssertNoErr(t, err)

	fmt.Println(n, len(qwq), a, c)
	s.AssertEqual(t, a.x, c.x)
	s.AssertTrue(t, bytes.Equal(a.y, c.y))

	qwq = qwq[n:]

	n, c.x, c.y, err = UnserializeAccountInterface(qwq)
	s.AssertNoErr(t, err)

	fmt.Println(n, len(qwq), b, c)
	s.AssertEqual(t, b.x, c.x)
	s.AssertTrue(t, bytes.Equal(b.y, c.y))

	qwq = qwq[n:]

	n, c.x, c.y, err = UnserializeAccountInterface(qwq)
	fmt.Println(n, len(qwq), c)
	s.AssertEqual(t, err, io.EOF)
}

func TestThreeAccounts(t *testing.T) {
	a := []types.Account{
		&mA{1, []byte{6, 6}},
		&mA{2, []byte{6, 2, 2, 2, 3, 4, 6}},
		&mA{1, []byte{6, 6}},
		&mA{0xaa00aa00, []byte{0xa0, 0x0a, 0x5a, 0xa5, 0x50, 0x05}},
		&mA{0xaa0000aa, []byte{0xa0, 0x0a, 0x5a, 0xa5, 0x50, 0x05}},
	}

	qwq, err := SerializeAccountsInterface(a)
	s.AssertNoErr(t, err)

	var c = new(mA)
	var n int64
	for _, aa := range a {
		n, c.x, c.y, err = UnserializeAccountInterface(qwq)
		s.AssertNoErr(t, err)

		fmt.Println(n, len(qwq), aa, c)
		s.AssertEqual(t, aa.GetChainId(), c.x)
		s.AssertTrue(t, bytes.Equal(aa.GetAddress(), c.y))
		qwq = qwq[n:]
	}

	n, c.x, c.y, err = UnserializeAccountInterface(qwq)
	fmt.Println(n, len(qwq), c)
	s.AssertEqual(t, err, io.EOF)
}

func TestSerializeAttestationContent(t *testing.T) {

	qwq, err := SerializeAttestationContent(233, 1, []byte("raw transaction"))
	s.AssertNoErr(t, err)
	var a uint64
	var b uint8
	var c []byte
	a, b, c, err = UnserializeAttestationContent(qwq)
	s.AssertNoErr(t, err)
	fmt.Println(a, b, string(c))

	s.AssertEqual(t, a, uint64(233))
	s.AssertEqual(t, b, uint8(1))
	s.AssertTrue(t, bytes.Equal(c, []byte("raw transaction")))
}
