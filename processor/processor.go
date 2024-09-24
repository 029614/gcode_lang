package processor

import (
	gcode "github.com/029614/gcode_lang/lib"
)

type Processor interface {
	PostProcessCommand(cmd gcode.Command) error
	PostProcessInstruction(ins gcode.Instruction) error
	PostProcessToken(tok gcode.Token) error
	Export(filepath string) (string, error)
}

type ScodeProcessor struct {
	text string
	Tree gcode.Tree
}

func (p *ScodeProcessor) Export(path string) error {
	return nil
}

func (p *ScodeProcessor) PostProcessCommand(cmd gcode.Command) error {
	if cmd.Instructions[0].Tokens[0].Rune == 'G' && cmd.Instructions[0].Tokens[0].Value == "81" {
		// block drill command
	}
	return nil
}

func (p *ScodeProcessor) PostProcessInstruction(ins gcode.Instruction) error {
	return nil
}

func (p *ScodeProcessor) PostProcessToken(tok gcode.Token) error {
	return nil
}

// TODO: Implement PreProcessor
//func PreProcessor(processor Processor, tokens []gcode.Token) ([]gcode.Token, error) {
//	return "", nil
//}

func PostProcess(processor Processor, tree *gcode.Tree) error {
	for _, command := range tree.Commands {
		// process command

		for _, instruction := range command.Instructions {
			// process instruction

			for _, token := range instruction.Tokens {
				// process token
				processor.PostProcessToken(token)

			}

			processor.PostProcessInstruction(instruction)
		}

		processor.PostProcessCommand(command)
	}
	return nil
}
