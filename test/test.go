package main

import (
	"fmt"
)

type ww struct {}

func (c ww) String() string {
	return "wwwwww"
}

func main() {
	var x ww
	fmt.Println(x)
}
