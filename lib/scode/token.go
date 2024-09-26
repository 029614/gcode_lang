package scode

import "fmt"

const ID_JOB_START = TokenID("JOBSTART")  // G-code line that starts the job
const ID_JOB_END = TokenID("JOBEND")      // G-code line that ends the job
const ID_SPINDLE = TokenID("SPINDLE")     // G-code line that starts the spindle
const ID_DRILL = TokenID("DRILL")         // G-code line that drills a hole
const ID_MOVE = TokenID("MOVE")           // G-code line that moves the machine without cutting
const ID_CUT = TokenID("CUT")             // G-code line that sets the machine parameters
const ID_ARC_CW_2D = TokenID("ARC2DCW")   // G-code line that cuts an arc in the clockwise direction
const ID_ARC_CCW_2D = TokenID("ARC2DCCW") // G-code line that cuts an arc in the counter-clockwise direction
const ID_COMMENT = TokenID(";")           // G-code line that is a comment
const ID_LINEBREAK = TokenID("\n")        // G-code line break
const ID_TAB = TokenID("\t")              // G-code tab
const ID_PARAMETER_X = TokenID("X")       // G-code line that sets the X parameter
const ID_PARAMETER_Y = TokenID("Y")       // G-code line that sets the Y parameter
const ID_PARAMETER_Z = TokenID("Z")       // G-code line that sets the Z parameter
const ID_PARAMETER_I = TokenID("I")       // G-code line that sets the I (Arc X) parameter
const ID_PARAMETER_J = TokenID("J")       // G-code line that sets the J (Arc Y) parameter
const ID_PARAMETER_K = TokenID("K")       // G-code line that sets the K (Arc Z) parameter
const ID_PARAMETER_TOOL = TokenID("T")    // G-code line that sets the tool parameter
const ID_PARAMETER_SPEED = TokenID("S")   // G-code line that sets the RPM parameter
const ID_PARAMETER_FEED = TokenID("F")    // G-code line that sets the feed rate parameter

type TokenID string

// Token Logic
type Token struct {
	Identifier TokenID
	Value      string
}

func (t Token) String() string {
	return fmt.Sprintf("%s%s", t.Identifier, t.Value)
}
