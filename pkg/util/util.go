package util

import (
	"os"
	"unicode"
)

func ReadFile(name string) (string, error) {
	b, err := os.ReadFile(name)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func IsIdentStart(r rune) bool {
	return r == '_' || unicode.IsLetter(r)
}

func IsIdentTail(r rune) bool {
	return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r)
}
