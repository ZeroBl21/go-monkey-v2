package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/ZeroBl21/go-monkey/src/compiler"
	"github.com/ZeroBl21/go-monkey/src/evaluator"
	"github.com/ZeroBl21/go-monkey/src/lexer"
	"github.com/ZeroBl21/go-monkey/src/object"
	"github.com/ZeroBl21/go-monkey/src/parser"
	"github.com/ZeroBl21/go-monkey/src/token"
	"github.com/ZeroBl21/go-monkey/src/vm"
)

const (
	RESET  = "\033[0m"
	YELLOW = "\033[33m"
	BLUE   = "\033[34m"
	PROMPT = ">> "
)

type REPL struct {
	env     *object.Environment
	scanner *bufio.Scanner
	out     io.Writer

	constants   []object.Object
	globals     []object.Object
	symbolTable *compiler.SymbolTable
}

func New(in io.Reader, out io.Writer) *REPL {
	return &REPL{
		env:     object.NewEnvironment(),
		scanner: bufio.NewScanner(in),
		out:     out,

		// Compiler
		constants:   []object.Object{},
		symbolTable: compiler.NewSymbolTable(),

		// VM
		globals: make([]object.Object, vm.GlobalSize),
	}
}

func (r *REPL) Start() {
	for {
		fmt.Fprint(r.out, applyColor(BLUE, PROMPT))
		scanned := r.scanner.Scan()
		if !scanned {
			return
		}

		line := r.scanner.Text()
		r.EvaluateLine(line)
	}
}

func (r *REPL) EvaluateLine(line string) {
	l := lexer.New(line)
	p := parser.New(l)

	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		r.printParserErrors(p.Errors())
		return
	}

	evaluated := evaluator.Eval(program, r.env)
	if evaluated != nil {
		io.WriteString(r.out, applyColor(YELLOW, evaluated.Inspect()))
		io.WriteString(r.out, "\n")
	}
}

func (r *REPL) StartCompiled() {
	for {
		fmt.Fprint(r.out, applyColor(BLUE, PROMPT))
		scanned := r.scanner.Scan()
		if !scanned {
			return
		}

		line := r.scanner.Text()
		r.EvaluateLineCompiled(line)
	}
}

func (r *REPL) StartLexer() {
	for {
		fmt.Fprint(r.out, applyColor(BLUE, PROMPT))
		scanned := r.scanner.Scan()
		if !scanned {
			return
		}

		line := r.scanner.Text()
		r.PrintTokens(line)
	}
}

func (r *REPL) EvaluateLineCompiled(line string) {
	l := lexer.New(line)
	p := parser.New(l)

	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		r.printParserErrors(p.Errors())
		return
	}

	comp := compiler.NewWithState(r.symbolTable, r.constants)
	if err := comp.Compile(program); err != nil {
		fmt.Fprintf(r.out, "Woops! Compilation failed:\n %s\n", err)
	}

	code := comp.Bytecode()
	r.constants = code.Constants

	machine := vm.NewWithGlobalStore(code, r.globals)
	if err := machine.Run(); err != nil {
		fmt.Fprintf(r.out, "Woops! Executing bytecode failed:\n %s\n", err)
	}

	lastPopped := machine.LastPoppedStackElem()
	io.WriteString(r.out, applyColor(YELLOW, lastPopped.Inspect()))
	io.WriteString(r.out, "\n")
}

func (r *REPL) printParserErrors(errors []string) {
	io.WriteString(r.out, "Woops!, We ran into some monkey business here!\n")
	io.WriteString(r.out, " parser errors:\n")
	for _, msg := range errors {
		io.WriteString(r.out, "\t"+msg+"\n")
	}
}

func (r *REPL) PrintTokens(line string) {
	l := lexer.New(line)

	for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
		fmt.Fprintf(r.out, "%+v\n", tok)
	}
}

func applyColor(color, text string) string {
	return color + text + RESET
}
