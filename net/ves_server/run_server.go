package main

import ves_server "github.com/Myriad-Dreamin/go-ves/ves"
import "log"

const port = ":23351"

func main() {
	if err := ves_server.ListenAndServe(port); err != nil {
		log.Fatal(err)
	}
}
