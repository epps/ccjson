package main

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	Illegal TokenType = "ILLEGAL"
	EOF     TokenType = "EOF"
	// Structural Characters
	BeginArray     TokenType = "["
	BeginObject    TokenType = "{"
	EndArray       TokenType = "]"
	EndObject      TokenType = "}"
	NameSeparator  TokenType = ":"
	ValueSeparator TokenType = ","
	// Values
	False  TokenType = "false"
	Null   TokenType = "null"
	True   TokenType = "true"
	String TokenType = "string"
	Number TokenType = "number"
	Object TokenType = "object"
	Array  TokenType = "array"
	// Numbers
	Minus        TokenType = "-"
	DecimalPoint TokenType = "."
)

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
}

func NewLexer(input string) *Lexer {
	l := &Lexer{input: input}
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
	l.readPosition += 1
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

func (l *Lexer) NextToken() Token {
	var tok Token

	l.skipWhitespace()

	switch l.ch {
	case '[':
		tok = Token{BeginArray, string(l.ch)}
	case ']':
		tok = Token{EndArray, string(l.ch)}
	case '{':
		tok = Token{BeginObject, string(l.ch)}
	case '}':
		tok = Token{EndObject, string(l.ch)}
	case ':':
		tok = Token{NameSeparator, string(l.ch)}
	case ',':
		tok = Token{ValueSeparator, string(l.ch)}
	case '"':
		tok = l.readString()
	case 0:
		tok = Token{EOF, ""}
	default:
		if l.isDigit(l.ch) || l.ch == '-' {
			tok = l.readNumber()
		} else if l.isLiteralName(l.ch) {
			tok = l.readLiteral()
		} else {
			tok = Token{Illegal, string(l.ch)}
		}
	}

	l.readChar()

	return tok
}

func (l *Lexer) readString() Token {
	position := l.position + 1
	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}
	return Token{String, l.input[position:l.position]}
}

func (l *Lexer) readNumber() Token {
	start := l.position
	for l.isDigit(l.peekChar()) || l.peekChar() == '.' || l.peekChar() == '-' {
		l.readChar()
	}
	return Token{Number, l.input[start:l.readPosition]}
}

func (l *Lexer) isLiteralName(ch byte) bool {
	return ch == 't' || ch == 'f' || ch == 'n'
}

func (l *Lexer) readLiteral() Token {
	start := l.position
	var tok Token
	switch l.ch {
	case 't':
		for i, c := range True[1:] {
			if c != rune(l.peekChar()) {
				return Token{Illegal, l.input[start : start+i]}
			}
			l.readChar()
		}
		tok = Token{True, "true"}
	case 'f':
		for i, c := range False[1:] {
			if c != rune(l.peekChar()) {
				return Token{Illegal, l.input[start : start+i]}
			}
			l.readChar()
		}
		tok = Token{False, "false"}
	case 'n':
		for i, c := range Null[1:] {
			if c != rune(l.peekChar()) {
				return Token{Illegal, l.input[start : start+i]}
			}
			l.readChar()
		}
		tok = Token{Null, "null"}
	default:
		return Token{Illegal, string(l.ch)}
	}
	return tok
}

func (l *Lexer) isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func (l *Lexer) skipWhitespace() {
	for {
		if l.ch == ' ' {
			l.readChar()
			continue
		}
		if l.ch == '\\' && (l.peekChar() == 't' || l.peekChar() == 'n' || l.peekChar() == 'r') {
			l.readChar()
			l.readChar()
			continue
		}
		break
	}
}
