package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	_ "modernc.org/sqlite"
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

	for _, url := range urls{
		fmt.Printf("\n_____fetching link_____: %v\n", url)
		for i:= 1; i<=5; i++{
			res := pingOnce(client, url)

			if err := savePing(db, res); err != nil{
				fmt.Printf("Ping was not saved to database: url: %v		rank: %v	Error: %v\n", res.URL, i, err)
			}
			
			if res.Err != nil{
				fmt.Printf("  [Ping %d/5] ERROR (%.2f ms): %v\n", i, res.LatencyMs, res.Err)
				continue
			}
			fmt.Printf("  [Ping %d/5] Status: %d | Latency: %.2f ms\n", i, res.StatusCode, res.LatencyMs)
			
			
		}
	}



}