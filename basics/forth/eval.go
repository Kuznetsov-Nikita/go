//go:build !solution

package main

import (
	"errors"
	"strconv"
	"strings"
)

type Evaluator struct {
	words map[string][]string
	stack []int
}

// NewEvaluator creates evaluator.
func NewEvaluator() *Evaluator {
	return &Evaluator{
		words: make(map[string][]string),
		stack: make([]int, 0),
	}
}

// Process evaluates sequence of words or definition.
//
// Returns resulting stack state and an error.
func (e *Evaluator) Process(row string) ([]int, error) {
	var tiles = strings.Split(row, " ")

	if len(tiles) == 0 {
		return e.stack, nil
	}

	var i = 0
	for i < len(tiles) {
		if tiles[i] == ":" {
			i++
			var definingWord = strings.ToLower(tiles[i])

			_, err := strconv.Atoi(definingWord)
			if err == nil {
				return e.stack, errors.New("number redefinition")
			}

			prevDefinition, wasInMap := e.words[definingWord]
			e.words[definingWord] = make([]string, 0)

			i++
			for i < len(tiles) && tiles[i] != ";" {
				var tile = strings.ToLower(tiles[i])
				commands, ok := e.words[tile]
				if ok == false {
					if tile != "+" && tile != "-" && tile != "*" && tile != "/" &&
						tile != "dup" && tile != "over" && tile != "drop" && tile != "swap" {
						_, err := strconv.Atoi(tile)
						if err != nil {
							return e.stack, errors.New(tile + ": unknown word")
						}
					}
					e.words[definingWord] = append(e.words[definingWord], tile)
				} else {
					if wasInMap && tile == definingWord {
						e.words[definingWord] = append(e.words[definingWord], prevDefinition...)
					} else {
						e.words[definingWord] = append(e.words[definingWord], commands...)
					}
				}
				i++
			}

			if len(e.words[definingWord]) == 0 || tiles[i] != ";" {
				return e.stack, nil
			}
		} else {
			var tile = strings.ToLower(tiles[i])
			commands, ok := e.words[tile]
			if ok == false {
				if tile == "+" || tile == "-" || tile == "*" || tile == "/" {
					if len(e.stack) < 2 {
						return e.stack, errors.New("incorrect operation")
					}

					var first = e.stack[len(e.stack)-1]
					var second = e.stack[len(e.stack)-2]
					e.stack = e.stack[:len(e.stack)-2]

					var result int
					switch tile {
					case "+":
						result = second + first
					case "-":
						result = second - first
					case "*":
						result = second * first
					case "/":
						if first == 0 {
							return e.stack, errors.New("division by zero")
						}
						result = second / first
					}
					e.stack = append(e.stack, result)
				} else if tile == "dup" {
					if len(e.stack) == 0 {
						return e.stack, errors.New("incorrect operation")
					}

					e.stack = append(e.stack, e.stack[len(e.stack)-1])
				} else if tile == "over" {
					if len(e.stack) < 2 {
						return e.stack, errors.New("incorrect operation")
					}

					e.stack = append(e.stack, e.stack[len(e.stack)-2])
				} else if tile == "drop" {
					if len(e.stack) == 0 {
						return e.stack, errors.New("incorrect operation")
					}

					e.stack = e.stack[:len(e.stack)-1]
				} else if tile == "swap" {
					if len(e.stack) < 2 {
						return e.stack, errors.New("incorrect operation")
					}

					var first = e.stack[len(e.stack)-1]
					var second = e.stack[len(e.stack)-2]
					e.stack = e.stack[:len(e.stack)-2]
					e.stack = append(e.stack, first, second)
				} else {
					number, err := strconv.Atoi(tile)
					if err != nil {
						return e.stack, errors.New(tile + ": unknown word")
					} else {
						e.stack = append(e.stack, number)
					}
				}
			} else {
				for _, command := range commands {
					if command == "+" || command == "-" || command == "*" || command == "/" {
						if len(e.stack) < 2 {
							return e.stack, errors.New("incorrect operation")
						}

						var first = e.stack[len(e.stack)-1]
						var second = e.stack[len(e.stack)-2]
						e.stack = e.stack[:len(e.stack)-2]

						var result int
						switch command {
						case "+":
							result = second + first
						case "-":
							result = second - first
						case "*":
							result = second * first
						case "/":
							if first == 0 {
								return e.stack, errors.New("division by zero")
							}
							result = second / first
						}
						e.stack = append(e.stack, result)
					} else if command == "dup" {
						if len(e.stack) == 0 {
							return e.stack, errors.New("incorrect operation")
						}

						e.stack = append(e.stack, e.stack[len(e.stack)-1])
					} else if command == "over" {
						if len(e.stack) < 2 {
							return e.stack, errors.New("incorrect operation")
						}

						e.stack = append(e.stack, e.stack[len(e.stack)-2])
					} else if command == "drop" {
						if len(e.stack) == 0 {
							return e.stack, errors.New("incorrect operation")
						}

						e.stack = e.stack[:len(e.stack)-1]
					} else if command == "swap" {
						if len(e.stack) < 2 {
							return e.stack, errors.New("incorrect operation")
						}

						var first = e.stack[len(e.stack)-1]
						var second = e.stack[len(e.stack)-2]
						e.stack = e.stack[:len(e.stack)-2]
						e.stack = append(e.stack, first, second)
					} else {
						number, err := strconv.Atoi(command)
						if err != nil {
							return e.stack, errors.New(command + ": unknown word")
						} else {
							e.stack = append(e.stack, number)
						}
					}
				}
			}
		}
		i++
	}

	return e.stack, nil
}
