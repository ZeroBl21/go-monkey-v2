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
	flag.Parse()

	if *fileFlag != "" {
		data, err := os.ReadFile(*fileFlag)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading file:", err)
			return
		}
		repl.Evaluate(string(data), os.Stdout)

		return
	}

	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Hello %s! This is the Monkey programming language!\n", user.Username)
	fmt.Printf("Feel free too type in commands\n")

	repl.Start(os.Stdin, os.Stdout)
}
