package error

import (
	"fmt"
	"os"

	"github.com/adaiasmagdiel/beremiz-go/internal/enums"
)

func red(content string) string {
	return fmt.Sprintf("\x1b[31m%s\x1b[0m", content)
}

func LexerError(file string, loc enums.Loc, message string) {
	fmt.Fprintf(os.Stderr, "%s %s", red("LexerError: "), message)
}
