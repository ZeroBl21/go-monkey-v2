package compiler

import (
	"fmt"
	"sort"

	"github.com/ZeroBl21/go-monkey/src/ast"
	"github.com/ZeroBl21/go-monkey/src/code"
	"github.com/ZeroBl21/go-monkey/src/object"
)

type EmittedInstruction struct {
	Opcode   code.Opcode
	Position int
}

type CompilationScope struct {
	instructions        code.Instructions
	lastInstruction     EmittedInstruction
	previousInstruction EmittedInstruction
}

type Compiler struct {
	constants []object.Object

	symbolTable *SymbolTable

	scopes     []CompilationScope
	scopeIndex int
}

func New() *Compiler {
	mainScope := CompilationScope{
		instructions:        code.Instructions{},
		lastInstruction:     EmittedInstruction{},
		previousInstruction: EmittedInstruction{},
	}

	symbolTable := NewSymbolTable()
	for i, v := range object.Builtins {
		symbolTable.DefineBuiltin(i, v.Name)
	}

	return &Compiler{
		constants:   []object.Object{},
		symbolTable: symbolTable,

		scopes:     []CompilationScope{mainScope},
		scopeIndex: 0,
	}
}

func NewWithState(s *SymbolTable, constants []object.Object) *Compiler {
	compiler := New()
	compiler.symbolTable = s
	compiler.constants = constants

	return compiler
}

func (c *Compiler) Compile(node ast.Node) error {
	switch node := node.(type) {
	case *ast.Program:
		for _, s := range node.Statements {
			if err := c.Compile(s); err != nil {
				return err
			}
		}

	case *ast.Identifier:
		symbol, ok := c.symbolTable.Resolve(node.Value)
		if !ok {
			return fmt.Errorf("undefined variable %s", node.Value)
		}

		c.loadSymbol(symbol)

	// Statements

	case *ast.LetStatement:
		if err := c.Compile(node.Value); err != nil {
			return err
		}
		symbol := c.symbolTable.Define(node.Name.Value)
		if symbol.Scope == GlobalScope {
			c.emit(code.OpSetGlobal, symbol.Index)
		} else {
			c.emit(code.OpSetLocal, symbol.Index)
		}

	case *ast.ExpressionStatement:
		if err := c.Compile(node.Expression); err != nil {
			return err
		}
		c.emit(code.OpPop)

	case *ast.BlockStatement:
		for _, s := range node.Statements {
			if err := c.Compile(s); err != nil {
				return err
			}
		}

	case *ast.ReturnStatement:
		if err := c.Compile(node.ReturnValue); err != nil {
			return err
		}

		c.emit(code.OpReturnValue)

	// Expression

	case *ast.InfixExpression:
		if node.Operator == "<" {
			if err := c.Compile(node.Right); err != nil {
				return err
			}
			if err := c.Compile(node.Left); err != nil {
				return err
			}
			c.emit(code.OpGreaterThan)
			return nil
		}

		if err := c.Compile(node.Left); err != nil {
			return err
		}
		if err := c.Compile(node.Right); err != nil {
			return err
		}
		switch node.Operator {
		case "+":
			c.emit(code.OpAdd)
		case "-":
			c.emit(code.OpSub)
		case "*":
			c.emit(code.OpMul)
		case "/":
			c.emit(code.OpDiv)
		case ">":
			c.emit(code.OpGreaterThan)
		case "==":
			c.emit(code.OpEqual)
		case "!=":
			c.emit(code.OpNotEqual)
		default:
			return fmt.Errorf("unknown operator %s", node.Operator)
		}

	case *ast.PrefixExpression:
		if err := c.Compile(node.Right); err != nil {
			return err
		}

		switch node.Operator {
		case "!":
			c.emit(code.OpBang)
		case "-":
			c.emit(code.OpMinus)
		default:
			return fmt.Errorf("unknown operator %s", node.Operator)
		}

	case *ast.IfExpression:
		if err := c.Compile(node.Condition); err != nil {
			return err
		}

		// Emit an `OpJumpNotTruthy` with a bogus value to patch later.
		jumpNotTruthyPos := c.emit(code.OpJumpNotTruthy, 9999)

		if err := c.Compile(node.Consequence); err != nil {
			return err
		}

		if c.lastInstructionIs(code.OpPop) {
			c.removeLastPop()
		}

		// Emit an `OpJump` with a bogus value to patch later.
		jumpPos := c.emit(code.OpJump, 9999)

		afterConsequncePos := len(c.currentInstructions())
		c.changeOperand(jumpNotTruthyPos, afterConsequncePos)

		if node.Alternative.Statements == nil {
			c.emit(code.OpNull)
		} else {
			if err := c.Compile(node.Alternative); err != nil {
				return err
			}

			if c.lastInstructionIs(code.OpPop) {
				c.removeLastPop()
			}
		}

		afterAlternativePos := len(c.currentInstructions())
		c.changeOperand(jumpPos, afterAlternativePos)

	case *ast.IndexExpression:
		if err := c.Compile(node.Left); err != nil {
			return err
		}

		if err := c.Compile(node.Index); err != nil {
			return err
		}

		c.emit(code.OpIndex)

	// Literals

	case *ast.IntegerLiteral:
		integer := &object.Integer{Value: node.Value}
		c.emit(code.OpConstant, c.addConstant(integer))

	case *ast.StringLiteral:
		string := &object.String{Value: node.Value}
		c.emit(code.OpConstant, c.addConstant(string))

	case *ast.Boolean:
		if node.Value {
			c.emit(code.OpTrue)
		} else {
			c.emit(code.OpFalse)
		}

	case *ast.ArrayLiteral:
		for _, el := range node.Elements {
			if err := c.Compile(el); err != nil {
				return err
			}
		}
		c.emit(code.OpArray, len(node.Elements))

	case *ast.HashLiteral:
		keys := []ast.Expression{}
		for k := range node.Pairs {
			keys = append(keys, k)
		}

		sort.Slice(keys, func(i, j int) bool {
			return keys[i].String() < keys[j].String()
		})

		for _, key := range keys {
			if err := c.Compile(key); err != nil {
				return err
			}

			if err := c.Compile(node.Pairs[key]); err != nil {
				return err
			}
		}
		c.emit(code.OpHash, len(node.Pairs)*2)

	case *ast.FunctionLiteral:
		c.enterScope()

		for _, parameter := range node.Parameters {
			c.symbolTable.Define(parameter.Value)
		}

		if err := c.Compile(node.Body); err != nil {
			return err
		}

		if c.lastInstructionIs(code.OpPop) {
			c.replaceLastPopWithReturn()
		}

		if !c.lastInstructionIs(code.OpReturnValue) {
			c.emit(code.OpReturn)
		}

		numLocals := c.symbolTable.numDefinitions
		instructions := c.leaveScope()

		compiledFn := &object.CompiledFunction{
			Instructions:  instructions,
			NumLocals:     numLocals,
			NumParameters: len(node.Parameters),
		}

		c.emit(code.OpConstant, c.addConstant(compiledFn))

	case *ast.CallExpression:
		if err := c.Compile(node.Function); err != nil {
			return err
		}

		for _, arg := range node.Arguments {
			if err := c.Compile(arg); err != nil {
				return err
			}
		}

		c.emit(code.OpCall, len(node.Arguments))

	}

	return nil
}

