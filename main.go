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
	MAX_GOROUTINES = 50
	DEFAULT_N      = 1
	DEFAULT_URL    = "https://www.example.com"
)

var (
	httpClient = &http.Client{Timeout: 30 * time.Second}
)

func main() {
	l := log.Default()
	var rawURL string
	var n int

	flag.StringVar(&rawURL, "url", DEFAULT_URL, "")
	flag.IntVar(&n, "n", DEFAULT_N, "N describes the number of times requests are made to the specified URL")
	flag.Parse()

	parsedURL, err := url.ParseRequestURI(rawURL)
	if err != nil {
		l.Fatalf("invalid url")
	}
	if n < 1 {
		l.Fatalf("N must be greater than 0")
	}

	wg := sync.WaitGroup{}
	wg.Add(MAX_GOROUTINES)

	results := make(chan int, 5)
	tasks := make(chan string, 5)

	go func() {
		for range n {
			tasks <- parsedURL.String()
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
		fmt.Println(r)
	}
}

func worker(tasks chan string, results chan int) {
	for t := range tasks {
		code, err := getStatusCode(t)
		if err != nil {
			fmt.Printf("worker failed for task: %v error: %v\n", t, err)
			continue
		}
		results <- code
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
