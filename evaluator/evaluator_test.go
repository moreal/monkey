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
		{"1 < 2", true},
		{"1 <= 2", true},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 != 2", true},
		{"1 == 2", false},
		{"1 > 2", false},
		{"1 >= 2", false},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"true || false", true},
		{"true && false", false},
		{"true && true", true},
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

func TestIfExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if (1 > 2) { 1 } else { 2 }", 2},
		{"if (1 < 2) { 1 } else { 2 }", 1},
		{"if (1 < 2) {} else {}", nil},
		{"if (1  2) {}", nil},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		if expected, ok := tt.expected.(int); ok {
			testIntegerObject(t, evaluated, int64(expected))
		} else if tt.expected != evaluated {
			t.Fatalf("Expected %+v but %+v", tt.expected, evaluated)
		}
	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{`
if (1 < 10) {
  if (1 < 10) {
    return 10;
  }

  return 1;
}`, 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{"5 + true;", "type mismatch: INTEGER + BOOLEAN"},
		{"5 + true; 5;", "type mismatch: INTEGER + BOOLEAN"},
		{"-true", "unknown operator: -BOOLEAN"},
		{"true + false;", "unknown operator: BOOLEAN + BOOLEAN"},
		{"5; true + false; 5;", "unknown operator: BOOLEAN + BOOLEAN"},
		{"5; 1 && 20; 5;", "unknown operator: INTEGER && INTEGER"},
		{`
if (1 < 10) {
  if (1 < 10) {
    return true + false;
  }

  return 1;
}`, "unknown operator: BOOLEAN + BOOLEAN"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("Expected 'Error' but '%T'", evaluated)
			continue
		}

		if errObj.Message != tt.expectedMessage {
			t.Errorf("Expected error message '%s' but '%s'", tt.expectedMessage, errObj.Message)
		}
	}
}

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let a = 10; a", 10},
		{"let a = 10; a; 5;", 5},
		{"let a = 10; let b = a; b;", 10},
		{"let a = 10; let b = a + 5; b;", 15},
		{`
let a = 5 * 5;
if (1 < 10) {
  if (1 < 10) {
    return a;
  }

  return 1;
}`, 25},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestFunctionObject(t *testing.T) {
	input := "fn(x) { x + 2; };"
	expectedBody := "(x + 2)"

	evaluated := testEval(input)
	fn, ok := evaluated.(*object.Function)
	if !ok {
		t.Fatalf("Expected 'Function' but '%T'", evaluated)
	}

	if len(fn.Parameters) != 1 {
		t.Fatalf("Expected 1 parameters but %d parameters", len(fn.Parameters))
	}

	if fn.Parameters[0].String() != "x" {
		t.Fatalf("The function should have 'x' parameter but '%s'", fn.Parameters[0].String())
	}

	if fn.Body.String() != expectedBody {
		t.Fatalf("Expected body '%s' but '%s'", expectedBody, fn.Body.String())
	}
}

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)

	env := object.NewEnvironment()

	return Eval(p.ParseProgram(), env)
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
