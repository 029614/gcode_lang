package scode

// Command Logic
type Command struct {
	Instructions []*Instruction
	Type         CommandType
}

type CommandType int

const (
	CT_ANY CommandType = iota
	CT_START
	CT_STOP

	CT_SPINDLESET
	CT_SPINDLEMOTION

	CT_DRILLSET
	CT_DRILLMOTION
)

func (com *Command) NewInstruction(tok ...*Token) *Instruction {
	ins := &Instruction{}
	com.Instructions = append(com.Instructions, ins)
	ins.AddToken(tok...)
	return ins
}

func (com *Command) AddInstruction(ins ...*Instruction) {
	com.Instructions = append(com.Instructions, ins...)
}
