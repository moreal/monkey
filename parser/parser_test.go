package parser

import (
	"fmt"
	"github.com/moreal/monkey/ast"
	"github.com/moreal/monkey/lexer"
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

	tests := []struct {
	}{
		{},
		{},
		{},
	}

	for i, _ := range tests {
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

	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("Expected Identifier but '%T'", stmt.Expression)
	}

	if ident.Value != "foobar" {
		t.Fatalf("Expected 'foobar' identifier but '%s'", ident.Value)
	}

	if ident.TokenLiteral() != "foobar" {
		t.Fatalf("Expected 'foobar' TokenLiteral but '%s'", ident.TokenLiteral())
	}
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
