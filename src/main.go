package main

import (
	"flag"
	"fmt"
	"os"
	"os/user"

	"github.com/ZeroBl21/go-monkey/src/repl"
)

func main() {
	fileFlag := flag.String("file", "", "Path to a file to be evaluated")
	compileFlag := flag.Bool("compile", false, "Enable compilation mode")
	lexerFlag := flag.Bool("lexer", false, "Enable lexer mode to print tokens")
	precedenceFlag := flag.Bool(
		"precedence",
		false,
		"Enable precedence mode to show parsed program",
	)
	flag.Parse()

	replFlagsMap := map[bool]int{
		*compileFlag:    repl.CompileFlag,
		*lexerFlag:      repl.LexerFlag,
		*precedenceFlag: repl.PrecedenceFlag,
	}

	flags := 0
	for condition, flag := range replFlagsMap {
		if condition {
			flags |= flag
		}
	}

	replInstance := repl.New(os.Stdin, os.Stdout)
	replInstance.SetFlags(flags)

	if *fileFlag != "" {
		data, err := os.ReadFile(*fileFlag)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading file:", err)
			return
		}

		if *compileFlag {
			replInstance.EvaluateLineCompiled(string(data))
		} else {
			replInstance.EvaluateLine(string(data))
		}

		return
	}

	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	fmt.Printf(
		"Hello %s! This is the Monkey programming language!\n",
		user.Username,
	)
	fmt.Printf("Feel free too type in commands\n")

	replInstance.Start()
}
