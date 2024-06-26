package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/ZeroBl21/go-monkey/src/evaluator"
	"github.com/ZeroBl21/go-monkey/src/lexer"
	"github.com/ZeroBl21/go-monkey/src/object"
	"github.com/ZeroBl21/go-monkey/src/parser"
)

const (
	RESET  = "\033[0m"
	YELLOW = "\033[33m"
	BLUE   = "\033[34m"
	PROMPT = BLUE + ">> " + RESET
)

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	for {
		fmt.Fprintf(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		EvaluateLine(line, out, env)
	}
}

func Evaluate(input string, out io.Writer) {
	env := object.NewEnvironment()
	EvaluateLine(input, out, env)
}

func EvaluateLine(line string, out io.Writer, env *object.Environment) {
	l := lexer.New(line)
	p := parser.New(l)

	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		printParserErrors(out, p.Errors())
		return
	}

	evaluated := evaluator.Eval(program, env)
	if evaluated != nil {
		io.WriteString(out, YELLOW+evaluated.Inspect()+RESET)
		io.WriteString(out, "\n")
	}
}

func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, "Woops!, We ran into some monkey business here!\n")
	io.WriteString(out, " parser errors:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
