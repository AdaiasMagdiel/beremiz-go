package main

import (
	"fmt"
	"os"
	"path"
	"github.com/adaiasmagdiel/beremiz-go/internal/enums"
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

	tokens []Token = [];

	fmt.Println("Program:\n\n")
}