func (c *Compiler) addConstant(obj object.Object) int {
	c.constants = append(c.constants, obj)
	return len(c.constants) - 1
}

func (c *Compiler) emit(op code.Opcode, operands ...int) int {
	ins := code.Make(op, operands...)
	pos := c.addInstruction(ins)

	c.setLastInstruction(op, pos)

	return pos
}

func (c *Compiler) currentInstructions() code.Instructions {
	return c.scopes[c.scopeIndex].instructions
}

func (c *Compiler) addInstruction(ins []byte) int {
	posNewInstruction := len(c.currentInstructions())
	updatedInstructions := append(c.currentInstructions(), ins...)

	c.scopes[c.scopeIndex].instructions = updatedInstructions

	return posNewInstruction
}

func (c *Compiler) setLastInstruction(op code.Opcode, pos int) {
	previous := c.scopes[c.scopeIndex].lastInstruction
	last := EmittedInstruction{Opcode: op, Position: pos}

	c.scopes[c.scopeIndex].previousInstruction = previous
	c.scopes[c.scopeIndex].lastInstruction = last
}

func (c *Compiler) changeOperand(opPos int, operand int) {
	op := code.Opcode(c.currentInstructions()[opPos])
	newInstruction := code.Make(op, operand)

	c.replaceInstruction(opPos, newInstruction)
}

func (c *Compiler) replaceInstruction(
	pos int,
	newInstruction code.Instructions,
) {
	ins := c.currentInstructions()

	for i := 0; i < len(newInstruction); i++ {
		ins[pos+i] = newInstruction[i]
	}
}

func (c *Compiler) replaceLastPopWithReturn() {
	lastPop := c.scopes[c.scopeIndex].lastInstruction.Position
	c.replaceInstruction(lastPop, code.Make(code.OpReturnValue))

	c.scopes[c.scopeIndex].lastInstruction.Opcode = code.OpReturnValue
}

func (c *Compiler) lastInstructionIs(op code.Opcode) bool {
	if len(c.currentInstructions()) == 00 {
		return false
	}

	return c.scopes[c.scopeIndex].lastInstruction.Opcode == op
}

func (c *Compiler) removeLastPop() {
	last := c.scopes[c.scopeIndex].lastInstruction
	previous := c.scopes[c.scopeIndex].previousInstruction

	old := c.currentInstructions()
	new := old[:last.Position]

	c.scopes[c.scopeIndex].instructions = new
	c.scopes[c.scopeIndex].lastInstruction = previous
}

func (c *Compiler) loadSymbol(s Symbol) {
	switch s.Scope {
	case GlobalScope:
		c.emit(code.OpGetGlobal, s.Index)
	case LocalScope:
		c.emit(code.OpGetLocal, s.Index)
	case BuilinScope:
		c.emit(code.OpGetBuilin, s.Index)
	default:
		panic(fmt.Sprintf("unexpected compiler.SymbolScope: %#v", s.Scope))
	}
}

func (c *Compiler) enterScope() {
	scope := CompilationScope{
		instructions:        code.Instructions{},
		lastInstruction:     EmittedInstruction{},
		previousInstruction: EmittedInstruction{},
	}

	c.scopes = append(c.scopes, scope)
	c.scopeIndex++

	c.symbolTable = NewEnclosedSymbolTable(c.symbolTable)
}

func (c *Compiler) leaveScope() code.Instructions {
	instructions := c.currentInstructions()

	c.scopes = c.scopes[:len(c.scopes)-1]
	c.scopeIndex--

	c.symbolTable = c.symbolTable.Outer

	return instructions
}

func (c *Compiler) Bytecode() *Bytecode {
	return &Bytecode{
		Instructions: c.currentInstructions(),
		Constants:    c.constants,
	}
}

type Bytecode struct {
	Instructions code.Instructions
	Constants    []object.Object
}
