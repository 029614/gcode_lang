package processor

import "github.com/029614/gcode_lang/pkg/scode"

type ProcessorBase struct{}

func (pb *ProcessorBase) PostProcessOperation(operation *scode.Operation)       {}
func (pb *ProcessorBase) PostProcessCommand(command *scode.Command)             {}
func (pb *ProcessorBase) PostProcessInstruction(instruction *scode.Instruction) {}
func (pb *ProcessorBase) PostProcessToken(token *scode.Token)                   {}
func (pb *ProcessorBase) PostProcessHandleDrill(operation *scode.Operation)     {}

func (pb *ProcessorBase) PostProcess(ot *scode.OperationTree) {
	for _, op := range *ot {
		for _, com := range op.Commands {
			for _, ins := range com.Instructions {
				for _, tok := range ins.Tokens {
					pb.PostProcessToken(tok)
				}
				pb.PostProcessInstruction(ins)
			}
			pb.PostProcessCommand(com)
		}
		pb.PostProcessOperation(op)
	}
}
