package gcode

import (
	"errors"
	"os"
	"unicode"
)

type TreeState struct {
	Parameters map[rune]string
	Previous   *TreeState
}

func (t *TreeState) SetParameter(key rune, value string) {
	t.Parameters[key] = value
}

func (t *TreeState) GetParameter(key rune) string {
	// Check if the key exists in the map
	if value, exists := t.Parameters[key]; exists {
		return value
	}
	// Return "0" if the key is not found
	return "0"
}

func (t *TreeState) HasParameter(key rune) bool {
	_, ok := t.Parameters[key]
	return ok
}

func (t *TreeState) Extend(parameters map[rune]string) TreeState {
	// Create a new map to store the merged parameters
	newParams := make(map[rune]string)

	// First, copy the original parameters into the new map
	for k, v := range t.Parameters {
		newParams[k] = v
	}

	// Then, override with the new parameters
	for k, v := range parameters {
		newParams[k] = v
	}

	// Return a new TreeState with the merged parameters
	return TreeState{
		Parameters: newParams,
		Previous:   t,
	}
}

// Token Literals
const TL_JOBSTART = "JOBSTART"         // G-code line that starts the job
const TL_JOBEND = "JOBEND"             // G-code line that ends the job
const TL_SPINDLESTART = "SPINDLESTART" // G-code line that starts the spindle
const TL_SPINDLEEND = "SPINDLEEND"     // G-code line that ends the spindle
const TL_TOOLCHANGE = "TOOLCHANGE"     // G-code line that changes the tool
const TL_GANGON = "GANGON"             // G-code line that turns on the gang tooling
const TL_GANGOFF = "GANGOFF"           // G-code line that turns off the gang tooling
const TL_GANGSET = "GANGSET"           // G-code line that sets the gang tooling
const TL_BEGINCUT = "BEGINCUT"         // G-code line that starts the cut
const TL_ENDCUT = "ENDCUT"             // G-code line that ends the cut

type TokenValue string

func NewTokenValue(value string) (TokenValue, error) {
	var decimalFound bool
	for _, ch := range value {
		if ch == '.' {
			if decimalFound {
				return TokenValue(""), errors.New("invalid token value")
			}
			decimalFound = true
		} else if !unicode.IsDigit(ch) {
			return TokenValue(""), errors.New("invalid token value")
		}
	}
	return TokenValue(value), nil
}

type Token struct {
	Rune  rune
	Value TokenValue
	State TreeState
}

func RuneTokenFromSubstring(substring string) (Token, error) {
	tv, err := NewTokenValue(substring[1:])
	if err != nil {
		return Token{}, err
	}
	return Token{Rune: rune(substring[0]), Value: tv}, nil
}

type Instruction struct {
	LineNumber int
	Tokens     []Token
}

type Command struct {
	LineNumber   int
	LineCount    int
	Instructions []Instruction
}

type Tree struct {
	Path     string
	Commands []Command
}

func NewTree(path string) *Tree {
	return &Tree{Path: path}
}

func tokenize(text string) []Token {
	var result []Token
	r := '#'
	s := ""
	state := TreeState{Parameters: make(map[rune]string)}
	bracketDepth := 0
	quoteCount := 0

	for _, ch := range text {

		// ignore stuff inside brackets, braces, and parens
		if ch == '(' || ch == '[' || ch == '{' {
			bracketDepth++
			continue
		} else if ch == ')' || ch == ']' || ch == '}' {
			bracketDepth--
			continue
		} else if bracketDepth > 0 {
			continue
		}

		// ignore stuff inside quotes
		if ch == '"' {
			quoteCount++
		}
		if quoteCount%2 == 1 {
			continue
		}

		// ignore spaces
		if unicode.IsSpace(ch) && ch != '\n' && ch != '\r' {
			continue

			// handle new lines and carriage returns
		} else if ch == '\n' || ch == '\r' {
			t, err := handleTokenization(r, s, state)
			if err == nil {
				result = append(result, t)
				state = t.State
			}
			result = append(result, Token{Rune: ch, Value: ""})
			s = ""
			r = '#'
		} else if ch == ';' {
			t, err := handleTokenization(r, s, state)
			if err == nil {
				result = append(result, t)
				state = t.State
			}
			r = ch
			s = ""
		} else if unicode.IsLetter(ch) && r != ';' {
			t, err := handleTokenization(r, s, state)
			if err == nil {
				result = append(result, t)
				state = t.State
			}
			r = ch
			s = ""
		} else {
			s += string(ch)
		}
	}
	return result
}

func handleTokenization(r rune, s string, state TreeState) (Token, error) {
	if s != "" {
		t := Token{Rune: r, Value: TokenValue(s)}
		if TokenIsParameter(t) {
			t.State = state.Extend(map[rune]string{r: s})
		}
		return t, nil
	} else {
		return Token{}, errors.New("empty token value")
	}
}

func (t *Tree) Parse() error {
	text, err := t.GetFileText()
	if err != nil {
		return err
	}

	tokens := tokenize(text)

	println(len(tokens))

	ln := 0
	nl := false
	com := Command{LineNumber: ln}
	ins := Instruction{LineNumber: ln}

	for _, tk := range tokens {
		if tk.Rune == '\n' || tk.Rune == '\r' {
			// new line
			ln++
			nl = true
			if isInstructionValid(ins) {
				com.Instructions = append(com.Instructions, ins)
			}
			ins = Instruction{LineNumber: ln}
		} else if tk.Rune == ';' || unicode.IsLetter(tk.Rune) {

			if nl && (tk.Rune == 'G' || tk.Rune == 'M') { // If a new line starts with a G or M code, then a new command block has begun
				com.LineCount = ln - com.LineNumber
				if isInstructionValid(ins) {
					com.Instructions = append(com.Instructions, ins)
				}
				if len(com.Instructions) > 0 {
					t.Commands = append(t.Commands, com)
				}
				com = Command{LineNumber: ln}
				ins = Instruction{LineNumber: ln}
				nl = false
			}
			ins.Tokens = append(ins.Tokens, tk)
		}
	}
	return nil
}

func (t *Tree) GetFileText() (string, error) {
	// Open the file
	data, err := os.ReadFile(t.Path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func isInstructionValid(ins Instruction) bool {
	if len(ins.Tokens) == 0 {
		return false
	}
	t := ins.Tokens[0].Rune
	if unicode.IsLetter(t) {
		return true
	} else if t == ';' {
		return true
	}
	return false
}

func TokenIsParameter(token Token) bool {
	if token.Rune == 'G' || token.Rune == 'M' {
		return false
	} else if unicode.IsLetter(token.Rune) {
		return true
	} else {
		return false
	}
}

func TokenIsComment(token Token) bool {
	return token.Rune == ';'
}

func TokenIsCode(token Token) bool {
	return token.Rune == 'G' || token.Rune == 'M'
}
