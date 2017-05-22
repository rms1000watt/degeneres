package generate

import (
	"bytes"
	"fmt"
	"io/ioutil"
)

func Scan(filepath string) (out string, err error) {
	fmt.Println("Starting Scan...")

	fileBytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		fmt.Errorf("Failed reading file: %s: %s", filepath, err)
		return
	}

	s := Scanner{
		Tokens:     []Token{},
		LineNumber: 1,
		InputBuf:   bytes.NewBuffer(fileBytes),
	}

	s.State = s.FileState
	for s.State != nil {
		s.State = s.State()
	}

	return
}

type State func() State

type Token struct {
	LineNumber int
	Name       string
	Value      string
}

type Scanner struct {
	State        State
	LineNumber   int
	Tokens       []Token
	CurrentValue []rune
	InputBuf     *bytes.Buffer
}

var eof = rune(0)

func (s Scanner) read() rune {
	r, _, err := s.InputBuf.ReadRune()
	if err != nil {
		return eof
	}
	return r
}

func (s Scanner) FileState() State {
	r := s.read()
	if r == eof {
		return nil
	}

	inType := []rune{}
	inVal := []rune{}
	if isLetter(r) {
		// Parse inType
		inType = append(inType, r)
		for {
			ru := s.read()

			if isWhitespace(ru) {
				continue
			}

			if ru == eof || isEqual(ru) || isDoubleQuote(ru) {
				break
			}

			inType = append(inType, ru)
		}

		// Check the inType.. if Option, do another round for k, v

		// Parse inVal
		for {
			ru := s.read()

			if isWhitespace(ru) || isDoubleQuote(ru) {
				continue
			}

			if ru == eof || isSemicolon(ru) {
				break
			}

			inVal = append(inVal, ru)
		}

		// Handle the type and val
		fmt.Println(string(inType), string(inVal))
	}

	return s.FileState
}

func (s Scanner) ServiceState(r rune) State {
	return nil
}

func (s Scanner) MessageState(r rune) State {
	return nil
}

func isWhitespace(r rune) bool {
	return r == ' ' ||
		r == '\t' ||
		r == '\r' ||
		r == '\n'
}

func isNewline(r rune) bool {
	return r == '\n'
}

func isNumber(r rune) bool {
	return r > 47 && r < 58
}

func isLetter(r rune) bool {
	return (r > 64 && r < 91) || (r > 96 && r < 123)
}

func isEqual(r rune) bool {
	return r == '='
}

func isDoubleQuote(r rune) bool {
	return r == '"'
}

func isSemicolon(r rune) bool {
	return r == ';'
}
