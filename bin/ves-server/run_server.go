package main

import ves_server "github.com/Myriad-Dreamin/go-ves/ves"
import "log"

const port = ":23351"
const centerAddress = "127.0.0.1:23352"

func main() {
	if err := ves_server.ListenAndServe(port, centerAddress); err != nil {
		log.Fatal(err)
	}
}
