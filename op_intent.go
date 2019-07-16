package main

import "fmt"

import GoPy "github.com/Myriad-Dreamin/go-py"

import opintent "github.com/Myriad-Dreamin/go-uip/op-intent"
import fileload "github.com/Myriad-Dreamin/go-uip/file-load"

func main() {
	defer GoPy.AtExit()

	x := fileload.LoadJson("./json_op_intents/opintents5.json")

	fmt.Println(x)
	y := opintent.BuildGraph(x)
	z := GoPy.GetItem(y, 0)
	fmt.Printf(opintent.Jsonize(z))

	GoPy.DecRef(&x)
	GoPy.DecRef(&y)
}
