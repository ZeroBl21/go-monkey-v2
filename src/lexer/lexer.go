package lexer

import (
	"errors"

	"github.com/ZeroBl21/go-monkey/src/token"
)

type Lexer struct {
	input        string
	position     int  // Current position in input
	readPosition int  // Current reading position in input (after current char)
	ch           byte // current char under examination
}

func New(input string) *Lexer {
	l := &Lexer{
		input:        input,
		position:     0,
		readPosition: 0,
		ch:           0,
	}
	l.readChar()

	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}

	l.position = l.readPosition
	l.readPosition++
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}

	return l.input[l.readPosition]
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()
	if l.ch == '/' && l.peekChar() == '/' {
		l.skipComment()
	}

	switch l.ch {

	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{
				Type:    token.EQ,
				Literal: literal,
			}
		} else {
			tok = newToken(token.ASSIGN, l.ch)
		}

	// Identifiers + literals
	case '"':
		str, err := l.readString()
		if err != nil {
			tok.Type = token.ILLEGAL
			tok.Literal = err.Error()
		} else {
			tok.Type = token.STRING
			tok.Literal = str
		}

	// Operators
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{
				Type:    token.NOT_EQ,
				Literal: literal,
			}
		} else {
			tok = newToken(token.BANG, l.ch)
		}
	case '/':
		tok = newToken(token.SLASH, l.ch)
	case '*':
		tok = newToken(token.ASTERISK, l.ch)
	case '<':
		tok = newToken(token.LT, l.ch)
	case '>':
		tok = newToken(token.RT, l.ch)

	// Delimiters
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case ':':
		tok = newToken(token.COLON, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case '[':
		tok = newToken(token.LBRACKET, l.ch)
	case ']':
		tok = newToken(token.RBRACKET, l.ch)

	// Default
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		}

		if isDigit(l.ch) {
			digit, err := l.readNumber()
			if err != nil {
				tok.Type = token.ILLEGAL
				tok.Literal = err.Error()
			} else {
				tok.Type = token.INT
				tok.Literal = digit
			}
			return tok
		}

		tok = newToken(token.ILLEGAL, l.ch)
	}

	l.readChar()
	return tok
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) skipComment() {
	for l.ch != '\n' || l.ch == 0 {
		l.readChar()
	}
	l.skipWhitespace()
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}

	return l.input[position:l.position]
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func (l *Lexer) readString() (string, error) {
	position := l.position + 1
	for {
		l.readChar()
		if l.ch == '"' {
			break
		}

		if l.ch == 0 {
			return "", errors.New("unterminated string")
		}
	}

	return l.input[position:l.position], nil
}

func (l *Lexer) readNumber() (string, error) {
	position := l.position
	previousCharWasUnderscore := false

	for isDigit(l.ch) {
		if l.ch == '_' {
			if previousCharWasUnderscore || l.position == position {
				return "", errors.New("invalid number: trailing underscore")
			}
			previousCharWasUnderscore = true
		} else {
			previousCharWasUnderscore = false
		}

		l.readChar()
	}

	if previousCharWasUnderscore {
		return "", errors.New("invalid number: trailing underscore")
	}

	return l.input[position:l.position], nil
}

func isDigit(ch byte) bool {
	return ('0' <= ch && ch <= '9') || ch == '_'
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{
		Type:    tokenType,
		Literal: string(ch),
	}
}
