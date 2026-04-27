# Guardian Scraper

## What is this project?

This is a simple web scraper written in Go that crawls The Guardian website (https://www.theguardian.com). It starts from the homepage and discovers all the links on the site that belong to the same domain. This is like a web crawler that explores a website automatically.

## Why would someone use this?

- To learn about web scraping in Go
- To collect all the URLs from a website
- As a starting point for more advanced scraping projects (like extracting articles or data)

## How does it work? (Beginner-friendly explanation)

Let's break down the code step by step, explaining the key Go concepts used.

### 1. **Structs in Go**
   - **What are structs?** Think of structs as blueprints for creating objects that hold data. They're like custom data types.
   - **In this project:** We use structs from Go's standard library, like `http.Client` and `http.Response`. For example:
     ```go
     client := &http.Client{Timeout: 10 * time.Second}
     ```
     This creates an HTTP client with a timeout. The `&` means we're creating a pointer to the struct.

### 2. **Goroutines**
   - **What are goroutines?** Goroutines are Go's way of running code concurrently (at the same time). They're lightweight threads that make it easy to do multiple things at once without blocking the main program.
   - **In this project:** We use goroutines to crawl multiple web pages at the same time. For example:
     ```go
     go func(link string) {
         found := crawl(link, baseURL, tokens)
         if found != nil {
             worklist <- found
         }
     }(link)
     ```
     The `go` keyword starts a new goroutine for each link we find.

### 3. **Channels**
   - **What are channels?** Channels are like pipes that goroutines use to communicate and synchronize. You can send data through channels or wait for data.
   - **In this project:** We use two types of channels:
     - `worklist chan []string`: A channel that holds lists of URLs to crawl.
     - `tokens chan struct{}`: A channel that limits how many goroutines can run at once (concurrency control). It's like a ticket system - only 5 goroutines can crawl at a time.

### 4. **Maps**
   - **What are maps?** Maps are like dictionaries in other languages - they store key-value pairs. In Go, they're efficient for lookups.
   - **In this project:** We use a map to keep track of URLs we've already seen:
     ```go
     seen := make(map[string]bool)
     ```
     This prevents us from crawling the same page twice.

### 5. **Other important concepts**
   - **HTTP Requests:** We make GET requests to fetch web pages using `http.NewRequest` and `http.Client`.
   - **HTML Parsing:** We use the `goquery` library (like jQuery for Go) to parse HTML and find links.
   - **URL Handling:** We use Go's `net/url` package to parse and resolve URLs properly.
   - **Randomization:** We rotate user agents to make our requests look more like a real browser.

## How to run the project

1. Make sure you have Go installed (version 1.25.0 or later).
2. Clone or download this project.
3. Run `go mod tidy` to install dependencies.
4. Run `go run main.go` to start the scraper.

The scraper will start crawling from The Guardian's homepage and print out each URL it visits.

## Important notes
- This scraper respects the website by limiting concurrent requests and adding delays.
- It's for educational purposes only. Always check a website's robots.txt and terms of service before scraping.
- Web scraping can be against the terms of service of some websites, so use responsibly!

## Learning resources
- [Go Tour](https://tour.golang.org/) - Great for beginners
- [Effective Go](https://golang.org/doc/effective_go.html) - Best practices
- [Go by Example](https://gobyexample.com/) - Practical examples

Feel free to experiment with the code and learn more about Go!</content>
<parameter name="filePath">/home/pushpit-saluja/go/src/guardian-scraper/README.md
