package main

import (
	"log"

	ves_server "github.com/Myriad-Dreamin/go-ves/ves"
)

const port = ":23351"
const centerAddress = "127.0.0.1:23352"

func main() {
	if err := ves_server.ListenAndServe(port, centerAddress); err != nil {
		log.Fatal(err)
	}
}
