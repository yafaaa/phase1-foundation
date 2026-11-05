package main

import (
	"fmt"
	"io"
	"net/http"
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

func main(){

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	res1 := pingOnce (client, "https://google.com")
	res2 := pingOnce (client, "https://this-domian-does-not-exist-123.com")
	fmt.Printf("result of google: %+v \n", res1)
	fmt.Printf("result of domain: %+v \n", res2)

}