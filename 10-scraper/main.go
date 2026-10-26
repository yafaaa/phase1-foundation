package main

import (
	"fmt"
	"net/http"
	"sync"
)

func main() {
	urls := []string{
		"https://example.com",
		"https://go.dev",
		"https://github.com",
		"https://google.com",
		"https://wikipedia.org",
		"https://stackoverflow.com",
		"https://cloudflare.com",
		"https://mozilla.org",
		"https://ubuntu.com",
		"https://python.org",
	}

	// Pre-allocate slice so each goroutine 'i' can store its result at results[i]
	results := make([]string, len(urls))

	var wg sync.WaitGroup
	// TODO 1: Create a buffered channel 'sem' of capacity 3 (e.g., make(chan struct{}, 3))
	sem := make(chan struct{}, 3)

	for i, url := range urls {
		wg.Add(1)

		go func(index int, targetURL string) {
			defer wg.Done()

			// TODO 2: Acquire a slot by sending into 'sem', and defer releasing it (<-sem)
			sem <- struct{}{}
			defer func(){ <- sem}()
			resp, err := http.Get(targetURL)
			if err != nil {
				results[index] = fmt.Sprintf("[%d] %s -> ERROR: %v", index, targetURL, err)
				return
			}
			defer resp.Body.Close()

			results[index] = fmt.Sprintf("[%d] %s -> Status: %s", index, targetURL, resp.Status)
		}(i, url)
	}

	// TODO 3: Wait for all goroutines to finish using wg.Wait(), then loop over 'results' and print each one in order [0..9]!
	wg.Wait()

	for _, res := range results{
		fmt.Println(res)
	}

}