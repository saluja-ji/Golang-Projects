package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
)

// urlStore keeps short codes mapped to original URLs in memory.
var urlStore = make(map[string]string)

type ShortenRequest struct {
	OriginalURL string `json:"originalUrl"`
}

type ShortenResponse struct {
	ShortCode string `json:"shortCode"`
}

// generateCode returns a random 6-character short code.
func generateCode() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	b := make([]byte, 6) // allocate space for the code

	for i := range b {
		b[i] = charset[rand.Intn(len(charset))] // pick a random character
	}

	return string(b)
}

// shortenHandler accepts a long URL, stores it, and returns a short code.
func shortenHandler(w http.ResponseWriter, r *http.Request) {
	var req ShortenRequest

	decoder := json.NewDecoder(r.Body) // read JSON from the request body
	err := decoder.Decode(&req)        // decode JSON into our struct

	if err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	if req.OriginalURL == "" {
		http.Error(w, "Original URL is required", http.StatusBadRequest)
		return
	}

	code := generateCode()
	urlStore[code] = req.OriginalURL // store the mapping in memory

	res := ShortenResponse{ShortCode: code}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res) // send JSON response back
}

// redirectHandler looks up the short code and redirects to the saved URL.
func redirectHandler(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code") // get code from the request path

	longURL, exists := urlStore[code] // check if code exists
	if !exists {
		http.Error(w, "404 - Short URL not found", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, longURL, http.StatusFound) // redirect to original URL
}

func main() {
	http.HandleFunc("POST /api/shorten", shortenHandler) // route for shortening URLs
	http.HandleFunc("GET /{code}", redirectHandler)      // route for redirects

	fmt.Println("🚀 Server running on port 8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("Server crashed: %v\n", err)
	}
}
