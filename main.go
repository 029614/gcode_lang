package main

import (
	"github.com/029614/gcode_lang/lib/data"
)

func main() {
	rtrlib := data.GetRouterLibrary()
	rtr, err := rtrlib.GetRouterByName("Multicam")
	if err != nil {
		panic(err)
	} else {
		println("fuck yeah: ", rtr.ID)
	}
	// Example demonstrating declarative use
	// println(scode.DrillingExample().GetScript())

	// Example demonstrating dynamic use with loops
	// pts := make([][2]float32, 0)
	// pts = append(pts, [2]float32{0, 0})
	// pts = append(pts, [2]float32{10, 0})
	// pts = append(pts, [2]float32{10, 10})
	// pts = append(pts, [2]float32{0, 10})
	// pts = append(pts, [2]float32{0, 0})
	// ex := scode.CuttingExample(pts...)

	//println(ex.GetScript())

	// proc := processor.MulticamProcessor{}
	// proc.PostProcess(ex)

	// println(ex.GetScript())
}
