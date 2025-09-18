package lexer

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/adaiasmagdiel/beremiz-go/internal/err"
	"github.com/adaiasmagdiel/beremiz-go/internal/tokens"
)

func (l *Lexer) extractNumber() tokens.Token {
	var literal string = ""
	var isInt bool = true
	var isNegative bool = false
	var isHex bool = false
	var isOctal bool = false
	var isBinary bool = false

	line := l.line
	col := l.col

	if l.peek() == '+' {
		l.consume()
	} else if l.peek() == '-' {
		isNegative = true
		l.consume()
	} else if l.peek() == '.' {
		isInt = false
		literal += "0."
		l.consume()
	} else if l.peek() == '0' && l.isNum(l.next()) {
		isOctal = true
	}

loop:
	for {
		if l.isAtEnd() {
			break
		}

		ch := l.peek()

		switch {
		case l.isWhitespace(ch):
			break loop

		case ch == '_':
			l.consume()
			continue

		case ch == '.' && !isHex && !isOctal && !isBinary:
			if isInt {
				isInt = false
			} else {
				err.LexerError(l.lines, tokens.Loc{File: l.file, Line: l.line, Col: max(l.col-1, 0)},
					"More than one decimal point in number.", 0)
				l.errorHandler()
				break loop
			}

		case ch == 'x' || ch == 'X':
			if isHex || len(literal) != 1 || literal[0] != '0' {
				err.LexerError(l.lines, tokens.Loc{File: l.file, Line: l.line, Col: max(l.col-1, 0)},
					fmt.Sprintf("Invalid hexadecimal literal: expected '0' before '%c', but found '%c'.", ch, l.prev()), 0)
				l.errorHandler()
				break loop
			}
			isHex = true

		case ch == 'o' || ch == 'O':
			if isOctal || len(literal) != 1 || literal[0] != '0' {
				err.LexerError(l.lines, tokens.Loc{File: l.file, Line: l.line, Col: max(l.col-1, 0)},
					fmt.Sprintf("Invalid octal literal: expected '0' before '%c', but found '%c'.", ch, l.prev()), 0)
				l.errorHandler()
				break loop
			}
			isOctal = true

		case !isHex && (ch == 'b' || ch == 'B'):
			if isBinary || len(literal) != 1 || literal[0] != '0' {
				err.LexerError(l.lines, tokens.Loc{File: l.file, Line: l.line, Col: max(l.col-1, 0)},
					fmt.Sprintf("Invalid binary literal: expected '0' before '%c', but found '%c'.", ch, l.prev()), 0)
				l.errorHandler()
				break loop
			}
			isBinary = true

		case isHex && !l.isValidHexadecimal(ch) && l.isAlpha(ch):
			err.LexerError(l.lines, tokens.Loc{File: l.file, Line: l.line, Col: max(l.col-1, 0)},
				fmt.Sprintf("Invalid character '%c' in hexadecimal literal.", ch), 0)
			l.errorHandler()
			break loop

		case isOctal && !l.isValidOctal(ch) && l.isAlpha(ch):
			err.LexerError(l.lines, tokens.Loc{File: l.file, Line: l.line, Col: max(l.col-1, 0)},
				fmt.Sprintf("Invalid character '%c' in octal literal.", ch), 0)
			l.errorHandler()
			break loop

		case isBinary && !l.isValidBinary(ch) && l.isAlpha(ch):
			err.LexerError(l.lines, tokens.Loc{File: l.file, Line: l.line, Col: max(l.col-1, 0)},
				fmt.Sprintf("Invalid character '%c' in binary literal.", ch), 0)
			l.errorHandler()
			break loop

		case !l.isNum(ch) && !isHex && !isOctal:
			break loop
		}

		literal += strings.ToLower(string(l.consume()))
	}

	if isInt {
		var base int = 10
		switch {
		case isHex:
			literal = literal[2:]
			base = 16
		case isOctal:
			if len(literal) > 1 && literal[1] == 'o' {
				literal = literal[2:]
			}
			base = 8
		case isBinary:
			literal = literal[2:]
			base = 2
		}

		if isNegative {
			literal = "-" + literal
		}

		n, e := strconv.ParseInt(literal, base, 64)
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
		if isNegative {
			n *= -1
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
	} else if literal == "true" {
		kind = tokens.True
	} else if literal == "false" {
		kind = tokens.False
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
