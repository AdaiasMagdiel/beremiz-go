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
	fmt.Fprintf(os.Stderr, "%s%s\n", red("LexerError: "), message)

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

func SyntaxError(token tokens.Token, message string, lines []string) {
	fmt.Fprintf(os.Stderr, "%s%s\n\n", red("SyntaxError: "), message)

	line := lines[token.Loc.Line-1]
	fmt.Println(token.Loc.File + ":\n")
	prefix := fmt.Sprintf("%d:%d | ", token.Loc.Line, token.Loc.Col)

	fmt.Printf("\x1b[31m%s\x1b[0m%s\n", prefix, line)

	trail := strings.Repeat(" ", token.Loc.Col-1+len(prefix)-1)
	fmt.Print(trail)

	lit, ok := token.Literal.(string)
	if !ok {
		lit = fmt.Sprint(token.Literal)
	}
	fmt.Print(red("^"))
	fmt.Print(red(strings.Repeat("~", len(lit)-1)))

	fmt.Print("\n")
}
