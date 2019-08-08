package main

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"os"
)

func main() {
	switch os.Args[1] {
	case "base64":
		s := []byte(os.Args[2])
		dst := make([]byte, 65535)
		base64.StdEncoding.Decode(dst, s)
		fmt.Println(string(dst[:base64.StdEncoding.DecodedLen(len(s))]))
		fmt.Println(hex.EncodeToString(dst[:base64.StdEncoding.DecodedLen(len(s))]))
	case "stdstring":
		fmt.Println(hex.EncodeToString([]byte(os.Args[2])))
	case "transaction":
		s := []byte(os.Args[2])
		dst := make([]byte, 65535)
		base64.StdEncoding.Decode(dst, s)
		b := dst[:base64.StdEncoding.DecodedLen(len(s))]
		t := bytes.Split(b, []byte{0x18})[1]
		fmt.Println(string(t), hex.EncodeToString(t))
	}
}
