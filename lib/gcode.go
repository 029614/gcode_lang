package lib

import (
	"errors"
	"log"
	"strings"
	"unicode"
)

const TOKEN_SLEW = 0
const TOKEN_MACHINE = 1
const TOKEN_CW_MOVE_2D = 2
const TOKEN_CCW_MOVE_2D = 3
const TOKEN_DWELL = 4

const TOKEN_GANG_ON = "GANGON"

const TOKEN_FIND_HOME = 37
const TOKEN_INCHES_PROGRAMMING = 70
const TOKEN_MILLIMETERS_PROGRAMMING = 71

const TOKEN_CW_MOVE_3D = 72
const TOKEN_CCW_MOVE_3D = 73
const TOKEN_ARC_INCREMENTAL_MODE = 74
const TOKEN_ARC_ABSOLUTE_MODE = 75

const TOKEN_ONE_STROKE_DRILL_CYCLE = 81
const TOKEN_PECK_DRILL_CYCLE = 83
const TOKEN_TAP_DRILL_CYCLE = 84

const TOKEN_ABSOLUTE_COORDINATE_MODE = 90
const TOKEN_INCREMENTAL_COORDINATE_MODE = 91
const TOKEN_SOFT_HOME = 92
const TOKEN_SPINDLE_SPEED = 97

const TOKEN_FOREIGN = 1000

const RUNE_COMMENT = ';'
const RUNE_LINEBREAK = '\n'
const RUNE_CARRIAGE_RETURN = '\r'
const RUNE_CODE_M = 'M'
const RUNE_CODE_G = 'G'

// Tool-specific and motion-related tokens
const RUNE_TOOL_CHANGE_OPERATOR_MESSAGE = 'C' // Tool Change Operator Message
const RUNE_PECK_DRILL_DATA = 'D'              // Peck Drill Data
const RUNE_FEEDRATE = 'F'                     // Feedrate in Units per Second

// G-code commands
const RUNE_PREPARATORY_FUNCTION = 'G' // Preparatory Function

// Circular Interpolation Values
const RUNE_CIRCULAR_INTERPOLATION_X = 'I' // Circular Interpolation X (G02, G03)
const RUNE_CIRCULAR_INTERPOLATION_Y = 'J' // Circular Interpolation Y (G02, G03)
const RUNE_CIRCULAR_INTERPOLATION_Z = 'K' // Circular Interpolation Z (G02, G03)

// Miscellaneous or control function
const RUNE_MISC_CONTROL = 'M' // Miscellaneous Control Function

// Sequence number
const RUNE_SEQUENCE_NUMBER = 'N' // Sequence Number

// Motion parameters
const RUNE_Z_MOTION_START = 'R' // Beginning Z Motion (G83)
const RUNE_SPINDLE_RPM = 'S'    // Spindle RPM (G97)
const RUNE_TOOL_CHANGE = 'T'    // Tool Change (G00)

// X, Y, Z motion dimensions
const RUNE_X_MOTION = 'X' // X Axis Motion Dimension
const RUNE_Y_MOTION = 'Y' // Y Axis Motion Dimension
const RUNE_Z_MOTION = 'Z' // Z Axis Motion Dimension

const RUNE_COMMENT_SUBSTRING = 'A' // The First unused letter available
const RUNE_FOREIGN_SUBSTRING = 'B' // The next unused letter available

// Parameters holds G-code parameters like X, Y, Z, etc.
type Parameters map[string]string

// ASTNode represents a G-code node in the AST
type ASTNode struct {
	Code       string
	Rune       rune
	Parameters *Parameters
}

// AST represents the collection of G-code nodes
type AST struct {
	Nodes []*ASTNode
}

// IsValidInt checks if a string is a valid integer
func IsValidInt(s string) bool {
	for _, ch := range s {
		if !unicode.IsDigit(ch) {
			return false
		}
	}
	return len(s) > 0
}

// IsValidFloat checks if a string is a valid floating-point number
func IsValidFloat(s string) bool {
	decimalFound := false
	for _, ch := range s {
		if ch == '.' {
			if decimalFound {
				return false
			}
			decimalFound = true
		} else if !unicode.IsDigit(ch) {
			return false
		}
	}
	return len(s) > 0
}

