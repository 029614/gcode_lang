package main

import (
	gcode "github.com/029614/gcode_lang/lib"
)

func main() {
	println("Hello World!")

	tree := gcode.NewTree("./tests/GCODE_ROSETTASTONE/nextech/__7 PrefinBirchPly_008.anc")
	err := tree.Parse()
	if err != nil {
		println(err.Error())
	} else {
		println("Parsed successfully")

		println(len(tree.Commands))

		for i, command := range tree.Commands {
			if i < 25 {
				println("Command: ")
				for _, instruction := range command.Instructions {
					println("\tInstruction: ")
					for _, token := range instruction.Tokens {
						println("\t\tToken: [", string(token.Rune), "] ", token.Value)
					}
				}
			}
		}
	}
}
