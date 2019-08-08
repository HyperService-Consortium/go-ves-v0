package main

import (
	"log"

	vesclient "github.com/Myriad-Dreamin/go-ves/lib/net/ves-client"
)

func main() {
	log.SetFlags(log.Lshortfile | log.Ltime | log.Ldate)
	vesclient.Main()
}
