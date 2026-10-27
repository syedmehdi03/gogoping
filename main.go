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
			pingURL(tasks, results)
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

func pingURL(tasks chan string, results chan int) {
	client := getHttpClient()

	for t := range tasks {
		resp, err := client.Get(t)
		if err != nil {
			fmt.Printf("error in goroutine: %v\n", err)
			continue
		}
		results <- resp.StatusCode
		resp.Body.Close()
	}
}

func getHttpClient() *http.Client {
	return &http.Client{Timeout: 30 * time.Second}
}
