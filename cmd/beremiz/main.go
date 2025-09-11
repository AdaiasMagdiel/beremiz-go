package main

import (
	"fmt"
	"os"
	"path"

	"github.com/adaiasmagdiel/beremiz-go/internal/lexer"
)

func main() {
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Unable to get the current directory.")
		return
	}

	file := path.Join(pwd, "main.brz")
	bytes, err := os.ReadFile(file)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Unable to get the current directory.")
		return
	}
	content := string(bytes)

	fmt.Printf("Program:\n    %s\n", content)

	lexer := lexer.New(content, "main.brz", func() {
		os.Exit(1)
	})
	lexer.Tokenize()

}
