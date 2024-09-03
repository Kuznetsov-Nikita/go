//go:build !solution

package varfmt

import (
	"fmt"
	"strings"
)

func Sprintf(format string, args ...interface{}) string {
	var builder = strings.Builder{}

	var number int32 = -1
	var power int32 = 0
	var arg int32 = 0

	for _, r := range format {
		if r == '{' {
			power = 10
		} else if power != 0 {
			if r == '}' {
				if number == -1 {
					builder.WriteString(fmt.Sprint(args[arg]))
				} else {
					builder.WriteString(fmt.Sprint(args[number]))
					number = -1
				}
				arg++
				power = 0
			} else if r < '0' || r > '9' {
				panic("")
			} else {
				if number == -1 {
					number = 0
				}
				number = number*power + r - '0'
			}
		} else if r == '}' {
			panic("")
		} else {
			builder.WriteRune(r)
		}
	}

	return builder.String()
}
