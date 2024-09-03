//go:build !solution

package spacecollapse

import (
	"strings"
	"unicode/utf8"
)

func CollapseSpaces(input string) string {
	var builder strings.Builder
	builder.Grow(len(input))

	var isPreviousSpace = false
	var runesCnt = utf8.RuneCountInString(input)

	for i := 0; i < runesCnt; i++ {
		r, s := utf8.DecodeRuneInString(input)
		input = input[s:]
		if r == ' ' || r == '\r' || r == '\t' || r == '\n' {
			if isPreviousSpace {
				continue
			} else {
				isPreviousSpace = true
				r = ' '
			}
		} else {
			isPreviousSpace = false
		}
		builder.WriteRune(r)
	}

	return builder.String()
}
