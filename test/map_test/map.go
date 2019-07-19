package main

import "fmt"

type QQ struct {
	A uint64
	B string
}

type GG struct {
	A uint64
	B []byte
}

func (q *QQ) GetB() []byte {
	return []byte(q.B)
}

func (g *GG) GetB() []byte {
	return g.B
}

func main() {
	var s = make(map[QQ]bool)
	var a = QQ{A: 1, B: "qwq"}
	s[a] = true
	var b = QQ{A: 1, B: "qwq"}
	_, ok := s[b]
	fmt.Println(ok)
}
