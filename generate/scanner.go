package generate

import (
	"bytes"
	"strings"

	log "github.com/sirupsen/logrus"
)

var eof = rune(0)

type State func() State

type Scanner struct {
	State    State
	TokenCh  chan Token
	InputBuf *bytes.Buffer
}

func Scan(fileBytes []byte) (tokens chan Token) {
	log.Debug("Starting scanner")
	defer log.Debug("Scanner done")

	tokens = make(chan Token)
	s := NewScanner(fileBytes, tokens)
	go s.Start()

	return
}

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
	s.TokenCh <- token
}

func (s Scanner) read() rune {
	r, _, err := s.InputBuf.ReadRune()
	if err != nil {
		if err.Error() != "EOF" {
			log.Error("Failed reading rune: ", err)
		}
		return eof
	}
	return r
}

func (s Scanner) unread() {
	err := s.InputBuf.UnreadRune()
	if err != nil {
		log.Error("Failed unreading rune: ", err)
	}
}

func (s Scanner) FileState() State {
	r := s.read()

	for !isLetter(r) {
		if isEOF(r) {
			return nil
		}

		if isForwardSlash(r) {
			r = s.read()
			if isForwardSlash(r) {
				s.scanComment()
			}
		}

		r = s.read()
	}

	// Scan for key
	key := s.getKey(r)

	if string(key) == "service" {
		val := s.getVal()
		s.Emit(Token{
			Name:  TokenServiceKey,
			Value: string(val),
		})

		return s.ServiceState
	}

	if string(key) == "message" {
		val := s.getVal()
		s.Emit(Token{
			Name:  TokenMessageKey,
			Value: string(val),
		})

		return s.MessageState
	}

	// Scan for option
	if string(key) == "option" {
		option := s.getSingleOption()
		s.Emit(Token{
			Name:  TokenFileOptionKey,
			Value: string(option),
		})

		val := s.getVal()
		s.Emit(Token{
			Name:  TokenFileOptionVal,
			Value: string(val),
		})
		return s.FileState
	}

	s.Emit(Token{
		Name:  TokenFileKey,
		Value: string(key),
	})

	// Scan for val
	val := s.getVal()
	s.Emit(Token{
		Name:  TokenFileVal,
		Value: string(val),
	})

	return s.FileState
}

func (s Scanner) ServiceState() State {
	r := s.read()

	for !isLetter(r) {
		if isEOF(r) {
			return nil
		}

		if isForwardSlash(r) {
			r = s.read()
			if isForwardSlash(r) {
				s.scanComment()
			}
		}

		if isCloseCurleyBrace(r) {
			s.Emit(Token{Name: TokenServiceDone})
			return s.FileState
		}

		r = s.read()
	}

	key := s.getKey(r)

	if string(key) == "rpc" {
		s.getRPCVals()
		return s.RPCState
	}

	if string(key) == "option" {
		option := s.getSingleOption()
		s.Emit(Token{
			Name:  TokenServiceOptionKey,
			Value: string(option),
		})

		val := s.getVal()
		s.Emit(Token{
			Name:  TokenServiceOptionVal,
			Value: string(val),
		})
		return s.ServiceState
	}

	return s.ServiceState
}

func (s Scanner) RPCState() State {
	r := s.read()

	for !isLetter(r) {
		if isEOF(r) {
			return nil
		}

		if isForwardSlash(r) {
			r = s.read()
			if isForwardSlash(r) {
				s.scanComment()
			}
		}

		if isCloseCurleyBrace(r) {
			s.Emit(Token{Name: TokenRPCDone})
			return s.ServiceState
		}

		r = s.read()
	}

	key := s.getKey(r)

	if string(key) == "option" {
		option := s.getSingleOption()
		s.Emit(Token{
			Name:  TokenRPCOptionKey,
			Value: string(option),
		})

		val := s.getVal()
		s.Emit(Token{
			Name:  TokenRPCOptionVal,
			Value: string(val),
		})
	}

	return s.RPCState
}

