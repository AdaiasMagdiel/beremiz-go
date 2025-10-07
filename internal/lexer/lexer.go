package lexer

import (
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

func (l *Lexer) prev() byte {
	if l.pos-1 >= 0 {
		return l.content[l.pos-1]
	} else {
		return 0
	}
}

func (l *Lexer) next() byte {
	if l.pos+1 >= len(l.content) {
		return 0
	} else {
		return l.content[l.pos+1]
	}
}

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

func (l *Lexer) isAtEnd() bool {
	return l.pos >= len(l.content)
}

func (l *Lexer) GetLines() []string {
	return l.lines
}

func (l *Lexer) Tokenize() []tokens.Token {
	var ts = []tokens.Token{}

	for {
		if l.isAtEnd() {
			break
		}

		ch := l.peek()

		if l.isNum(ch) ||
			ch == '.' && l.isNum(l.next()) ||
			ch == '-' && l.isNum(l.next()) ||
			ch == '+' && l.isNum(l.next()) {
			token := l.extractNumber()
			ts = append(ts, token)
		} else if tokens.IsOperator(ch) {
			if ch == '*' && l.next() == '*' {
				ts = append(ts, tokens.Token{
					Type:    tokens.Exp,
					Literal: "**",
					Loc:     l.getLoc(),
				})
				l.consume()
				l.consume()

			} else if ch == '<' && l.next() == '=' {
				ts = append(ts, tokens.Token{
					Type:    tokens.Le,
					Literal: "<=",
					Loc:     l.getLoc(),
				})
				l.consume()
				l.consume()
			} else if ch == '>' && l.next() == '=' {
				ts = append(ts, tokens.Token{
					Type:    tokens.Ge,
					Literal: ">=",
					Loc:     l.getLoc(),
				})
				l.consume()
				l.consume()
			} else {
				ts = append(ts, tokens.Token{
					Type:    tokens.Operators[string(ch)],
					Literal: string(l.consume()),
					Loc:     l.getLoc(),
				})
			}
		} else if l.isAlpha(ch) || ch == '_' {
			token := l.extractIdentifier()
			ts = append(ts, token)
		} else if ch == '#' {
			l.extractComment()
		} else if ch == '\'' || ch == '"' {
			token := l.extractString()
			ts = append(ts, token)
		} else if l.isWhitespace(ch) {
			l.consume()
		} else {
			err.LexerError(l.lines, l.getLoc(), "invalid character '"+string(ch)+"'", 0)
			l.errorHandler()
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
