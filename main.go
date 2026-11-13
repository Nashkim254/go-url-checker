package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

type Result struct {
	URL     string
	Status  int
	Elapsed time.Duration
	Err     error
}

func main() {
	fmt.Println("Hello URLs")
	urls := []string{
		"https://google.com",
		"https://github.com",
		"https://golang.org",
		"https://stackoverflow.com",
		"https://reddit.com",
		"https://wikipedia.org",
		"https://youtube.com",
		"https://amazon.com",
		"https://microsoft.com",
		"https://apple.com",
		"https://cloudflare.com",
		"https://mozilla.org",
		"https://python.org",
		"https://rust-lang.org",
		"https://nodejs.org",
		"https://docker.com",
		"https://kubernetes.io",
		"https://npmjs.com",
		"https://medium.com",
		"https://news.ycombinator.com",
	}
	for _, n := range []int{1, 5, 10, 20} {
		start := time.Now()
		UrlChecker(urls, n)
		fmt.Printf("=== workers=%d took %s ===\n\n", n, time.Since(start))
	}
}

func UrlChecker(urls []string, numWorkers int) {
	jobs := make(chan string, len(urls))
	results := make(chan Result)

	// One client shared by all workers. It's safe for concurrent use
	// and reuses connections between requests.
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	// 1. Start the workers (Add before anyone calls Wait).
	var wg sync.WaitGroup
	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go Worker(i, client, jobs, results, &wg)
	}

	// 2. Queue all the jobs. jobs is buffered, so this doesn't block.
	for _, url := range urls {
		jobs <- url
	}
	close(jobs)

	// 3. Close results once every worker has finished sending.
	go func() {
		wg.Wait()
		close(results)
	}()

	// 4. Collect results. The loop ends when results is closed.
	for r := range results {
		if r.Err != nil {
			fmt.Printf("[ERROR] %s: %v\n", r.URL, r.Err)
			continue
		}
		fmt.Printf("[OK] %s - Status: %d - Time: %s\n", r.URL, r.Status, r.Elapsed)
	}
}

func CheckUrl(client *http.Client, url string) Result {
	start := time.Now()
	response, err := client.Get(url)
	elapsed := time.Since(start)
	if err != nil {
		return Result{URL: url, Elapsed: elapsed, Err: err}
	}
	defer response.Body.Close()

	return Result{URL: url, Status: response.StatusCode, Elapsed: elapsed}
}

func Worker(id int, client *http.Client, jobs <-chan string, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()

	for url := range jobs {
		results <- CheckUrl(client, url)
	}
}
