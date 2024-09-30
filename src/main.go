package main

import (
	"flag"
	"fmt"
	"os"
	"os/user"

	"github.com/ZeroBl21/go-monkey/src/repl"
)

func main() {
	// TODO: Refactor this into a bit flags in REPL struct
	fileFlag := flag.String("file", "", "Path to a file to be evaluated")
	compileFlag := flag.Bool("compile", false, "Enable compilation mode")
	lexerFlag := flag.Bool("lexer", false, "Enable lexer mode to print tokens")
	flag.Parse()

	if *fileFlag != "" {
		data, err := os.ReadFile(*fileFlag)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading file:", err)
			return
		}

		replInstance := repl.New(os.Stdin, os.Stdout)
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

	fmt.Printf("Hello %s! This is the Monkey programming language!\n", user.Username)
	fmt.Printf("Feel free too type in commands\n")

	replInstance := repl.New(os.Stdin, os.Stdout)
	// TODO: Make only only one "Start" function
	if *compileFlag {
		replInstance.StartCompiled()
	} else if *lexerFlag {
		replInstance.StartLexer()
	} else {
		replInstance.Start()
	}
}
