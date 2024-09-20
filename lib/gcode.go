package gcode

import (
	"fmt"
	"unicode"
)

const TOKEN_EOF = -1
const TOKEN_SLEW = 0
const TOKEN_MACHINE = 1
const TOKEN_CWMOVE_2D = 2
const TOKEN_CCWMOVE_2D = 3
const TOKEN_DWELL = 4

const TOKEN_FIND_HOME = 37
const TOKEN_INCHES_PROGRAMMING = 70
const TOKEN_MILLIMETERS_PROGRAMMING = 71

const TOKEN_CWMOVE_3D = 72
const TOKEN_CCWMOVE_3D = 73
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

const TOKEN_LITERAL_COMMENT = ';'
const TOKEN_LITERAL_PARENS_OPEN = '('
const TOKEN_LITERAL_PARENS_CLOSE = ')'
const TOKEN_LITERAL_QUOTE = '"'
const TOKEN_LITERAL_SINGLE_QUOTE = '\''
const TOKEN_LITERAL_BRACKET_OPEN = '['
const TOKEN_LITERAL_BRACKET_CLOSE = ']'
const TOKEN_LITERAL_BRACE_OPEN = '{'
const TOKEN_LITERAL_BRACE_CLOSE = '}'
const TOKEN_LITERAL_LINE_BREAK = '\n'
const TOKEN_LITERAL_CARRIAGE_RETURN = '\r'
const TOKEN_LITERAL_M_CODE = 'M'
const TOKEN_LITERAL_G_CODE = 'G'

func ParseSubstring(code int, letter, substring string) (*ASTNode, error) {
	var node *ASTNode = &ASTNode{
		Token:      code,
		Letter:     letter,
		Parameters: &Parameters{},
	}
	return node, nil
}

// ScanString scans a string character by character, skipping whitespace
func Lexer(s string) *AST {
	var ast *AST = &AST{}

	var is_inside_comment bool
	var block_depth int = 0
	var is_inside_quotes bool
	var is_inside_single_quotes bool
	var is_inside_g_code bool
	var is_inside_m_code bool

	var current_type string
	var current_code int
	var current_substr string
	for _, ch := range s {
		var is_eol bool = ch == TOKEN_LITERAL_LINE_BREAK || ch == TOKEN_LITERAL_CARRIAGE_RETURN
		if is_inside_g_code || is_inside_m_code {
			if is_eol {
				is_inside_g_code = false
				is_inside_m_code = false
				node, err := ParseSubstring(current_code, current_type, current_substr)
				if err != nil {
					fmt.Println(err)
				} else {
					ast.AddNode(node)
				}
			} else {
				current_substr += string(ch)
			}
		} else {
			if unicode.IsSpace(ch) {
				continue // Skip whitespace characters
			} else if block_depth > 0 {
				continue
			} else if is_inside_comment {
				if is_eol {
					is_inside_comment = false
				} else {
					continue
				}
			} else if is_inside_quotes {
				if ch == TOKEN_LITERAL_QUOTE {
					is_inside_quotes = false
				} else {
					continue
				}
			} else if is_inside_single_quotes {
				if ch == TOKEN_LITERAL_SINGLE_QUOTE {
					is_inside_single_quotes = false
				} else {
					continue
				}
			} else if ch == TOKEN_LITERAL_PARENS_OPEN || ch == TOKEN_LITERAL_BRACE_OPEN || ch == TOKEN_LITERAL_BRACKET_OPEN {
				block_depth++
			} else if ch == TOKEN_LITERAL_PARENS_CLOSE || ch == TOKEN_LITERAL_BRACE_CLOSE || ch == TOKEN_LITERAL_BRACKET_CLOSE {
				block_depth--
			} else if ch == TOKEN_LITERAL_QUOTE {
				is_inside_quotes = true
			} else if ch == TOKEN_LITERAL_SINGLE_QUOTE {
				is_inside_single_quotes = true
			} else if ch == TOKEN_LITERAL_COMMENT {
				is_inside_comment = true
			} else if ch == TOKEN_LITERAL_G_CODE {
				is_inside_g_code = true
				fmt.Println("Found G-code command")
			} else if ch == TOKEN_LITERAL_M_CODE {
				is_inside_m_code = true
				fmt.Println("Found M-code command")
			} else if ch == TOKEN_LITERAL_LINE_BREAK || ch == TOKEN_LITERAL_CARRIAGE_RETURN {
				is_inside_g_code = false
				is_inside_m_code = false
				is_inside_comment = false
			}
		}
	}
	return ast
}
