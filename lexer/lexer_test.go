package lexer

import (
	"github.com/moreal/monkey/token"
	"testing"
)

func TestNextToken(t *testing.T) {
	input := `let five = 5;
let ten = 10;
let add = fn(x, y) {
  x + y;
};
let result = add(five, ten);
!-/*5;
5 < 10 > 5;
if (5 < 10) {
  return true;
} else {
  return false;
}
return 1|2 == 2 && 1 >= 2 || 2 <= 3 && 2&0 != 3;
`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.LET, "let"},
		{token.IDENT, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "ten"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "add"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "fn"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.IDENT, "x"},
		{token.PLUS, "+"},
		{token.IDENT, "y"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "result"},
		{token.ASSIGN, "="},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.IDENT, "five"},
		{token.COMMA, ","},
		{token.IDENT, "ten"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.BANG, "!"},
		{token.MINUS, "-"},
		{token.SLASH, "/"},
		{token.ASTERISK, "*"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.INT, "5"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.GT, ">"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		//if (5 < 10) {
		//  return true;
		//} else {
		//  return false;
		//}
		{token.IF, "if"},
		{token.LPAREN, "("},
		{token.INT, "5"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.TRUE, "true"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.ELSE, "else"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.FALSE, "false"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		// return 1|2 == 2 && 1 >= 2 || 2 <= 3 && 2&0 != 3;
		{token.RETURN, "return"},
		{token.INT, "1"},
		{token.BOR, "|"},
		{token.INT, "2"},
		{token.EQ, "=="},
		{token.INT, "2"},
		{token.LAND, "&&"},
		{token.INT, "1"},
		{token.GTE, ">="},
		{token.INT, "2"},
		{token.LOR, "||"},
		{token.INT, "2"},
		{token.LTE, "<="},
		{token.INT, "3"},
		{token.LAND, "&&"},
		{token.INT, "2"},
		{token.BAND, "&"},
		{token.INT, "0"},
		{token.NEQ, "!="},
		{token.INT, "3"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()
		t.Logf("args %q", tok)
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - Token.Type is wrong. (%q != %q) (expected != actual)", i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - Token.Literal is wrong. (%q != %q) (expected != actual)", i, tt.expectedLiteral, tok.Literal)
		}
	}
}
