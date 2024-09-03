//go:build !solution

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	wordsCount := make(map[string]int)

	for i := 1; i < len(os.Args); i++ {
		file, err := os.Open(os.Args[i])
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		input := bufio.NewScanner(file)
		for input.Scan() {
			wordsCount[input.Text()]++
		}
	}

	for word, count := range wordsCount {
		if count > 1 {
			fmt.Printf("%v\t%s\n", count, word)
		}
	}
}
