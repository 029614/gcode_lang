package scode

import "errors"

const TOKEN_JOB_START = "JOBSTART" // G-code line that starts the job
const TOKEN_JOB_END = "JOBEND"     // G-code line that ends the job

const TOKEN_SPINDLE = "SPINDLE" // G-code line that starts the spindle
const TOKEN_DRILL = "DRILL"     // G-code line that drills a hole

const TOKEN_SLEW = "SLEW"           // G-code line that moves the machine without cutting
const TOKEN_MACHINE = "MACHINE"     // G-code line that sets the machine parameters
const TOKEN_ARC_CW_2D = "ARC2DCW"   // G-code line that cuts an arc in the clockwise direction
const TOKEN_ARC_CCW_2D = "ARC2DCCW" // G-code line that cuts an arc in the counter-clockwise direction

const TOKEN_COMMENT = ";" // G-code line that is a comment

const TOKEN_PARAMETER_X = "X"     // G-code line that sets the X parameter
const TOKEN_PARAMETER_Y = "Y"     // G-code line that sets the Y parameter
const TOKEN_PARAMETER_Z = "Z"     // G-code line that sets the Z parameter
const TOKEN_PARAMETER_I = "I"     // G-code line that sets the I (Arc X) parameter
const TOKEN_PARAMETER_J = "J"     // G-code line that sets the J (Arc Y) parameter
const TOKEN_PARAMETER_K = "K"     // G-code line that sets the K (Arc Z) parameter
const TOKEN_PARAMETER_TOOL = "T"  // G-code line that sets the tool parameter
const TOKEN_PARAMETER_SPEED = "S" // G-code line that sets the RPM parameter
const TOKEN_PARAMETER_FEED = "F"  // G-code line that sets the feed rate parameter

type Token string
type Instruction []Token
type Command []Instruction
type Operation []Command
type OperationTree []Operation
type Coder struct {
	Operations *OperationTree
}

func (t Token) NewToken(value string) (Token, error) {
	err := IsValidToken(value)
	if err != nil {
		return Token(""), err
	} else {
		return Token(value), nil
	}
}

func (i Instruction) NewInstruction(tokens ...Token) Instruction {
	return Instruction(tokens)
}

func (c Command) NewCommand(instructions ...Instruction) Command {
	return Command(instructions)
}

func (o Operation) NewOperation(commands ...Command) Operation {
	return Operation(commands)
}

func (ot OperationTree) NewOperationTree(operations ...Operation) OperationTree {
	return OperationTree(operations)
}

func IsValidToken(value string) error {
	return errors.New("invalid token")
}

func (t Token) IsValidInstruction(tok Token) error {
	return errors.New("invalid Instruction")
}

func (t Token) IsValidCommand(tok Token) error {
	return errors.New("invalid Command")
}

func (t Token) IsValidOperation(tok Token) error {
	return errors.New("invalid Operation")
}

func (t Token) IsValidComment(tok Token) error {
	return errors.New("invalid Comment")
}

func (c Coder) NewOperation(tok Token) {

}
