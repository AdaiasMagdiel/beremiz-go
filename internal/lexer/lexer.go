package lexer

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/adaiasmagdiel/beremiz-go/internal/err"
	"github.com/adaiasmagdiel/beremiz-go/internal/tokens"
)

type Lexer struct {
	content      string
	file         string
	lines        []string
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
		lines:        strings.Split(content, "\n"),
		pos:          0,
		col:          1,
		line:         1,
		errorHandler: errorHandler,
	}
}

func (l *Lexer) getLoc() tokens.Loc {
	return tokens.Loc{
		File: l.file,
		Line: l.line,
		Col:  l.col,
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

	if !l.isAtEnd() && l.peek() == '\n' {
		l.col = 1
		l.line++
	}

	return ch
}

func (l *Lexer) extractNumber() tokens.Token {
	var literal string = ""
	var isInt bool = true

	line := l.line
	col := l.col

	var start int = l.pos
	for {
		if l.isAtEnd() {
			break
		}

		ch := l.peek()
		if !l.isNum(ch) {
			if ch == '.' {
				if isInt {
					isInt = false
				} else {
					err.LexerError(l.lines, l.getLoc(), "more than one decimal point in number.", -1)
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
		n, e := strconv.ParseInt(literal, 10, 64)
		if e != nil {
			n = 0.0
			err.Error(fmt.Sprintf("Unable to convert literal '%s' to int64.", literal))
			l.errorHandler()
		}
		return tokens.Token{
			Type:    tokens.Int,
			Literal: n,
			Loc:     tokens.Loc{File: l.file, Line: line, Col: col},
		}
	} else {
		n, e := strconv.ParseFloat(literal, 64)
		if e != nil {
			n = 0.0
			err.Error(fmt.Sprintf("Unable to convert literal '%s' to float64.", literal))
			l.errorHandler()
		}
		return tokens.Token{
			Type:    tokens.Float,
			Literal: n,
			Loc:     tokens.Loc{File: l.file, Line: line, Col: col},
		}
	}
}

func (l *Lexer) extractIdentifier() tokens.Token {
	var literal string = ""

	line := l.line
	col := l.col

	var start int = l.pos
	for {
		if l.isAtEnd() {
			break
		}

		ch := l.peek()
		if !l.isValidIdentifier(ch) {
			break
		}

		l.consume()
	}

	literal = l.content[start:l.pos]
	kind := tokens.Identifier

	if tokens.IsKeyword(literal) {
		kind = tokens.Keywords[literal]
	}

	return tokens.Token{
		Type:    kind,
		Literal: literal,
		Loc:     tokens.Loc{File: l.file, Line: line, Col: col},
	}
}

func (l *Lexer) isNum(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func (l *Lexer) isAlpha(ch byte) bool {
	return ch >= 'a' && ch <= 'z' || ch >= 'A' && ch <= 'Z'
}

func (l *Lexer) isAlphaNum(ch byte) bool {
	return l.isAlpha(ch) || l.isNum(ch)
}

func (l *Lexer) isValidIdentifier(ch byte) bool {
	return l.isAlphaNum(ch) || ch == '_'
}

func (l *Lexer) isAtEnd() bool {
	return l.pos >= len(l.content)
}

func (l *Lexer) Tokenize() []tokens.Token {
	var ts = []tokens.Token{}

	for {
		if l.isAtEnd() {
			break
		}

		ch := l.peek()

		if l.isNum(ch) {
			token := l.extractNumber()
			ts = append(ts, token)
		} else if tokens.IsOperator(ch) {
			ts = append(ts, tokens.Token{
				Type:    tokens.Operators[string(ch)],
				Literal: string(l.consume()),
				Loc:     l.getLoc(),
			})
		} else if l.isAlpha(ch) || ch == '_' {
			token := l.extractIdentifier()
			ts = append(ts, token)
		} else {
			l.consume()
		}
	}

	ts = append(ts, tokens.Token{
		Type:    tokens.EOF,
		Literal: "EOF",
		Loc:     l.getLoc(),
	})

	return ts
}
