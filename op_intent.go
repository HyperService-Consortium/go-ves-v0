package main


import "fmt"

import "github.com/Myriad-Dreamin/go-ves/go-py"

import "github.com/Myriad-Dreamin/go-ves/go-uiputils/op-intent"
import "github.com/Myriad-Dreamin/go-ves/go-uiputils/file-load"


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