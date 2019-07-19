package main

import (
	"testing"
)

func BenchmarkGetMap1(b *testing.B) {
	var s = make(map[QQ]bool)
	var A = QQ{3, string(make([]byte, 2e9))}
	s[A] = true
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = s[A]
	}
}

func BenchmarkGetMap2(b *testing.B) {
	var s = make(map[QQ]bool)
	var A = QQ{3, string(make([]byte, 2e9))}
	s[A] = true
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		a, _ := s[A]
		_ = a
	}
}