func (s Scanner) MessageState() State {
	r := s.read()

	for !isLetter(r) {
		if isEOF(r) {
			return nil
		}

		if isForwardSlash(r) {
			r = s.read()
			if isForwardSlash(r) {
				s.scanComment()
			}
		}

		if isCloseCurleyBrace(r) {
			s.Emit(Token{Name: TokenMessageDone})
			return s.FileState
		}

		r = s.read()
	}

	dataType := s.getFieldDataType(r)
	if isFieldRule(string(dataType)) {
		s.Emit(Token{
			Name:  TokenFieldRule,
			Value: string(dataType),
		})
		dataType = s.getFieldDataType()
	}
	s.Emit(Token{
		Name:  TokenFieldDataType,
		Value: string(dataType),
	})

	key := s.getFieldKey()
	s.Emit(Token{
		Name:  TokenFieldKey,
		Value: string(key),
	})

	// Scan until the options begin or end of field
	for {
		r := s.read()

		if isEOF(r) {
			return nil
		}

		if isSemicolon(r) {
			s.Emit(Token{Name: TokenFieldDone})
			break
		}

		if isOpenSquareBracket(r) {
			return s.FieldOptionsState
		}
	}

	return s.MessageState
}

func (s Scanner) FieldOptionsState() State {
	// Scan until we get to option key
	for {
		r := s.read()

		if isEOF(r) {
			return nil
		}

		if isSemicolon(r) || isCloseSquareBracket(r) {
			s.Emit(Token{Name: TokenFieldDone})
			return s.MessageState
		}

		if isOpenParen(r) {
			break
		}
	}

	fieldOptionKey := []rune{}
	for {
		r := s.read()

		if isEOF(r) {
			return nil
		}

		if isCloseParen(r) {
			break
		}

		if isWhitespace(r) {
			continue
		}

		fieldOptionKey = append(fieldOptionKey, r)
	}

	s.Emit(Token{
		Name:  TokenFieldOptionKey,
		Value: string(fieldOptionKey),
	})

	// Scan until we get to option val
	for {
		r := s.read()

		if isEOF(r) {
			return nil
		}

		if isCloseSquareBracket(r) || isSemicolon(r) {
			return s.MessageState
		}

		if isDoubleQuote(r) {
			break
		}

		if isWhitespace(r) {
			continue
		}
	}

	fieldOptionVal := []rune{}
	for {
		r := s.read()

		if isEOF(r) {
			return nil
		}

		if isDoubleQuote(r) || isCloseSquareBracket(r) || isSemicolon(r) {
			break
		}

		fieldOptionVal = append(fieldOptionVal, r)
	}

	s.Emit(Token{
		Name:  TokenFieldOptionVal,
		Value: string(fieldOptionVal),
	})

	return s.FieldOptionsState
}

func (s Scanner) getKey(runes ...rune) (key []rune) {
	for _, r := range runes {
		key = append(key, r)
	}

	for {
		r := s.read()

		if isEOF(r) || isEqual(r) || isWhitespace(r) {
			break
		}

		if isDoubleQuote(r) || isCloseCurleyBrace(r) {
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

		if isEOF(r) || isCloseParen(r) {
			break
		}

		if isWhitespace(r) || isOpenParen(r) {
			continue
		}

		option = append(option, r)
	}
	return
}

func (s Scanner) scanComment() {
	for {
		r := s.read()

		if isEOF(r) {
			s.unread()
			break
		}

		if isNewline(r) {
			break
		}
	}
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

		if isEOF(r) || isSemicolon(r) || isOpenCurleyBrace(r) {
			break
		}

		val = append(val, r)
	}
	return
}

func (s Scanner) getFieldDataType(runes ...rune) (dataType []rune) {
	for _, r := range runes {
		dataType = append(dataType, r)
	}

	for {
		r := s.read()

		if isEOF(r) || isWhitespace(r) {
			break
		}

		dataType = append(dataType, r)
	}
	return
}

func (s Scanner) getFieldKey(runes ...rune) (key []rune) {
	for _, r := range runes {
		key = append(key, r)
	}

	for {
		r := s.read()

		if isWhitespace(r) {
			continue
		}

		if isEOF(r) || isSemicolon(r) || isEqual(r) {
			break
		}

		key = append(key, r)
	}
	return
}

func (s Scanner) scanLiteral() (literal []rune) {
	for {
		r := s.read()

		if isEOF(r) || isDoubleQuote(r) || isBacktick(r) {
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
		Name:  TokenRPCName,
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
		Name:  TokenRPCIn,
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
		Name:  TokenRPCOut,
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

func isOpenSquareBracket(r rune) bool {
	return r == '['
}

func isCloseSquareBracket(r rune) bool {
	return r == ']'
}

func isComma(r rune) bool {
	return r == ','
}

func isForwardSlash(r rune) bool {
	return r == '/'
}

func isFieldRule(in string) bool {
	in = strings.ToLower(in)
	return in == FieldRuleOptional ||
		in == FieldRuleRepeated ||
		in == FieldRuleRequired
}