// ParseCode parses a G-code or M-code substring
func ParseCode(substring string) (rune, string, string, error) {
	if len(substring) < 3 {
		return 'G', "-1", "", errors.New("substring too short")
	}

	r := rune(substring[0])
	if r != RUNE_CODE_G && r != RUNE_CODE_M {
		return 'G', "-1", "", errors.New("substring must start with 'M' or 'G'")
	}

	code := substring[1:3]
	if !IsValidInt(code) {
		return 'G', "-1", "", errors.New("invalid code, should be two digits")
	}

	parameterSubstring := substring[3:]
	return r, code, parameterSubstring, nil
}

// IsEndOfCode checks if the current G-code substring has ended
func IsEndOfCode(ch rune) bool {
	return unicode.IsSpace(ch) || ch == RUNE_COMMENT || ch == RUNE_CODE_G || ch == RUNE_CODE_M
}

// ParseParameters parses the parameter substring of G-code (e.g., "X10 Y20")
func ParseParameters(s string) Parameters {
	params := Parameters{}
	var currentParam rune
	var currentValue strings.Builder
	for _, ch := range s {
		if unicode.IsLetter(ch) {
			if currentParam != 0 && currentValue.Len() > 0 {
				params[string(currentParam)] = currentValue.String()
			}
			currentParam = ch
			currentValue.Reset()
		} else if unicode.IsDigit(ch) || ch == '.' {
			currentValue.WriteRune(ch)
		}
	}
	if currentParam != 0 && currentValue.Len() > 0 {
		params[string(currentParam)] = currentValue.String()
	}
	return params
}

// NewASTNode creates a new AST node
func NewASTNode(token string, letter rune, parameters *Parameters) *ASTNode {
	return &ASTNode{
		Code:       token,
		Rune:       letter,
		Parameters: parameters,
	}
}

// NewAST creates a new AST
func NewAST() *AST {
	return &AST{
		Nodes: []*ASTNode{},
	}
}

// AddNode adds a node to the AST
func (a *AST) AddNode(node *ASTNode) {
	a.Nodes = append(a.Nodes, node)
}

// ProcessSubstring processes a G-code substring and adds it to the AST
func (a *AST) ProcessSubstring(substring string) error {
	rcode, scode, parameterSubstring, err := ParseCode(substring)
	if err != nil {
		return err
	}
	params := ParseParameters(parameterSubstring)
	a.AddNode(NewASTNode(scode, rcode, &params))
	return nil
}

// Lexer processes the input string into tokens and builds the AST
func Lexer(s string) {
	ast := NewAST()
	var isInsideComment bool
	var isInsideCode bool
	var substring strings.Builder

	for _, ch := range s {
		if isInsideComment {
			if ch == RUNE_LINEBREAK || ch == RUNE_CARRIAGE_RETURN {
				// End of comment
				if err := ast.ProcessSubstring(substring.String()); err != nil {
					log.Printf("Error processing comment substring: %v", err)
				}
				substring.Reset()
				isInsideComment = false
			} else {
				// Continue accumulating comment
				substring.WriteRune(ch)
				continue
			}
		} else if isInsideCode {
			if IsEndOfCode(ch) {
				// Process current G-code or M-code command
				if err := ast.ProcessSubstring(substring.String()); err != nil {
					log.Printf("Error processing G-code substring: %v", err)
				}
				substring.Reset()
				isInsideCode = false

				// Start new command or enter comment
				if ch == RUNE_COMMENT {
					isInsideComment = true
				} else if ch == RUNE_CODE_G || ch == RUNE_CODE_M {
					// New G/M code starts here
					isInsideCode = true
					substring.WriteRune(ch)
				}
			} else {
				// Continue accumulating G-code or M-code command
				substring.WriteRune(ch)
				continue
			}
		} else {
			// Not inside any code or comment, check for new G/M or comment start
			if ch == RUNE_COMMENT {
				isInsideComment = true
			} else if ch == RUNE_CODE_G || ch == RUNE_CODE_M {
				// New G/M code starts here
				isInsideCode = true
				substring.WriteRune(ch)
			}
		}
	}

	// Process any remaining substring at the end of the input
	if substring.Len() > 0 {
		if err := ast.ProcessSubstring(substring.String()); err != nil {
			log.Printf("Error processing remaining substring: %v", err)
		}
	}
}
