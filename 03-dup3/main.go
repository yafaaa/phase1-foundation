package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	counts := make(map[string]int)

	for _, filename := range os.Args[1:]{
		data, err := os.ReadFile(filename)

		if err != nil{
			fmt.Fprintf(os.Stderr, "dup3: %v\n", err)
			continue
		}

		
		for _, line := range strings.Split(string(data), "\n"){
			line = strings.TrimSpace(line)
    
			if line == "" {
				continue
			}
			
			counts[line]++		
		}

		for line, freq := range counts{
			if freq > 1{
				fmt.Printf("%d\t%s\n", freq, line)
			}
		}


	}
}