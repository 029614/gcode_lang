package processor

import "github.com/029614/gcode_lang/lib/scode"

type MulticamProcessor struct {
	ProcessorBase
}

func (mp *MulticamProcessor) PostProcessOperation(operation *scode.Operation) {
	// Do something with the operation
	var tok1 = operation.GetFirstCodeToken()
	if tok1 == nil {
		// not a valid operation
		return
	}
	println("checking tok1: ", tok1.String())
	switch tok1.Identifier {

	case scode.ID_DRILL:
		mp.handleDrillOperation(operation)

	case scode.ID_SPINDLE:
		tok1.Identifier = "M03"
	}
}

func (mp *MulticamProcessor) PostProcessCommand(command *scode.Command) {
	switch command.Type {
	case scode.CT_ANY:

	case scode.CT_START:
		ins := make([]*scode.Instruction, 0)
		ins = append(ins, scode.NewInstruction(scode.NewToken(scode.ID_COMMENT, "// Multicam Start")))
		ins = append(ins, scode.NewInstruction(scode.NewToken(scode.TokenID("M"), "90")))
		ins = append(ins, scode.NewInstruction(scode.NewToken(scode.TokenID("G"), "90")))
		ins = append(ins, scode.NewInstruction(scode.NewToken(scode.TokenID("G"), "75")))
		command.Instructions = ins

	case scode.CT_STOP:
		iList := make([]*scode.Instruction, 0)
		iList = append(iList, scode.NewInstruction(scode.NewToken(scode.ID_COMMENT, "// Multicam End")))
		iList = append(iList, scode.NewInstruction(scode.NewToken(scode.TokenID("M"), "12")))
		iList = append(iList, scode.NewInstruction(scode.NewToken(scode.TokenID("M"), "05")))
		iList = append(iList, scode.NewInstruction(scode.NewToken(scode.TokenID("G"), "98"), scode.NewToken(scode.TokenID("P"), "147"), scode.NewToken(scode.TokenID("D"), "1")))
		iList = append(iList, scode.NewInstruction(scode.NewToken(scode.TokenID("M"), "02")))
		command.Instructions = iList

	case scode.CT_SPINDLESET:
		iList := make([]*scode.Instruction, 0)
		iList = append(iList, scode.NewInstruction(scode.NewToken(scode.ID_COMMENT, "// Setting Spindle Parameters")))
		for _, ins := range command.Instructions {
			for _, tok := range ins.Tokens {
				if tok.Identifier == scode.ID_PARAMETER_SPEED {
					iList = append(iList, scode.NewInstruction(
						scode.NewToken(scode.TokenID("G"), "97"),
						tok,
					))
				} else if tok.Identifier == scode.ID_PARAMETER_TOOL {
					iList = append(iList, scode.NewInstruction(
						scode.NewToken(scode.TokenID("G"), "00"),
						tok,
					))
				}
			}
		}
		command.Instructions = iList

	case scode.CT_SPINDLEMOTION:

	case scode.CT_DRILLSET:

	case scode.CT_DRILLMOTION:
	}
}

func (mp *MulticamProcessor) handleDrillOperation(operation *scode.Operation) {
}

func (mp *MulticamProcessor) PostProcessInstruction(instruction *scode.Instruction) {
}

func (mp *MulticamProcessor) PostProcessToken(token *scode.Token) {
	switch token.Identifier {
	case scode.ID_DRILL:
	case scode.ID_SPINDLE:
	case scode.ID_JOB_START:
	case scode.ID_JOB_END:
	case scode.ID_MOVE:
		token.Identifier = "G00"
	case scode.ID_CUT:
		token.Identifier = "G01"
	case scode.ID_ARC_CCW_2D:
		token.Identifier = "G03"
	case scode.ID_ARC_CW_2D:
		token.Identifier = "G02"
	}
}

// dereferencing hell, but it works.
func (pb *MulticamProcessor) PostProcess(ot *scode.OperationTree) {
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
