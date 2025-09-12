package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"path"
	"strings"
	"syscall"

	"github.com/adaiasmagdiel/beremiz-go/internal/err"
	"github.com/adaiasmagdiel/beremiz-go/internal/lexer"
	"github.com/adaiasmagdiel/beremiz-go/internal/parser"
	"github.com/adaiasmagdiel/beremiz-go/internal/pathutils"
)

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		runEval()
		return
	}

	filename := args[0]

	filePath, e := pathutils.ResolveFilePath(filename)
	if e != nil {
		err.Error("Error resolving file path.\n")
		return
	}

	runFile(filePath)
}

func runFile(filepath string) {
	bytes, e := os.ReadFile(filepath)
	if e != nil {
		err.Error("Unable to get the file content.")
		return
	}
	content := string(bytes)

	errorHandler := func() {
		os.Exit(1)
	}

	lexer := lexer.New(content, path.Base(filepath), errorHandler)
	tokens := lexer.Tokenize()

	parser := parser.New(tokens, errorHandler)
	parser.Eval()
}

func runEval() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	reader := bufio.NewReader(os.Stdin)

	go func() {
		<-sigChan
		os.Exit(0)
	}()

	for {
		fmt.Print("\n> ")
		input, e := reader.ReadString('\n')
		if e != nil {
			err.Error("Unable to read stdin.")
			continue
		}

		input = strings.TrimSpace(input)

		if shouldExit(input) {
			break
		}

		if input != "" {
			processInput(input)
		}
	}
}

func shouldExit(input string) bool {
	exitCommands := []string{".exit", "exit"}
	for _, cmd := range exitCommands {
		if strings.EqualFold(input, cmd) {
			return true
		}
	}
	return false
}

func processInput(input string) {
	switch input {
	case ".help":
		printHelp()
		return
	case ".clear":
		clearScreen()
		return
	}

	errorHandler := func() {
	}

	lexer := lexer.New(input, "stdin", errorHandler)
	tokens := lexer.Tokenize()

	parser := parser.New(tokens, errorHandler)
	parser.Eval()
}

func printHelp() {
	fmt.Println(`
Available commands:
  .help           - Show this help message
  .exit, exit     - Exit the program
  .clear          - Clear the screen

  Any other text will be processed normally`)
}

func clearScreen() {
	fmt.Print("\033[H\033[2J")
}
