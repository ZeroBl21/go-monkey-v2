package compiler

import (
	"fmt"
	"testing"

	"github.com/ZeroBl21/go-monkey/src/ast"
	"github.com/ZeroBl21/go-monkey/src/code"
	"github.com/ZeroBl21/go-monkey/src/lexer"
	"github.com/ZeroBl21/go-monkey/src/object"
	"github.com/ZeroBl21/go-monkey/src/parser"
)

type compilerTestCase struct {
	input                string
	expectedConstants    []any
	expectedInstructions []code.Instructions
}

func TestIntegerArithmetic(t *testing.T) {
	tests := []compilerTestCase{
		{
			input:             "1 + 2",
			expectedConstants: []any{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
			},
		},
	}

	runCompilerTests(t, tests)
}

func runCompilerTests(t *testing.T, tests []compilerTestCase) {
	t.Helper()

	for _, tt := range tests {
		program := parse(tt.input)

		compiler := New()
		if err := compiler.Compile(program); err != nil {
			t.Fatalf("compiler error: %s", err)
		}

		bytecode := compiler.Bytecode()
		err := testInstructions(tt.expectedInstructions, bytecode.Instructions)
		if err != nil {
			t.Fatalf("testInstructions failed: %s", err)
		}

		err = testConstants(t, tt.expectedConstants, bytecode.Constants)
		if err != nil {
			t.Fatalf("testConstants failed: %s", err)
		}
	}
}

func parse(input string) *ast.Program {
	l := lexer.New(input)
	p := parser.New(l)

	return p.ParseProgram()
}

func testInstructions(expected []code.Instructions, actual code.Instructions) error {
	concatted := concatInstructions(expected)

	if len(actual) != len(concatted) {
		return fmt.Errorf("wrong instruction length.\nwant=%q\ngot =%q",
			concatted, actual)
	}

	for i, ins := range concatted {
		if actual[i] != ins {
			return fmt.Errorf("wrong instruction at %d.\nwant=%q\ngot =%q",
				i, concatted, actual)
		}
	}

	return nil
}

func concatInstructions(s []code.Instructions) code.Instructions {
	out := code.Instructions{}

	for _, ins := range s {
		out = append(out, ins...)
	}

	return out
}

func testConstants(t *testing.T, expected []any, actual []object.Object) error {
	if len(expected) != len(actual) {
		return fmt.Errorf("wrong instruction length.\nwant=%q\ngot=%q",
			len(actual), len(expected))
	}

	for i, constant := range expected {
		switch constant := constant.(type) {
		case int:
			if err := testIntegerObject(int64(constant), actual[i]); err != nil {
				return fmt.Errorf("constant %d - testIntegerObject failed: %s",
					i, err)
			}
		}
	}

	return nil
}

func testIntegerObject(expected int64, actual object.Object) error {
	result, ok := actual.(*object.Integer)
	if !ok {
		return fmt.Errorf("object is not Integer. got=%T (%+v)",
			actual, actual)
	}

	if result.Value != expected {
		return fmt.Errorf("object has wrong value. got=%d, want=%d",
			result.Value, expected)
	}

	return nil
}