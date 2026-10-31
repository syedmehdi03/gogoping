package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"
)

const (
	MAX_GOROUTINES = 5
	DEFAULT_N      = 1
	DEFAULT_URL    = "https://www.example.com"
)

var (
	httpClient = &http.Client{Timeout: 30 * time.Second}
)

type Response struct {
	Status    int
	Err       error
	TimeTaken time.Duration
}

func main() {
	now := time.Now()
	l := log.Default()
	var rawURL string
	var n int

	flag.StringVar(&rawURL, "url", DEFAULT_URL, "")
	flag.IntVar(&n, "n", DEFAULT_N, "N describes the number of times requests are made to the specified URL")
	flag.Parse()

	parsedURL, n, err := validateFlags(rawURL, n)
	if err != nil {
		l.Fatal(err)
	}

	wg := sync.WaitGroup{}
	wg.Add(MAX_GOROUTINES)

	results := make(chan Response, 5)
	tasks := make(chan string, 5)

	go func() {
		for range n {
			tasks <- parsedURL
		}
		close(tasks)
	}()

	for range MAX_GOROUTINES {
		go func() {
			worker(tasks, results)
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(results)
	}()
	for r := range results {
		fmt.Printf("Status code: %v | Error: %v | TimeTaken: %v\n", r.Status, r.Err, r.TimeTaken)
	}

	fmt.Printf("Total time: %v", time.Since(now))
}

func worker(tasks chan string, results chan Response) {
	for t := range tasks {
		start := time.Now()
		code, err := getStatusCode(t)
		elapsed := time.Since(start)

		if err != nil {
			results <- Response{
				Err:       err,
				TimeTaken: elapsed,
			}
			continue
		}
		results <- Response{
			Status:    code,
			TimeTaken: elapsed,
		}
	}
}

func getStatusCode(url string) (int, error) {
	resp, err := httpClient.Get(url)
	if err != nil {
		return -1, err
	}
	defer resp.Body.Close()
	return resp.StatusCode, nil
}

func validateFlags(rawURL string, n int) (string, int, error) {
	parsedURL, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return "", 0, fmt.Errorf("invalid URL")
	}
	if n < 1 {
		return "", 0, fmt.Errorf("N must be greater than 1")
	}
	return parsedURL.String(), n, nil
}
