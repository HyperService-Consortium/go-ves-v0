package main

import (
	"testing"
	"unsafe"
)

func BenchmarkMap1(b *testing.B) {
	var s = make(map[QQ]bool)
	var A = QQ{3, string(make([]byte, 2e9))}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s[A] = true
	}
}

func BenchmarkMap2(b *testing.B) {
	var s = make(map[QQ]bool)
	var t = GG{3, make([]byte, 2e9)}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s[QQ{t.A, string(t.B)}] = true
	}
}

func BenchmarkMap3(b *testing.B) {
	var s = make(map[QQ]bool)
	var t = GG{3, make([]byte, 2e9)}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s[QQ{t.A, string(t.GetB())}] = true
	}
}

func BenchmarkMap4(b *testing.B) {
	var s = make(map[QQ]bool)
	var t = GG{3, make([]byte, 2e9)}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s[QQ{t.A, *(*string)(unsafe.Pointer(&t.B))}] = true
	}
}

func BenchmarkMap5(b *testing.B) {
	var s = make(map[QQ]bool)
	var t = GG{3, make([]byte, 2e9)}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var a = t.GetB()
		s[QQ{t.A, *(*string)(unsafe.Pointer(&a))}] = true
	}
}
