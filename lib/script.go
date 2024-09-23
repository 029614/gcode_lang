package gcode

import (
	"errors"
	"os"
	"strings"
	"unicode"
)

type Code struct {
	String        string
	Code          rune
	Numeral       string
	Parameters    map[rune]string
	InlineCode    *Code
	InlineComment *Code
}

func FromCode(line string) (Code, error) {
	if strings.HasPrefix(line, "G") || strings.HasPrefix(line, "M") {
		code := Code{String: line}
		code.ParseCode(line)
		return code, nil
	}
	return Code{}, errors.New("not a valid code")
}

func FromComment(line string) (Code, error) {
	if strings.HasPrefix(line, ";") {
		return Code{String: line, Code: ';'}, nil
	}
	return Code{}, errors.New("not a valid comment")
}

func FromUndefined(line string) (Code, error) {
	// Find the ';' substring and add everything after it to InlineComment
	if index := strings.Index(line, ";"); index != -1 {
		comment := line[index:]
		return Code{
			String:        line,
			Code:          '?',
			InlineComment: &Code{String: comment, Code: ';'},
		}, nil
	}
	return Code{String: line, Code: '?'}, nil // No comment found
}

// NewCode parses a G-code line and returns a Code struct with the extracted components.
func (c *Code) ParseCode(line string) error {
	// Initialize the Code struct
	c.Parameters = make(map[rune]string)

	// Trim whitespace from the line
	line = strings.TrimSpace(line)

	for i, ch := range line {
		// ensure first rune is a G or M
		if i == 0 {
			if ch != 'G' && ch != 'M' {
				return errors.New("invalid code, should start with 'G' or 'M'")
			}
			c.Code = ch
			continue
		}

		// ensure next two runes are digits
		if i > 0 && i < 3 {
			if !unicode.IsDigit(ch) {
				return errors.New("invalid code, should be two digits")
			}
			c.Numeral += string(ch)
			continue
		}

		var trackingParamter bool = false
		var p rune

		// Process the line
		if i > 2 {
			// skip whitespace
			if unicode.IsSpace(ch) {
				continue

				// Handle inline comments
			} else if ch == ';' {
				ic, err := FromComment(line[i:])
				if err == nil {
					c.InlineComment = &ic
				}
				return nil

				// Handle inline codes
			} else if ch == 'G' || ch == 'M' {
				inline, err := FromCode(line[i:])
				if err != nil {
					c.InlineCode = &inline
					c.InlineComment = c.InlineCode.InlineComment // Raise comment to root code
					return nil
				} else {
					return err
				}

				// Handle parameters
			} else if unicode.IsLetter(ch) {
				trackingParamter = true
				p = ch
				c.Parameters[p] = ""
			} else if trackingParamter {
				c.Parameters[p] += string(ch)
			}
		}
	}

	return nil
}

type Script struct {
	Path       string
	SyntaxTree []Code
}

func NewScript(path string) *Script {
	return &Script{Path: path}
}

func (s *Script) GetUndefined() []Code {
	var undefined []Code
	for _, code := range s.SyntaxTree {
		if code.Code == '?' {
			undefined = append(undefined, code)
		}
	}
	return undefined
}

func (s *Script) GetComments(includeInline bool) []Code {
	var comments []Code
	for _, code := range s.SyntaxTree {
		if code.InlineComment != nil && includeInline {
			comments = append(comments, *code.InlineComment)
		} else if code.Code == ';' {
			comments = append(comments, code)
		}
	}
	return comments
}

func (s *Script) GetGCodes(includeInline bool) []Code {
	var codes []Code
	for _, code := range s.SyntaxTree {
		if code.Code == 'G' {
			codes = append(codes, code)
		}
		if includeInline && code.InlineCode != nil && code.InlineCode.Code == 'G' {
			codes = append(codes, *code.InlineCode)
		}
	}
	return codes
}

func (s *Script) GetMCodes(includeInline bool) []Code {
	var codes []Code
	for _, code := range s.SyntaxTree {
		if code.Code == 'M' {
			codes = append(codes, code)
		}
		if includeInline && code.InlineCode != nil && code.InlineCode.Code == 'M' {
			codes = append(codes, *code.InlineCode)
		}
	}
	return codes
}

func (s *Script) GetFileText() (string, error) {
	// Open the file
	data, err := os.ReadFile(s.Path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// Parse categorizes the script into comments, G/M codes, and other lines.
func (s *Script) Parse() error {
	text, err := s.GetFileText()
	if err != nil {
		return err
	}
	// Declare slices for comments, G/M codes, and other lines

	// Check if the script is empty
	if strings.TrimSpace(text) == "" {
		return errors.New("script is empty")
	}

	// Split the script into lines
	lines := strings.Split(text, "\n")

	// Loop over each line and categorize it
	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)

		// Skip empty lines
		if trimmedLine == "" {
			continue
		} else {

			// Handle comments
			if strings.HasPrefix(trimmedLine, ";") {
				code, err := FromComment(trimmedLine)
				if err != nil {
					return err
				} else {
					s.SyntaxTree = append(s.SyntaxTree, code)
				}

				// Handle G/M codes
			} else if strings.HasPrefix(trimmedLine, "G") || strings.HasPrefix(trimmedLine, "M") {
				code, err := FromCode(trimmedLine)
				if err != nil {
					return err
				} else {
					s.SyntaxTree = append(s.SyntaxTree, code)
				}

				// Handle unknown cases
			} else {
				code, err := FromUndefined(trimmedLine)
				if err != nil {
					return err
				} else {
					s.SyntaxTree = append(s.SyntaxTree, code)
				}
			}
		}
	}

	return nil
}
