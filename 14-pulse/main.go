package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type PingResult struct {
	URL string
	StatusCode int
	LatencyMs float64
	Err error
}

func pingOnce (client *http.Client, url string) PingResult {
	start := time.Now()
	resp, err := client.Get(url)
	latencyMs := float64(time.Since(start).Microseconds()) /1000.0

	if err != nil {
		return PingResult{
			URL: url,
			LatencyMs: latencyMs,
			Err: err,
		}
	}

	defer resp.Body.Close()

	_, _ = io.Copy(io.Discard, resp.Body)

	return PingResult{
		URL: url,
		StatusCode: resp.StatusCode,
		LatencyMs: latencyMs,
	}

}

func readURLS(filepath string) ([] string, error){
	file, err := os.Open(filepath)

	if err != nil {
		return nil, err
	}
	
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var urls [] string
	for scanner.Scan(){
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#"){
			continue
		}
		urls = append(urls, line)
	}

	if err := scanner.Err(); err != nil{
		return nil, err
	}
	return urls, nil

}

func main(){

	if len(os.Args) < 2 {
		fmt.Println("Enter file name and try again")
		os.Exit(1)
	}

	urls, err := readURLS(os.Args[1])

	if err != nil{
		fmt.Println("Unable to open file")
		os.Exit(1)
	}

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	for _, url := range urls{
		fmt.Printf("\n_____fetching link_____: %v\n", url)
		for i:= 1; i<=5; i++{
			res := pingOnce(client, url)
			
			if res.Err != nil{
				fmt.Printf("  [Ping %d/5] ERROR (%.2f ms): %v\n", i, res.LatencyMs, res.Err)
				continue
			}
			fmt.Printf("  [Ping %d/5] Status: %d | Latency: %.2f ms\n", i, res.StatusCode, res.LatencyMs)

		}
	}



}