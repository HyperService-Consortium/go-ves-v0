package adaptor

import "testing"

type adaptee struct {
	hh []byte
}

func (e *adaptee) GetA() []byte {
	return e.hh
}

type adapti interface {
	GetA() []byte
}

type adaptor struct {
	b []byte `xorm:"my_hh"`
}

func BenchmarkPureSet(b *testing.B) {
	b.StopTimer()
	var x = adaptee{hh: make([]byte, 2e9)}

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		var xx adaptor
		xx.b = x.hh
	}
}

func BenchmarkInterfaceSet(b *testing.B) {
	b.StopTimer()
	var x = &adaptee{hh: make([]byte, 2e9)}
	var ix adapti = x

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		var xx adaptor
		xx.b = ix.GetA()
	}
}

func Adaptorize(a adapti) *adaptor {
	return &adaptor{
		b: a.GetA(),
	}
}

func BenchmarkInterfaceFuncSet(b *testing.B) {
	b.StopTimer()
	var x = &adaptee{hh: make([]byte, 2e9)}
	var ix adapti = x

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_ = Adaptorize(ix)
	}
}
