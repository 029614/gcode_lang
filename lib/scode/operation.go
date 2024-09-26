package scode

// Operation Logic
type Operation struct {
	Commands []*Command
	Type     OperationType
}

type OperationType int

const (
	OT_ANY OperationType = iota
	OT_START
	OT_END
	OT_DRILL
	OT_SPINDLE
)

func (op *Operation) NewCommand(ct CommandType, ins ...*Instruction) *Command {
	com := &Command{Instructions: []*Instruction{}, Type: ct}
	op.Commands = append(op.Commands, com)
	return com
}

func (op *Operation) AddCommand(com ...*Command) {
	op.Commands = append(op.Commands, com...)
}

func (op *Operation) GetFirstCodeToken() *Token {
	for _, tok := range op.GetTokenList() {
		valid := (tok.Identifier == ID_JOB_START ||
			tok.Identifier == ID_JOB_END ||
			tok.Identifier == ID_DRILL ||
			tok.Identifier == ID_MOVE ||
			tok.Identifier == ID_CUT ||
			tok.Identifier == ID_ARC_CCW_2D ||
			tok.Identifier == ID_ARC_CW_2D ||
			tok.Identifier == ID_SPINDLE)

		if valid {
			return tok
		}
	}
	return nil
}

func (op *Operation) GetTokenList() []*Token {
	tokens := []*Token{}
	for _, com := range op.Commands {
		for _, ins := range com.Instructions {
			tokens = append(ins.Tokens, ins.Tokens...)
		}
	}
	return tokens
}
