package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// List of user agents (add more for better rotation)
var userAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0 Safari/537.36",
}

// Return random user agent
func randomUserAgent() string {
	return userAgents[rand.Intn(len(userAgents))]
}

// Create HTTP request with headers
func getRequest(targetURL string) (*http.Response, error) {
	client := &http.Client{Timeout: 10 * time.Second}

	req, err := http.NewRequest("GET", targetURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", randomUserAgent())

	return client.Do(req)
}

// Extract and clean links
func discoverLinks(resp *http.Response, baseURL string) []string {
	foundUrls := []string{}

	if resp == nil {
		return foundUrls
	}

	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return foundUrls
	}

	base, err := url.Parse(baseURL)
	if err != nil {
		return foundUrls
	}

	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if !exists || href == "" || strings.HasPrefix(href, "#") {
			return
		}

		link, err := url.Parse(href)
		if err != nil {
			return
		}

		resolved := base.ResolveReference(link)

		// Restrict to same domain
		if strings.Contains(resolved.Host, base.Host) {
			foundUrls = append(foundUrls, resolved.String())
		}
	})

	return foundUrls
}

// Crawl a single page
func crawl(targetURL string, baseURL string, tokens chan struct{}) []string {
	// Acquire token (limit concurrency)
	tokens <- struct{}{}

	resp, err := getRequest(targetURL)

	// Release token when done
	defer func() { <-tokens }()

	if err != nil {
		fmt.Println("Request failed:", targetURL)
		return nil
	}

	// Optional: slow down requests (important)
	time.Sleep(300 * time.Millisecond)

	fmt.Println("Crawling:", targetURL)

	return discoverLinks(resp, baseURL)
}

func main() {
	rand.Seed(time.Now().UnixNano())

	baseURL := "https://www.theguardian.com"

	worklist := make(chan []string)
	seen := make(map[string]bool)

	// Limit concurrent requests
	tokens := make(chan struct{}, 5)

	// Seed initial URL
	go func() {
		worklist <- []string{baseURL}
	}()

	for list := range worklist {
		for _, link := range list {

			if !seen[link] {
				seen[link] = true

				go func(link string) {
					found := crawl(link, baseURL, tokens)

					if found != nil {
						worklist <- found
					}
				}(link)
			}
		}
	}
}
