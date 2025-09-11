package main

import (
	"fmt"
	"os"
	"path"

	"github.com/adaiasmagdiel/beremiz-go/internal/lexer"
	"github.com/adaiasmagdiel/beremiz-go/internal/parser"
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

	errorHandler := func() {
		os.Exit(1)
	}

	lexer := lexer.New(content, "main.brz", errorHandler)
	tokens := lexer.Tokenize()

	// fmt.Println(tokens)

	parser := parser.New(tokens, errorHandler)
	parser.Eval()
}
