package err

import (
	"fmt"
	"os"
	"strings"

	"github.com/adaiasmagdiel/beremiz-go/internal/tokens"
)

func red(content string) string {
	return fmt.Sprintf("\x1b[31m%s\x1b[0m", content)
}

func Error(message string) {
	fmt.Fprintf(os.Stderr, "%s%s\n", red("Error: "), message)
}

func LexerError(lines []string, loc tokens.Loc, message string, tailLength int) {
	fmt.Fprintf(os.Stderr, "%s %s\n\n", red("LexerError: "), message)

	line := lines[loc.Line-1]
	prefix := fmt.Sprintf("%s:%d:%d: ", loc.File, loc.Line, loc.Col)

	fmt.Printf("\x1b[31m%s\x1b[0m%s\n", prefix, line)

	trail := strings.Repeat(" ", loc.Col-1+len(prefix))
	fmt.Print(trail)
	fmt.Print(red("^"))
	if tailLength > 0 {
		fmt.Print(red(strings.Repeat("~", tailLength)))
	} else if tailLength == -1 {
		length := len(line) - loc.Col - 1
		fmt.Print(red(strings.Repeat("~", length)))
	}
	fmt.Print("\n")
}

func SyntaxError(token tokens.Token, message string) {
	fmt.Fprintf(os.Stderr, "%s%s\n", red("SyntaxError: "), message)
}
