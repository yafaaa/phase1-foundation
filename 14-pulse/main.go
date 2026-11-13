package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	_ "modernc.org/sqlite"
)

type PingResult struct {
	URL string
	PingNum int
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

func initDB(dbpath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dbpath)

	if err != nil{
		return nil, err
	}
	
	schema := `
	CREATE TABLE IF NOT EXISTS pings (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		url TEXT NOT NULL,
		status_code INTEGER NOT NULL,
		latency_ms REAL NOT NULL,
		error_msg TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	if _, err := db.Exec(schema); err != nil{
		db.Close()
		return nil, err
	}

	return db, nil
}

func savePing(db *sql.DB, res PingResult) error {
	var err_Msg string

	if res.Err != nil{
		err_Msg = res.Err.Error()
	}
	query := `Insert into pings (url, status_code, latency_ms, error_msg) VALUES(?, ?, ?, ?)`
	
	if _, err := db.Exec(query, res.URL, res.StatusCode, res.LatencyMs, err_Msg); err != nil{
		return err
	}
	return nil
}

func checkURL(client *http.Client, url string, wg *sync.WaitGroup, sem chan struct{}, results chan<- PingResult){
	defer wg.Done()

	sem <- struct{}{}
	defer func ()  { <-sem}()

	for i := 1; i <= 5; i++{
		res := pingOnce(client, url)
		res.PingNum = i
		results <- res
	}
}

func main(){

	if len(os.Args) < 2 {
		fmt.Println("Enter file name and try again")
		os.Exit(1)
	}

	urls, err := readURLS(os.Args[1])

	if err != nil{
		fmt.Printf("Unable to open file: %v\n", err)
		os.Exit(1)
	}

	db, err := initDB("pulse.db")

	if err != nil{
		fmt.Printf("Failed to intialize database: %v\n", err)
		os.Exit(1)
	}

	defer db.Close()

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	var wg sync.WaitGroup
	sem := make(chan struct{}, 5)
	results := make(chan PingResult)

	for _, url := range urls{
		wg.Add(1)
		go checkURL(client, url, &wg, sem, results)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for res := range results {
		if err := savePing(db, res); err != nil {
			fmt.Printf("  [DB Error] %s: %v\n", res.URL, err)
		}

		if res.Err != nil {
			fmt.Printf("  [Ping %d/5] %-35s | ERROR (%.2f ms): %v\n", res.PingNum, res.URL, res.LatencyMs, res.Err)
			continue
		}
		fmt.Printf("  [Ping %d/5] %-35s | Status: %d | Latency: %.2f ms\n", res.PingNum, res.URL, res.StatusCode, res.LatencyMs)
	}


}