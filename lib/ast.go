package gcode

// Tool-specific and motion-related tokens
const TOKEN_TOOL_CHANGE_OPERATOR_MESSAGE = "C" // Tool Change Operator Message
const TOKEN_PECK_DRILL_DATA = "D"              // Peck Drill Data
const TOKEN_FEEDRATE = "F"                     // Feedrate in Units per Second

// G-code commands
const TOKEN_PREPARATORY_FUNCTION = "G" // Preparatory Function

// Circular Interpolation Values
const TOKEN_CIRCULAR_INTERPOLATION_X = "I" // Circular Interpolation X (G02, G03)
const TOKEN_CIRCULAR_INTERPOLATION_Y = "J" // Circular Interpolation Y (G02, G03)
const TOKEN_CIRCULAR_INTERPOLATION_Z = "K" // Circular Interpolation Z (G02, G03)

// Miscellaneous or control function
const TOKEN_MISC_CONTROL = "M" // Miscellaneous Control Function

// Sequence number
const TOKEN_SEQUENCE_NUMBER = "N" // Sequence Number

// Motion parameters
const TOKEN_Z_MOTION_START = "R" // Beginning Z Motion (G83)
const TOKEN_SPINDLE_RPM = "S"    // Spindle RPM (G97)
const TOKEN_TOOL_CHANGE = "T"    // Tool Change (G00)

// X, Y, Z motion dimensions
const TOKEN_X_MOTION = "X" // X Axis Motion Dimension
const TOKEN_Y_MOTION = "Y" // Y Axis Motion Dimension
const TOKEN_Z_MOTION = "Z" // Z Axis Motion Dimension

const TOKEN_COMMENT_SUBSTRING = "COMMENT"
const TOKEN_FOREIGN_SUBSTRING = "FOREIGN"

type Parameters map[string]string

type ASTNode struct {
	Token      int
	Letter     string
	Parameters *Parameters
}

type AST struct {
	Nodes []*ASTNode
}

func (a *AST) AddNode(node *ASTNode) {
	a.Nodes = append(a.Nodes, node)
}
