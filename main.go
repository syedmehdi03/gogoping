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

	results := make(chan int, n)
	tasks := make(chan string, n)
	for range n {
		tasks <- parsedURL.String()
	}
	close(tasks)

	for range MAX_GOROUTINES {
		go pingURL2(&wg, tasks, results)
	}

	go func() {
		wg.Wait()
		close(results)
	}()
	for r := range results {
		fmt.Println(r)
	}
}

func pingURL(wg *sync.WaitGroup, mu *sync.Mutex, tasks *int, u string) int {
	defer wg.Done()
	mu.Lock()
	if *tasks > 0 {
		*tasks--
		mu.Unlock()

		client := getHttpClient()

		resp, err := client.Get(u)
		if err != nil {
			fmt.Printf("error occured: %v\n", err)
			return -1
		}
		fmt.Println(resp.StatusCode)
		return resp.StatusCode
	}
	mu.Unlock()
	return -1
}

func pingURL2(wg *sync.WaitGroup, tasks chan string, results chan int) int {
	defer wg.Done()
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

	return -1
}

func printFlags(wg *sync.WaitGroup, flags ...any) {
	defer wg.Done()
	for _, f := range flags {
		fmt.Println(f)
	}
}

func getHttpClient() *http.Client {
	return &http.Client{Timeout: 30 * time.Second}
}
