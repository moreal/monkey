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
		{"-5", -5},
		{"--5", 5},
		{"5 * 5 * 5", 125},
		{"-5 * 5 * -5", 125},
		{"-5 * 5 * 5", -125},
		{"-5 * 5 * 5 / 5 + 5", -20},
		{"-5 * 5 * 5 / (5 + 5)", -12},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!!false", false},
		{"!!true", true},
		{"!!5", true},
		{"!5", false},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
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

func testBooleanObject(t *testing.T, obj object.Object, expected bool) {
	boolean, ok := obj.(*object.Boolean)
	if !ok {
		t.Fatalf("Expected 'Boolean' but '%T'", obj)
	}

	if boolean.Value != expected {
		t.Fatalf("Expected %t but %t", expected, boolean.Value)
	}
}
