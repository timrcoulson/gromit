package printing

import (
	"strings"
	"unicode"
)

func Clean(input string) string  {
	return strings.Map(func(r rune) rune {
		if r <= unicode.MaxASCII{
			return r
		}
		return -1
	}, input)
}
