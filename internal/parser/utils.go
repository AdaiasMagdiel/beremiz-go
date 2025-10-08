package parser

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/adaiasmagdiel/beremiz-go/internal/tokens"
)

type BlockType uint8

const (
	BlockNone BlockType = iota
	BlockIf
	BlockFor
	BlockDefine
)

func Pop[T any](s []T) ([]T, T, error) {
	var zero T
	if len(s) == 0 {
		return s, zero, errors.New("cannot pop from empty slice")
	}
	last := s[len(s)-1]
	s = s[:len(s)-1]
	return s, last, nil
}

func toBool(t tokens.Token) bool {
	switch t.Type {
	case tokens.Bool:
		return t.Literal.(bool)
	case tokens.Nil:
		return false
	case tokens.Int:
		return t.Literal != 0
	case tokens.Float:
		return t.Literal != 0.0
	case tokens.String:
		return t.Literal != ""
	default:
		return false
	}
}

func toString(t tokens.Token) string {
	switch v := t.Literal.(type) {
	case nil:
		return "nil"
	case bool:
		if v {
			return "true"
		}
		return "false"
	case int64:
		return strconv.FormatInt(v, 10)
	case float64:
		return strconv.FormatFloat(v, 'g', -1, 64)
	case string:
		return v
	default:
		return fmt.Sprintf("%v", v)
	}
}
