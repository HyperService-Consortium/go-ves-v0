package main

import (
	"log"

	vesclient "github.com/Myriad-Dreamin/go-ves/net/ves_client"
)

func main() {
	log.SetFlags(log.Lshortfile | log.Ltime | log.Ldate)
	vesclient.Main()
}
