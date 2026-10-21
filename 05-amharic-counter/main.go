package main

import (
	"fmt"
	"sort"
	"strings"
	"unicode/utf8"
)

func main() {
	text := `ሰላም ኢትዮጵያ ሰላም አዲስ አበባ ኢትዮጵያ ሶፍትዌር ምህንድስና ሰላም ኢትዮጵያ ሶፍትዌር አዲስ አበባ ሰላም ኮምፒውተር ሳይንስ ኢትዮጵያ ሶፍትዌር ሰላም አዲስ ሶፍትዌር ምህንድስና`
	counts := make(map[string]int)
	allwords := strings.Fields(text)
	len_rune := 0

	for _, word := range allwords {
		counts[word]++
		len_rune += utf8.RuneCountInString(word)
	}

	fmt.Println("The total Rune: ", len_rune)

	var unique []string
	for  word := range counts{
		unique = append(unique, word)
	}

	sort.Slice(unique, func(i, j int) bool {
		return counts[unique[i]] > counts[unique[j]]
	})

	limit := 5

	if len(unique) < limit {
		limit = len(unique)
	}

	for i, word := range unique[:limit] {
		fmt.Println("Rank ", i+1, word, "frequency: ", counts[word], "Rune: ", utf8.RuneCountInString(word), "Byte: ", len(word))
	}

}
