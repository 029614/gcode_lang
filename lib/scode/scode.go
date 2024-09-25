package scode

import (
	"fmt"
)

type TokenID string

// Token Logic
type Token struct {
	Identifier TokenID
	Value      string
}

func (t Token) String() string {
	return fmt.Sprintf("%s%s", t.Identifier, t.Value)
}

// OperationTree Logic
type OperationTree []*Operation

func (ot *OperationTree) NewOperation(com ...Command) *Operation {
	op := &Operation{}
	*ot = append(*ot, op)
	for _, c := range com {
		*op = append(*op, &c)
	}
	return op
}

func (ot *OperationTree) AddOperation(op ...Operation) {
	for _, o := range op {
		*ot = append(*ot, &o)
	}
}

func (ot *OperationTree) GetScript() string {
	text := ""
	for _, op := range *ot {
		for _, com := range *op {
			for _, ins := range *com {
				for _, tok := range *ins {
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

// Operation Logic
type Operation []*Command

func (op *Operation) NewCommand(ins ...Instruction) *Command {
	com := &Command{}
	*op = append(*op, com)
	for _, i := range ins {
		*com = append(*com, &i)
	}
	return com
}

func (op *Operation) AddCommand(com ...Command) {
	for _, c := range com {
		*op = append(*op, &c)
	}
}

// Command Logic
type Command []*Instruction

func (com *Command) NewInstruction(tok ...Token) *Instruction {
	ins := &Instruction{}
	*com = append(*com, ins)
	ins.AddToken(tok...)
	return ins
}

func (com *Command) AddInstruction(ins ...Instruction) {
	for _, i := range ins {
		*com = append(*com, &i)
	}
}

// Instruction Logic
type Instruction []*Token

func (ins *Instruction) AddToken(tok ...Token) {
	for _, t := range tok {
		*ins = append(*ins, &t)
	}
}
