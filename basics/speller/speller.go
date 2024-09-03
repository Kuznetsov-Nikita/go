//go:build !solution

package speller

import "strings"

func Spell(n int64) string {
	var ones = []string{"zero", "one", "two", "three", "four",
		"five", "six", "seven", "eight", "nine",
		"ten", "eleven", "twelve", "thirteen", "fourteen",
		"fifteen", "sixteen", "seventeen", "eighteen", "nineteen"}
	var tens = []string{"zero", "ten", "twenty", "thirty", "forty",
		"fifty", "sixty", "seventy", "eighty", "ninety"}

	var result = strings.Builder{}

	if n < 0 {
		result.WriteString("minus ")
		result.WriteString(Spell(-n))
	} else {
		for _, elem := range []struct {
			order int64
			word  string
		}{
			{order: 1000000000000, word: " trillion"},
			{order: 1000000000, word: " billion"},
			{order: 1000000, word: " million"},
			{order: 1000, word: " thousand"},
			{order: 100, word: " hundred"},
		} {
			if n >= elem.order {
				result.WriteString(Spell(n / elem.order))
				result.WriteString(elem.word)
				if n%elem.order != 0 {
					result.WriteString(" ")
					result.WriteString(Spell(n % elem.order))
				}

				return result.String()
			}
		}

		if n < 20 {
			result.WriteString(ones[n])
		} else {
			result.WriteString(tens[n/10])
			if n%10 != 0 {
				result.WriteString("-")
				result.WriteString(Spell(n % 10))
			}
		}
	}

	return result.String()
}
