package gcode

import (
	"errors"
	"os"
	"unicode"
)

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
	r := '0'
	s := ""

	for _, ch := range text {
		if unicode.IsSpace(ch) && ch != '\n' && ch != '\r' {
			continue
		} else if ch == '\n' || ch == '\r' {
			t, err := handleTokenization(r, s)
			if err == nil {
				result = append(result, t)
			}
			result = append(result, Token{Rune: ch, Value: ""})
			s = ""
			r = '0'
		} else if ch == ';' {
			t, err := handleTokenization(r, s)
			if err == nil {
				result = append(result, t)
			}
			r = ch
			s = ""
		} else if unicode.IsLetter(ch) && r != ';' {
			t, err := handleTokenization(r, s)
			if err == nil {
				result = append(result, t)
			}
			r = ch
			s = ""
		} else {
			s += string(ch)
		}
	}
	return result
}

func handleTokenization(r rune, s string) (Token, error) {
	if s != "" {
		return Token{Rune: r, Value: TokenValue(s)}, nil
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
			if len(ins.Tokens) > 0 {
				com.Instructions = append(com.Instructions, ins)
			}
			ins = Instruction{LineNumber: ln}
		} else if tk.Rune == ';' || unicode.IsLetter(tk.Rune) {

			if nl && (tk.Rune == 'G' || tk.Rune == 'M') { // If a new line starts with a G or M code, then a new command block has begun
				com.LineCount = ln - com.LineNumber
				if len(ins.Tokens) > 0 {
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
