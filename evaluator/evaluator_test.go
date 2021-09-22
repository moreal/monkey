package evaluator

import (
	"github.com/moreal/monkey/lexer"
	"github.com/moreal/monkey/object"
	"github.com/moreal/monkey/parser"
	"testing"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"100", 100},
		{"0", 0},
		{"1648590", 1648590},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)

	return Eval(p.ParseProgram())
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) {
	integer, ok := obj.(*object.Integer)
	if !ok {
		t.Fatalf("Expected 'Integer' but '%T'", obj)
	}

	if integer.Value != expected {
		t.Fatalf("Expected %d but %d", expected, integer.Value)
	}
}
