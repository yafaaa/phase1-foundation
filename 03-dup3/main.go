package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	counts := make(map[string]int)

	for _, filename := range os.Args[1:] {
		data, err := os.ReadFile(filename)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error found: %v\n", err)
			continue
		}

		for _, word := range strings.Fields(strings.TrimSpace(string(data))) {
			counts[word]++
		}

		for word, freq := range counts {
			if freq > 1 {
				fmt.Println(word)
			}
		}
	}
}
