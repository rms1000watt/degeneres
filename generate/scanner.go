package generate

import (
	"bytes"
	"fmt"
	"io/ioutil"
)

const (
	TokenNameOption  = "option"
	TokenNameKey     = "key"
	TokenNameVal     = "val"
	TokenNameRPCName = "rpcName"
	TokenNameRPCIn   = "rpcIn"
	TokenNameRPCOut  = "rpcOut"
)

func Scan(filepath string) (out string, err error) {
	fmt.Println("Starting Scan...")

	fileBytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		fmt.Errorf("Failed reading file: %s: %s", filepath, err)
		return
	}
	tokenCh := make(chan Token)

	s := NewScanner(fileBytes, tokenCh)
	go s.Start()

	for token := range tokenCh {
		fmt.Println(token)
	}

	return
}

type State func() State

type Token struct {
	Name  string
	Value string
}

type Scanner struct {
	State    State
	TokenCh  chan Token
	InputBuf *bytes.Buffer
}

var eof = rune(0)

func NewScanner(inputBytes []byte, tokenCh chan Token) Scanner {
	return Scanner{
		InputBuf: bytes.NewBuffer(inputBytes),
		TokenCh:  tokenCh,
	}
}

func (s Scanner) Start() {
	s.State = s.FileState
	for s.State != nil {
		s.State = s.State()
	}
	close(s.TokenCh)
}

func (s Scanner) Emit(token Token) {
	// Change this to publish to a channel
	// s.Tokens = append(s.Tokens, token)
	s.TokenCh <- token
}

func (s Scanner) read() rune {
	r, _, err := s.InputBuf.ReadRune()
	if err != nil {
		return eof
	}
	return r
}

func (s Scanner) unread() {
	err := s.InputBuf.UnreadRune()
	if err != nil {
		fmt.Println("Error unreading rune:", err)
	}
}

func (s Scanner) FileState() State {
	r := s.read()

	for !isLetter(r) {
		if r == eof {
			return nil
		}

		r = s.read()
	}

	// Scan for key
	key := s.getKey(r)
	s.Emit(Token{
		Name:  TokenNameKey,
		Value: string(key),
	})

	if string(key) == "service" {
		val := s.getVal()
		s.Emit(Token{
			Name:  TokenNameVal,
			Value: string(val),
		})

		return s.ServiceState
	}

	// Scan for option
	if string(key) == "option" {
		option := s.getSingleOption()
		s.Emit(Token{
			Name:  TokenNameOption,
			Value: string(option),
		})
	}

	// Scan for val
	val := s.getVal()
	s.Emit(Token{
		Name:  TokenNameVal,
		Value: string(val),
	})

	return s.FileState
}

func (s Scanner) ServiceState() State {
	r := s.read()

	for !isLetter(r) {
		if r == eof {
			return nil
		}

		if isCloseCurleyBrace(r) {
			return s.FileState
		}

		r = s.read()
	}

	// Check for service options
	key := s.getKey(r)
	s.Emit(Token{
		Name:  TokenNameKey,
		Value: string(key),
	})

	if string(key) == "rpc" {
		s.getRPCVals()
		return s.RPCState
	}

	// Scan for option
	if string(key) == "option" {
		option := s.getSingleOption()
		s.Emit(Token{
			Name:  TokenNameOption,
			Value: string(option),
		})
	}

	// Scan for val
	val := s.getVal()
	s.Emit(Token{
		Name:  TokenNameVal,
		Value: string(val),
	})

	return s.ServiceState
}

func (s Scanner) RPCState() State {
	r := s.read()

	for !isLetter(r) {
		if r == eof {
			return nil
		}

		if isCloseCurleyBrace(r) {
			return s.ServiceState
		}

		r = s.read()
	}

	// Check for service options
	key := s.getKey(r)
	s.Emit(Token{
		Name:  TokenNameKey,
		Value: string(key),
	})

	// Scan for option
	if string(key) == "option" {
		option := s.getSingleOption()
		s.Emit(Token{
			Name:  TokenNameOption,
			Value: string(option),
		})
	}

	// Scan for val
	val := s.getVal()
	s.Emit(Token{
		Name:  TokenNameVal,
		Value: string(val),
	})

	return s.RPCState
}

func (s Scanner) MessageState() State {
	r := s.read()
	if r == eof {
		return nil
	}

	return nil
}

func (s Scanner) getKey(runes ...rune) (key []rune) {
	for _, r := range runes {
		key = append(key, r)
	}

	for {
		r := s.read()

		if r == eof || isEqual(r) || isWhitespace(r) {
			break
		}

		if isDoubleQuote(r) {
			s.unread()
			break
		}

		key = append(key, r)
	}
	return
}

func (s Scanner) getSingleOption() (option []rune) {
	for {
		r := s.read()

		if r == eof || isCloseParen(r) {
			break
		}

		if isWhitespace(r) || isOpenParen(r) {
			continue
		}

		option = append(option, r)
	}
	return
}

func (s Scanner) getVal(runes ...rune) (val []rune) {
	for _, r := range runes {
		val = append(val, r)
	}

	for {
		r := s.read()

		if isEqual(r) || isWhitespace(r) {
			continue
		}

		if isDoubleQuote(r) || isBacktick(r) {
			val = s.scanLiteral()
			continue
		}

		if r == eof || isSemicolon(r) || isOpenCurleyBrace(r) {
			break
		}

		val = append(val, r)
	}
	return
}

func (s Scanner) scanLiteral() (literal []rune) {
	for {
		r := s.read()

		if r == eof || isDoubleQuote(r) || isBacktick(r) {
			break
		}

		literal = append(literal, r)
	}
	return
}

func (s Scanner) getRPCVals() {
	var r rune

	// Parse RPC Name
	rpcName := []rune{}
	for {
		r = s.read()

		if isWhitespace(r) {
			continue
		}

		if isEOF(r) || isOpenParen(r) {
			break
		}

		rpcName = append(rpcName, r)
	}

	s.Emit(Token{
		Name:  TokenNameRPCName,
		Value: string(rpcName),
	})

	// Parse RPC In Func
	rpcIn := []rune{}
	for {
		r = s.read()

		if isWhitespace(r) {
			continue
		}

		if isEOF(r) || isCloseParen(r) {
			break
		}

		rpcIn = append(rpcIn, r)
	}

	s.Emit(Token{
		Name:  TokenNameRPCIn,
		Value: string(rpcIn),
	})

	// Skip the `returns` and whitespaces
	for {
		r = s.read()

		if isEOF(r) || isOpenParen(r) {
			break
		}
	}

	// Parse RPC Out func
	rpcOut := []rune{}
	for {
		r = s.read()

		if isWhitespace(r) {
			continue
		}

		if isEOF(r) || isCloseParen(r) {
			break
		}

		rpcOut = append(rpcOut, r)
	}

	s.Emit(Token{
		Name:  TokenNameRPCOut,
		Value: string(rpcOut),
	})
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
	return (r > 64 && r < 91) ||
		(r > 96 && r < 123)
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

func isOpenParen(r rune) bool {
	return r == '('
}

func isCloseParen(r rune) bool {
	return r == ')'
}

func isBacktick(r rune) bool {
	return r == '`'
}

func isOpenCurleyBrace(r rune) bool {
	return r == '{'
}

func isCloseCurleyBrace(r rune) bool {
	return r == '}'
}

func isEOF(r rune) bool {
	return r == eof
}
