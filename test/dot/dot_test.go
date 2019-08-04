package dot

import "testing"

type c struct {
	d int
}

type bbb struct {
	c c
}

type a struct {
	b bbb
}

func BenchmarkNo(b *testing.B) {
	var cc int
	b.N = 1 << 35
	for i := 0; i < b.N; i++ {
		cc = i
	}
	_ = cc
}

func BenchmarkDot(b *testing.B) {
	var cc c
	b.N = 1 << 35
	for i := 0; i < b.N; i++ {
		cc.d = i
	}
}

func BenchmarkDot2(b *testing.B) {
	var bb bbb
	b.N = 1 << 35
	for i := 0; i < b.N; i++ {
		bb.c.d = i
	}
}

func BenchmarkDot3(b *testing.B) {
	var aa a
	b.N = 1 << 35
	for i := 0; i < b.N; i++ {
		aa.b.c.d = i
	}
}
