package scode

// Instruction Logic
type Instruction struct {
	Tokens []*Token
}

func (ins *Instruction) AddToken(tok ...*Token) {
	ins.Tokens = append(ins.Tokens, tok...)
}
