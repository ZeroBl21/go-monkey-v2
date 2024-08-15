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
	if *compileFlag {
		replInstance.StartCompiled()
	} else {
		replInstance.Start()
	}
}
