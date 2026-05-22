package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

func main() {

	// Seed the random generator
	// NOTE:
	// In newer Go versions this is auto-seeded,
	// but doing this manually is still a good habit for learning.
	rand.Seed(time.Now().UnixNano())

	// =====================================================
	// 1. Sample prompts for the typing test
	// =====================================================

	prompts := []string{
		"Go is an open source programming language designed to make it easy to build simple reliable and efficient software.",
		"Concurrency is not parallelism it is a way to structure a program by breaking it into independently executing components.",
		"The Go standard library is powerful enough to build robust web servers with minimal external dependencies.",
		"Channels are used to communicate between goroutines allowing safe data sharing.",
	}

	// =====================================================
	// 2. Select a random prompt
	// =====================================================

	// Generate a random index
	randomIndex := rand.Intn(len(prompts))

	// Pick the target text from the slice
	targetText := prompts[randomIndex]

	// =====================================================
	// 3. Display UI
	// =====================================================

	fmt.Println("========================================")
	fmt.Println("      TERMINAL TYPING SPEED TEST        ")
	fmt.Println("========================================")
	fmt.Println("Type the following text exactly as shown:")
	fmt.Println()

	// ANSI escape code for cyan color
	fmt.Printf("\033[36m%s\033[0m\n\n", targetText)

	fmt.Println("Press [ENTER] when you are ready to begin...")

	// =====================================================
	// 4. Create a Reader for user input
	// =====================================================

	reader := bufio.NewReader(os.Stdin)

	// Wait for the user to press ENTER
	_, err := reader.ReadString('\n')

	// Always handle errors properly
	if err != nil {
		fmt.Printf("Error reading input: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("🚀 GO! Start typing... (Press [ENTER] when finished)")

	// =====================================================
	// 5. Start Timer
	// =====================================================

	startTime := time.Now()

	// Read user input
	userInput, err := reader.ReadString('\n')

	if err != nil {
		fmt.Printf("Error reading input: %v\n", err)
		os.Exit(1)
	}

	// Stop timer
	duration := time.Since(startTime)

	// =====================================================
	// 6. Clean Input
	// =====================================================

	// Remove trailing spaces and newline
	userInput = strings.TrimSpace(userInput)

	// Normalize text for fair comparison
	// Example:
	// "Go" and "go" should be treated equally
	targetText = strings.ToLower(targetText)
	userInput = strings.ToLower(userInput)

	// =====================================================
	// 7. Calculate Character-Level Accuracy
	// =====================================================

	// Why character-level?
	//
	// Word-level comparison is weak because
	// one missing word shifts everything.
	//
	// Character comparison is much more accurate
	// and closer to real typing tests.

	correctChars := 0

	// Compare character by character
	for i := 0; i < len(targetText) && i < len(userInput); i++ {

		if targetText[i] == userInput[i] {
			correctChars++
		}
	}

	// Prevent division-by-zero
	accuracy := 0.0

	if len(targetText) > 0 {
		accuracy = (float64(correctChars) / float64(len(targetText))) * 100
	}

	// =====================================================
	// 8. Calculate Net WPM
	// =====================================================

	// Standard typing formula:
	//
	// WPM = (Characters / 5) / Minutes
	//
	// We use CORRECT characters instead of total typed
	// to make the score realistic.

	wpm := 0.0

	if duration.Minutes() > 0 {
		wpm = (float64(correctChars) / 5.0) / duration.Minutes()
	}

	// =====================================================
	// 9. Calculate Word Statistics
	// =====================================================

	// Fields() is better than Split()
	//
	// Split("a  b", " ")
	// gives empty strings too.
	//
	// Fields() automatically handles multiple spaces.

	targetWords := strings.Fields(targetText)
	userWords := strings.Fields(userInput)

	correctWords := 0

	for i := 0; i < len(targetWords) && i < len(userWords); i++ {

		if targetWords[i] == userWords[i] {
			correctWords++
		}
	}

	// =====================================================
	// 10. Display Final Results
	// =====================================================

	fmt.Println("\n----------------------------------------")
	fmt.Println("             TEST RESULTS               ")
	fmt.Println("----------------------------------------")

	fmt.Printf("Time Elapsed : %.2f seconds\n", duration.Seconds())
	fmt.Printf("Net WPM      : %.0f words/minute\n", wpm)
	fmt.Printf("Accuracy     : %.1f%%\n", accuracy)
	fmt.Printf("Correct Chars: %d / %d\n", correctChars, len(targetText))
	fmt.Printf("Correct Words: %d / %d\n", correctWords, len(targetWords))

	fmt.Println("========================================")
}
