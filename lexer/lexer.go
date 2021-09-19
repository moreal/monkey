package lexer

import "github.com/moreal/monkey/token"

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) NextToken() (tok token.Token) {
	l.skipWhitespaces()

	switch l.ch {
	case '=':
		tok = newTokenWithChar(token.ASSIGN, l.ch)
	case '+':
		tok = newTokenWithChar(token.PLUS, l.ch)
	case '(':
		tok = newTokenWithChar(token.LPAREN, l.ch)
	case ')':
		tok = newTokenWithChar(token.RPAREN, l.ch)
	case '{':
		tok = newTokenWithChar(token.LBRACE, l.ch)
	case '}':
		tok = newTokenWithChar(token.RBRACE, l.ch)
	case ',':
		tok = newTokenWithChar(token.COMMA, l.ch)
	case ';':
		tok = newTokenWithChar(token.SEMICOLON, l.ch)
	case 0:
		tok = newToken(token.EOF, "")
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = identifierToTokenType(tok.Literal)
			return
		} else if isDigits(l.ch) {
			tok = newToken(token.INT, l.readInteger())
			return
		} else {
			tok = newTokenWithChar(token.ILLEGAL, l.ch)
		}
	}

	l.readChar()
	return
}

func identifierToTokenType(identifier string) token.TokenType {
	keywords := map[string]token.TokenType{
		"fn":  token.FUNCTION,
		"let": token.LET,
	}

	if tok, ok := keywords[identifier]; ok {
		return tok
	}

	return token.IDENT
}

func (l *Lexer) skipWhitespaces() {
	for isWhitespace(l.ch) {
		l.readChar()
	}
}

func isWhitespace(char byte) bool {
	return char == ' ' || char == '\n' || char == '\t' || char == '\r'
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}

	return l.input[position:l.position]
}

func (l *Lexer) readInteger() string {
	position := l.position
	for isDigits(l.ch) {
		l.readChar()
	}

	return l.input[position:l.position]
}

func isLetter(char byte) bool {
	return 'a' <= char && char <= 'z' || 'A' <= char && char <= 'Z' || char == '_'
}

func isDigits(char byte) bool {
	return '0' <= char && char <= '9'
}

func newTokenWithChar(tokenType token.TokenType, char byte) token.Token {
	return newToken(tokenType, string(char))
}

func newToken(tokenType token.TokenType, literal string) token.Token {
	return token.Token{Type: tokenType, Literal: literal}
}

func (l *Lexer) readChar() {
	// TODO: Support Unicode
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}
