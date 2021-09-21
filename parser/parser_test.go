package parser

import (
	"fmt"
	"github.com/moreal/monkey/ast"
	"github.com/moreal/monkey/lexer"
	"log"
	"testing"
)

func TestLetStatement(t *testing.T) {
	input := `
let x = 5;
let y = 10;
let foobar = 838383;
`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	if program == nil {
		t.Fatalf("ParseProgram() returned nil.")
	}

	if len(program.Statements) != 3 {
		t.Fatalf("It should have 3 statesments but %d", len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
	}
}

func TestReturnStatement(t *testing.T) {
	input := `
return 5;
return 10;
return 100;
`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	if program == nil {
		t.Fatalf("ParseProgram() returned nil.")
	}

	if len(program.Statements) != 3 {
		t.Fatalf("It should have 3 statesments but %d", len(program.Statements))
	}

	for i := 0; i < 3; i++ {
		stmt := program.Statements[i]
		if !testReturnStatement(t, stmt) {
			return
		}
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := `
foobar;
`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	if program == nil {
		t.Fatalf("ParseProgram() returned nil.")
	}

	if len(program.Statements) != 1 {
		t.Fatalf("It should have 1 statesments but %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Expected ExpressionStatement but '%T'", program.Statements[0])
	}

	testIdentifier(t, stmt.Expression, "foobar")
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := `
156497;
`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	if program == nil {
		t.Fatalf("ParseProgram() returned nil.")
	}

	if len(program.Statements) != 1 {
		t.Fatalf("It should have 1 statesments but %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Expected ExpressionStatement but '%T'", program.Statements[0])
	}

	testIntegerLiteral(t, stmt.Expression, 156497)
}

func TestBooleanExpression(t *testing.T) {
	tests := []struct {
		input string
		value bool
	}{
		{"true;", true},
		{"false;", false},
	}

	for _, test := range tests {
		l := lexer.New(test.input)
		p := New(l)

		program := p.ParseProgram()
		if program == nil {
			t.Fatalf("ParseProgram() returned nil.")
		}

		if len(program.Statements) != 1 {
			t.Fatalf("It should have 1 statesment but %d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			log.Fatalf("Expected 'ExpressionStatement' type but '%T'", program.Statements[0])
		}

		testBoolean(t, stmt.Expression, test.value)
	}
}

func TestOrderPrecedences(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"1+1;", "(1 + 1)"},
		{"-1+1;", "((-1) + 1)"},
		{"-1*1 +1;", "(((-1) * 1) + 1)"},
		{"!-1+1*1;", "((!(-1)) + (1 * 1))"},
		{"1+1*(1+2);", "(1 + (1 * (1 + 2)))"},
	}

	for _, test := range tests {
		l := lexer.New(test.input)
		p := New(l)

		program := p.ParseProgram()
		if program == nil {
			t.Fatalf("ParseProgram() returned nil.")
		}

		if len(program.Statements) != 1 {
			t.Fatalf("It should have 1 statement but %d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("Expected 'ExpressionStatement' type but '%T'", program.Statements[0])
		}

		if test.expected != stmt.Expression.String() {
			t.Fatalf("Expected '%s' but '%s'", test.expected, stmt.Expression.String())
		}
	}
}

func TestParsePrefixExpression(t *testing.T) {
	tests := []struct {
		input    string
		operator string
		right    interface{}
	}{
		{"-1;", "-", 1},
		{"!2;", "!", 2},
		{"!true;", "!", true},
		{"!false;", "!", false},
	}

	for _, test := range tests {
		l := lexer.New(test.input)
		p := New(l)

		program := p.ParseProgram()
		if program == nil {
			t.Fatalf("ParseProgram() returned nil.")
		}

		if len(program.Statements) != 1 {
			t.Fatalf("It should have 1 statesment but %d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			log.Fatalf("Expected 'ExpressionStatement' type but '%T'", program.Statements[0])
		}

		testPrefixExpression(t, stmt.Expression, test.operator, test.right)
	}
}

func TestParseInfixExpression(t *testing.T) {
	tests := []struct {
		input    string
		left     interface{}
		operator string
		right    interface{}
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 != 5;", 5, "!=", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 >= 5;", 5, ">=", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 <= 5;", 5, "<=", 5},
		{"foo <= bar;", "foo", "<=", "bar"},
	}

	for _, test := range tests {
		l := lexer.New(test.input)
		p := New(l)

		program := p.ParseProgram()
		if program == nil {
			t.Fatalf("ParseProgram() returned nil.")
		}

		if len(program.Statements) != 1 {
			t.Fatalf("It should have 1 statesment but %d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			log.Fatalf("Expected 'ExpressionStatement' type but '%T'", program.Statements[0])
		}

		testInfixExpression(t, stmt.Expression, test.left, test.operator, test.right)
	}
}

func testIdentifier(t *testing.T, expr ast.Expression, name string) {
	ident, ok := expr.(*ast.Identifier)
	if !ok {
		t.Fatalf("Expected Identifier but '%T'", expr)
	}

	if ident.Value != name {
		t.Fatalf("Expected '%s' identifier but '%s'", name, ident.Value)
	}

	if ident.TokenLiteral() != name {
		t.Fatalf("Expected '%s' TokenLiteral but '%s'", name, ident.TokenLiteral())
	}
}

func testLiteralExpression(t *testing.T, expr ast.Expression, expected interface{}) {
	switch v := expected.(type) {
	case int:
		testIntegerLiteral(t, expr, int64(v))
	case int64:
		testIntegerLiteral(t, expr, v)
	case string:
		testIdentifier(t, expr, v)
	case bool:
		testBoolean(t, expr, v)
	}
}

func testPrefixExpression(t *testing.T, expr ast.Expression, operator string, right interface{}) {
	prefixExpr, ok := expr.(*ast.PrefixExpression)
	if !ok {
		log.Fatalf("Expected 'PrefixExpression' type but '%T'", expr)
	}

	if prefixExpr.Operator != operator {
		log.Fatalf("Expected '%s' type but '%s'", operator, prefixExpr.Operator)
	}

	testLiteralExpression(t, prefixExpr.Right, right)
}

func testInfixExpression(t *testing.T, expr ast.Expression, left interface{}, operator string, right interface{}) {
	infixExpr, ok := expr.(*ast.InfixExpression)
	if !ok {
		t.Fatalf("Expected 'InfixExpression' but '%T'", expr)
	}

	testLiteralExpression(t, infixExpr.Left, left)

	if infixExpr.Operator != operator {
		t.Fatalf("Expected '%s' operator but '%s'", operator, infixExpr.Operator)
	}

	testLiteralExpression(t, infixExpr.Right, right)
}

func testIntegerLiteral(t *testing.T, expression ast.Expression, value int64) {
	integerLiteral, ok := expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("Expected IntegerLiteral but '%T'", expression)
	}

	if integerLiteral.Value != value {
		t.Fatalf("Expected %d identifier but '%d'", value, integerLiteral.Value)
	}

	if integerLiteral.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Fatalf("Expected '%d' TokenLiteral but '%s'", value, integerLiteral.TokenLiteral())
	}
}

func testBoolean(t *testing.T, expression ast.Expression, value bool) {
	boolean, ok := expression.(*ast.Boolean)
	if !ok {
		t.Fatalf("Expected IntegerLiteral but '%T'", expression)
	}

	if boolean.Value != value {
		t.Fatalf("Expected %t identifier but '%t'", value, boolean.Value)
	}

	if boolean.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Fatalf("Expected '%t' TokenLiteral but '%s'", value, boolean.TokenLiteral())
	}
}

func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("TokenLiteral expected 'let' but '%q'", s.TokenLiteral())
		return false
	}

	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("expected 'LetStatement' type but '%T'", s)
		return false
	}

	if letStmt.Name.Value != name {
		t.Errorf("LetStatement.Name.Value expected '%s' but '%s'", name, letStmt.Name.Value)
		return false
	}

	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("LetStatement.Name.TokenLiteral() expected to return '%s' but '%s'", name, letStmt.Name.Value)
		return false
	}

	return true
}

func testReturnStatement(t *testing.T, s ast.Statement) bool {
	if s.TokenLiteral() != "return" {
		t.Errorf("TokenLiteral expected 'let' but '%q'", s.TokenLiteral())
		return false
	}

	_, ok := s.(*ast.ReturnStatement)
	if !ok {
		t.Errorf("expected 'ReturnStatement' type but '%T'", s)
		return false
	}

	return true
}
