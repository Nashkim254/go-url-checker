package main

import (
	"fmt"
	"net/http"
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
	for _, url := range urls {
	go	CheckUrl(url)
	}
}

func CheckUrl(url string) {
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
