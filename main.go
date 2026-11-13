package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

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
	var wg sync.WaitGroup
	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go Worker(i, jobs, &wg)

	}
	for _, url := range urls {
		jobs <- url
	}
	close(jobs)
	wg.Wait()
}

func CheckUrl(url string) {
	start := time.Now()
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	response, err := client.Get(url)
	elapsed := time.Since(start)
	if err != nil {
		fmt.Printf("[ERROR] %s: %v\n", url, err)
		return
	}
	defer response.Body.Close()

	fmt.Printf("[OK] %s - Status: %d - Time: %s\n", url, response.StatusCode, elapsed)
}

func Worker(id int, jobs <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	for url := range jobs {
		fmt.Printf("Worker %d checking %s\n", id, url)
		CheckUrl(url)
	}

}
