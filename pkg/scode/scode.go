package scode

func NewToken(id TokenID, value string) *Token {
	return &Token{Identifier: id, Value: value}
}

func NewInstruction(tok ...*Token) *Instruction {
	ins := Instruction{}
	ins.AddToken(tok...)
	return &ins
}

func NewCommand(ct CommandType, ins ...*Instruction) *Command {
	com := Command{Type: ct}
	com.AddInstruction(ins...)
	return &com
}

func NewOperation(otype OperationType, com ...*Command) *Operation {
	op := Operation{Type: otype}
	op.AddCommand(com...)
	return &op
}

func NewOperationTree(op ...*Operation) *OperationTree {
	ot := OperationTree{}
	ot.AddOperation(op...)
	return &ot
}

// OperationTree Logic
type OperationTree []*Operation

func (ot *OperationTree) NewOperation(otype OperationType, com ...*Command) *Operation {
	op := &Operation{Type: otype}
	ot.AddOperation(op)
	op.AddCommand(com...)
	return op
}

func (ot *OperationTree) AddOperation(op ...*Operation) {
	for _, o := range op {
		*ot = append(*ot, o)
	}
}

func (ot *OperationTree) GetScript() string {
	text := ""
	for _, op := range *ot {
		for _, com := range op.Commands {
			for _, ins := range com.Instructions {
				for _, tok := range ins.Tokens {
					text += tok.String() + " "
				}
				text += "\n"
			}
			text += "\n"
		}
		text += "\n"
	}
	return text
}
