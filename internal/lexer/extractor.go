package lexer

import (
	"fmt"
	"strconv"

	"github.com/adaiasmagdiel/beremiz-go/internal/err"
	"github.com/adaiasmagdiel/beremiz-go/internal/tokens"
)

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

func (l *Lexer) extractComment() tokens.Token {
	l.consume() // Remove #

	isMultiline := false
	if l.peek() == '[' {
		isMultiline = true
	}

	for {
		if l.isAtEnd() {
			break
		}

		ch := l.peek()

		if ch == '\n' {
			if !isMultiline {
				break
			}
		}

		if ch == '#' {
			if isMultiline {
				l.consume()
				break
			}
		}

		l.consume()
	}

	return tokens.Token{}
}

func (l *Lexer) extractString() tokens.Token {
	var stringType byte = l.consume() // Remove ' or "
	var literal string = ""

	line := l.line
	col := l.col

	var start int = l.pos
	for {
		if l.isAtEnd() {
			err.LexerError(l.lines, tokens.Loc{
				File: l.file,
				Line: line,
				Col:  col - 1,
			}, "Unterminated string literal.", l.col-col)
			l.errorHandler()
			break
		}

		ch := l.peek()

		// Disallow multiline strings here
		if ch == '\n' {
			err.LexerError(l.lines, tokens.Loc{
				File: l.file,
				Line: line,
				Col:  col - 1,
			}, "Unterminated string literal.", l.col-col)
			l.errorHandler()
			break
		}

		if ch == stringType {
			if l.prev() != '\\' {
				l.consume() // Remove ' or "
				break
			}
		}

		l.consume()
	}

	literal = l.content[start : l.pos-1]

	parsed, e := strconv.Unquote(`"` + literal + `"`)
	if e != nil {
		err.Error("Unable to parse string literal.")
		l.errorHandler()
	}

	return tokens.Token{
		Type:    tokens.String,
		Literal: parsed,
		Loc:     tokens.Loc{File: l.file, Line: line, Col: col},
	}
}
