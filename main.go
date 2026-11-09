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
	}
	UrlChecker(urls)
}

func UrlChecker(urls []string) {
	wg := new(sync.WaitGroup)
	for _, url := range urls {
		wg.Add(1)
		go CheckUrl(url, wg)

	}
	wg.Wait()
}

func CheckUrl(url string, wg *sync.WaitGroup) {
	defer wg.Done()
	start := time.Now()
	response, err := http.Get(url)
	elapsed := time.Since(start)
	if err != nil {
		fmt.Printf("[ERROR] %s: %v\n", url, err)
		return
	}
	defer response.Body.Close()

	fmt.Printf("[OK] %s - Status: %d - Time: %s\n", url, response.StatusCode, elapsed)
}
