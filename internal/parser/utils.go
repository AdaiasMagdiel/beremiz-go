package parser

import "errors"

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
