package main

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"os"
)

func main() {
	b, err := hex.DecodeString(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}
	dst := make([]byte, 20000)
	base64.StdEncoding.Encode(dst, b)
	fmt.Println(string(dst[:base64.StdEncoding.EncodedLen(len(b))]))
}
