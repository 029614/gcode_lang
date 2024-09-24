package main

import (
	"fmt"

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
				//println("Command: ")
				println(string(command.Instructions[0].Tokens[0].Rune), command.Instructions[0].Tokens[0].Value)
				for _, instruction := range command.Instructions {
					//println("\tInstruction: ")
					//println(string(instruction.Tokens[0].Rune), instruction.Tokens[0].Value)
					//fmt.Println("Parameters:")
					for _, token := range instruction.Tokens {
						// Print each key-value pair in the map with the rune as a character (string) and its value
						fmt.Println(" ", string(token.Rune), token.Value)
						for k, v := range token.State.Parameters {
							fmt.Printf("  %q: %s\n", k, v) // %q prints the rune as a quoted character
						}
					}
				}
			}
		}
	}
}
