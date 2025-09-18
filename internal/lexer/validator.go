package lexer

func (l *Lexer) isWhitespace(ch byte) bool {
	switch ch {
	case ' ', '\t', '\n', '\r', '\f', '\v':
		return true
	}
	return false
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

func (l *Lexer) isValidHexadecimal(ch byte) bool {
	return l.isNum(ch) || ch >= 'a' && ch <= 'f' || ch >= 'A' && ch <= 'F'
}

func (l *Lexer) isValidOctal(ch byte) bool {
	return ch >= '0' && ch <= '7'
}

func (l *Lexer) isValidBinary(ch byte) bool {
	return ch == '0' || ch == '1'
}
