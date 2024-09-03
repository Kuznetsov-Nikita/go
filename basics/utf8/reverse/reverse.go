//go:build !solution

package reverse

import (
	"strings"
	"unicode/utf8"
)

func Reverse(input string) string {
	var builder strings.Builder
	builder.Grow(len(input))

	for i := 0; i < utf8.RuneCountInString(input); i++ {
		r, s := utf8.DecodeLastRuneInString(input)
		input = input[:len(input)-s]
		builder.WriteRune(r)
	}

	return builder.String()
}
