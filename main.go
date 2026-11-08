package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
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

type PingURLTask struct {
	URL        string
	StatusCode int
	Err        error
}

func NewPingURLTask(url string) PingURLTask {
	return PingURLTask{
		URL: url,
	}
}

func (p PingURLTask) Process() Task {
	t := PingURLTask{
		URL: p.URL,
	}
	resp, err := httpClient.Get(t.URL)
	if err != nil {
		t.Err = err
		return t
	}
	defer resp.Body.Close()
	t.StatusCode = resp.StatusCode
	return t
}

func main() {
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

	queue := NewQueueClient(10, 10, MAX_GOROUTINES)

	queue.AddToProducer()
	go func() {
		defer queue.ProducerDone()
		for range n {
			newTask := NewPingURLTask(parsedURL)
			queue.EnQueue(newTask)
		}
	}()

	queue.Run()

	for r := range queue.GetResults() {
		res := r.(PingURLTask)
		fmt.Printf("Status code: %v | Error: %v\n", res.StatusCode, res.Err)
	}
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
