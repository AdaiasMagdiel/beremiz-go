package lexer

import (
	"fmt"

	"github.com/adaiasmagdiel/beremiz-go/internal/enums"
	"github.com/adaiasmagdiel/beremiz-go/internal/error"
)

type Lexer struct {
	content      string
	file         string
	pos          int
	col          int
	line         int
	errorHandler func()
}

func New(content string, file string, errorHandler func()) *Lexer {
	if errorHandler == nil {
		errorHandler = func() {}
	}

	return &Lexer{
		content:      content,
		file:         file,
		pos:          0,
		col:          1,
		line:         1,
		errorHandler: errorHandler,
	}
}

func (l *Lexer) peek() byte {
	return l.content[l.pos]
}

// func (l *Lexer) next() byte {
// 	return l.content[l.pos+1]
// }

func (l *Lexer) consume() byte {
	ch := l.peek()

	l.pos++
	l.col++

	if !l.IsAtEnd() && l.peek() == '\n' {
		l.col = 1
		l.line++
	}

	return ch
}

func (l *Lexer) extractNumber() enums.Token {
	var literal string = ""
	var isInt bool = true

	line := l.line
	col := l.col

	var start int = l.pos
	for {
		if l.IsAtEnd() {
			break
		}

		ch := l.peek()
		if !l.isNum(ch) {
			if ch == '.' {
				if isInt {
					isInt = false
				} else {
					error.LexerError(l.file, enums.Loc{Line: l.line, Col: l.col}, "more than one decimal point in number.")
					l.errorHandler()
				}
			} else {
				break
			}
		}
		l.consume()
	}

	literal = l.content[start:l.pos]
	if isInt {
		return enums.Token{
			Type:    enums.Int,
			Literal: literal,
			Loc:     enums.Loc{Line: line, Col: col},
		}
	} else {
		return enums.Token{
			Type:    enums.Float,
			Literal: literal,
			Loc:     enums.Loc{Line: line, Col: col},
		}
	}
}

func (l *Lexer) isNum(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func (l *Lexer) IsAtEnd() bool {
	return l.pos >= len(l.content)
}

func (l *Lexer) Tokenize() []enums.Token {
	var tokens = []enums.Token{}

	for {
		if l.IsAtEnd() {
			break
		}

		ch := l.peek()

		if l.isNum(ch) {
			// isInt, litInt, litFloat := l.extractNumber()
			fmt.Println(l.extractNumber())
		} else {
			l.consume()
		}
	}

	return tokens
}
