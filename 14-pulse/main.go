package main

import (
	"bufio"
	"database/sql"
	"flag"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"strings"
	"sync"
	"text/tabwriter"
	"time"
	_ "modernc.org/sqlite"
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorCyan   = "\033[36m"
	colorBold   = "\033[1m"
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
	);
	
	-- NEW: Composite B-Tree Index for fast time-series queries
	CREATE INDEX IF NOT EXISTS idx_pings_url_created_at ON pings (url, created_at);
	`

	if _, err := db.Exec(schema); err != nil{
		db.Close()
		return nil, err
	}

	return db, nil
}

func runHistoryCommand(db *sql.DB, args []string) {
	historyCmd := flag.NewFlagSet("history", flag.ExitOnError)
	urlFlag := historyCmd.String("url", "", "Target URL to inspect (e.g., https://google.com)")
	lastFlag := historyCmd.String("last", "1h", "Lookback duration (e.g., 15m, 1h, 24h)")

	
	_ = historyCmd.Parse(args)

	duration, err := time.ParseDuration(*lastFlag)
	if err != nil {
		fmt.Printf("Invalid --last duration %q: %v (use formats like 30m, 1h, 24h)\n", *lastFlag, err)
		os.Exit(1)
	}
	
	targetURL := strings.TrimSpace(*urlFlag)
	if targetURL == "" {
		fmt.Println("Error: --url flag is required (e.g., pulse history --url https://google.com --last 1h)")
		os.Exit(1)
	}

	if !strings.HasPrefix(targetURL, "http://") && !strings.HasPrefix(targetURL, "https://") {
		targetURL = "https://" + targetURL
	}
	
	cutoff := time.Now().UTC().Add(-duration).Format("2006-01-02 15:04:05")

	query := `
	SELECT
		COUNT(*),
		AVG(latency_ms),
		MIN(latency_ms),
		MAX(latency_ms)
	FROM pings
	WHERE url = ? AND created_at >= ? AND status_code > 0;`

	var count int
	var avgLatency, minLatency, maxLatency sql.NullFloat64

	err = db.QueryRow(query, targetURL, cutoff).Scan(&count, &avgLatency, &minLatency, &maxLatency)
	if err != nil {
		fmt.Printf("Database query failed: %v\n", err)
		os.Exit(1)
	}
	
	fmt.Printf("\n=== Pulse History: %s (last %s) ===\n", targetURL, duration)
	if count == 0 || !avgLatency.Valid {
		fmt.Println("No successful pings recorded in this time window.")
		return
	}

	fmt.Printf("Successful Pings : %d\n", count)
	fmt.Printf("Average Latency  : %.2f ms\n", avgLatency.Float64)
	fmt.Printf("Min Latency      : %.2f ms\n", minLatency.Float64)
	fmt.Printf("Max Latency      : %.2f ms\n", maxLatency.Float64)

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

type PingRow struct {
	StatusCode int
	LatencyMs  float64
	CreatedAt  string
}

type StatsReport struct {
	TotalCount   int
	SuccessCount int
	PacketLoss   float64
	AvgLatency   float64
	MinLatency   float64
	MaxLatency   float64
	JitterMs     float64
	Trend        string
}

func calculateStats(row []PingRow) StatsReport{
	report := StatsReport{TotalCount: len(row)}

	if len(row) == 0{
		return report
	}
	
	var latencies []float64
	var sum float64

	for _, r := range row{

		if r.StatusCode > 0{
			latencies = append(latencies, r.LatencyMs)
			sum += r.LatencyMs

			if len(latencies) == 1 || r.LatencyMs < report.MinLatency{
				report.MinLatency = r.LatencyMs
			}

			if r.LatencyMs > report.MaxLatency{
				report.MaxLatency = r.LatencyMs
			}	
		}
	}
	report.PacketLoss = float64((report.TotalCount - report.SuccessCount)*100)/float64(report.TotalCount)
	
	if report.SuccessCount == 0{
		report.Trend = "DOWN"
		return report
	}
	
	report.AvgLatency = sum/float64(report.SuccessCount)

	if report.SuccessCount >= 2{
		var varianceSum float64
		for _, lat := range latencies{
			diff := lat - report.AvgLatency
			varianceSum += diff * diff
		}
		
		report.JitterMs = math.Sqrt(varianceSum/float64(len(latencies)))
	}
	

	if report.SuccessCount >= 6{
		SplitIdx := report.SuccessCount - 5

		baseline := mean(latencies[:SplitIdx])
		recentAvg := mean(latencies[SplitIdx:])

		if recentAvg > baseline * 1.25{
			report.Trend = "DEGRADING"
		}else if recentAvg < baseline * 0.85{
			report.Trend = "IMPROVING"
		}else {
			report.Trend = "STABLE"
		}

	}else{
		report.Trend = "INSUFFICIENT DATA"
	}

	return report
}

func mean(vals []float64) float64{
	var sum float64

	for _, val := range vals{
		sum += val
	}

	if len(vals) == 0 {
		return sum
	}
	return sum/float64(len(vals))
}

func calculateHealthScore(r StatsReport) int{
	if r.TotalCount == 0 || r.SuccessCount == 0 {
		return 0
	}
	score := 100.0

	score -= r.PacketLoss * 2.0

	if r.AvgLatency > 100.0 {
		pt := (r.AvgLatency - 100.0) * 0.05

		if pt > 35.0{
			pt = 35.0
		}
		score -= pt
	}

	if r.JitterMs > 15.0{
		pt:= (r.JitterMs - 15.0) * 0.25

		if pt > 15.0{
			pt = 15.0
		}
		score -= pt
	}

	if score < 0{
		return 0
	}

	return int(math.Round(score))
}

func printStatsTable(stats StatsReport, score int) {
	scoreColor := colorGreen
	statusLabel := "HEALTHY"
	if score < 50 {
		scoreColor = colorRed
		statusLabel = "CRITICAL"
	} else if score < 80 {
		scoreColor = colorYellow
		statusLabel = "DEGRADED"
	}

	trendColor := colorCyan
	if stats.Trend == "DEGRADING" || stats.Trend == "DOWN" {
		trendColor = colorRed
	} else if stats.Trend == "IMPROVING" {
		trendColor = colorGreen
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "METRIC\tVALUE\tSTATUS")
	fmt.Fprintln(w, "------\t-----\t------")
	fmt.Fprintf(w, "Health Score\t%s%d / 100%s\t%s\n", scoreColor, score, colorReset, statusLabel)
	fmt.Fprintf(w, "Packet Loss\t%.1f%%\t(%d/%d succeeded)\n", stats.PacketLoss, stats.SuccessCount, stats.TotalCount)
	fmt.Fprintf(w, "Avg Latency\t%.2f ms\tmin: %.2f ms | max: %.2f ms\n", stats.AvgLatency, stats.MinLatency, stats.MaxLatency)
	fmt.Fprintf(w, "Jitter (StdDev)\t%.2f ms\t-\n", stats.JitterMs)
	fmt.Fprintf(w, "Trend\t%s%s%s\t(last 5 vs baseline)\n", trendColor, stats.Trend, colorReset)
	w.Flush()
}

func printASCIIGraph(rows []PingRow) {
	fmt.Printf("\n%sLast %d Pings (Oldest -> Newest):%s\n", colorBold, len(rows), colorReset)

	var maxLat float64
	for _, r := range rows {
		if r.StatusCode > 0 && r.LatencyMs > maxLat {
			maxLat = r.LatencyMs
		}
	}
	if maxLat == 0 {
		maxLat = 1.0
	}

	const maxBarWidth = 35
	for i, r := range rows {
		if r.StatusCode == 0 {
			fmt.Printf("  #%02d | %s✖ FAILED / TIMEOUT (%.0f ms)%s\n", i+1, colorRed, r.LatencyMs, colorReset)
			continue
		}

		barLen := int(math.Round((r.LatencyMs / maxLat) * maxBarWidth))
		if barLen < 1 {
			barLen = 1
		}

		barColor := colorGreen
		if r.LatencyMs > 250.0 {
			barColor = colorYellow
		}

		bar := strings.Repeat("█", barLen)
		fmt.Printf("  #%02d | %s%s%s %.2f ms\n", i+1, barColor, bar, colorReset, r.LatencyMs)
	}
}

func runAnalyzeCommand(db *sql.DB, args []string){
	analysisCmd := flag.NewFlagSet("analyze", flag.ExitOnError)
	urlFlag := analysisCmd.String("url", "", "Target URL to analyze (e.g., google.com)")
	_ = analysisCmd.Parse(args)

	targetURL := strings.TrimSpace(*urlFlag)
	if targetURL == ""{
		fmt.Println("Error: --url flag is required")
		os.Exit(1)
	}
	if !strings.HasPrefix(targetURL, "https://") || !strings.HasPrefix(targetURL, "https://"){
		targetURL = "https://" + targetURL
	}

	query := `
	SELECT status_code, latency_ms, created_at
	FROM pings
	WHERE url = ?
	ORDER BY id DESC
	LIMIT 20;`

	dbRows, err := db.Query(query, targetURL)

	if err != nil{
		fmt.Println("Failed to process query")
		os.Exit(1)
	}

	defer dbRows.Close()

	var rows []PingRow
	for dbRows.Next(){
		var r PingRow
		if err := dbRows.Scan(&r.StatusCode, &r.LatencyMs, &r.CreatedAt); err != nil{
			fmt.Printf("Row scan failed: %v\n", err)
			os.Exit(1)
		}
		rows = append(rows, r)
	}
	if err := dbRows.Err(); err != nil {
		fmt.Printf("Row iteration error: %v\n", err)
		os.Exit(1)
	}

	if len(rows) == 0{
		fmt.Printf("No ping records found for %s. Run `./pulse urls.txt` first.\n", targetURL)
		return
	}

	for i, j := 0, len(rows)-1; i < j; i, j = i+1, j-1 {
		rows[i], rows[j] = rows[j], rows[i]
	}

	stats := calculateStats(rows)
	score := calculateHealthScore(stats)

	fmt.Printf("\n%s=== Pulse Network Analysis: %s ===%s\n\n", colorBold, targetURL, colorReset)
	printStatsTable(stats, score)
	printASCIIGraph(rows)

}


func main(){

	if len(os.Args) < 2 {
		fmt.Println("Usage:")
		fmt.Println("  pulse <urls-file>")
		fmt.Println("  pulse history --url <url> --last <duration>")
		fmt.Println("  pulse analyze --url <url>")
		os.Exit(1)
	}

	db, err := initDB("pulse.db")
	if err != nil {
		fmt.Printf("Failed to initialize database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	if os.Args[1] == "history" {
		runHistoryCommand(db, os.Args[2:])
		return
	}
	if os.Args[1] == "analyze" {
		runAnalyzeCommand(db, os.Args[2:])
		return
	}
	urls, err := readURLS(os.Args[1])

	if err != nil{
		fmt.Printf("Unable to open file: %v\n", err)
		os.Exit(1)
	}

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