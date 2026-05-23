package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

// Session is a single timed task record stored as JSON.
// Fields include the task name, the date (YYYY-MM-DD), and the duration.
type Session struct {
	// JSON tags control the key names when saving/loading.
	Task     string `json:"task"`
	Date     string `json:"date"`
	Duration string `json:"duration"`
}

const dataFile = "codeclock.json"

func main() {
	fmt.Println("========================================")
	fmt.Println("             CODECLOCK CLI              ")
	fmt.Println("========================================")

	// Phase 1 — start timing the task
	fmt.Print("What task are you starting? : ")

	// Create a reader to capture text typed by the user in the terminal.
	reader := bufio.NewReader(os.Stdin)

	// Read the task name the user typed (up to newline) and clean it.
	taskName, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading the input:", err)
		return
	}

	taskName = strings.TrimSpace(taskName)

	// Record the current time as the session start time.
	startTime := time.Now()

	// Show the clock-in time in HH:MM:SS format.
	fmt.Printf("\nClocked IN at %s\n", startTime.Format("15:04:05"))

	fmt.Println("Press [ENTER] when you are finished...")

	// Pause here until the user presses Enter to stop the timer.
	bufio.NewReader(os.Stdin).ReadByte()

	// Phase 2 — stop timing and compute durations
	// Compute the exact elapsed time since `startTime`.
	exactDuration := time.Since(startTime)

	// Round to the nearest second for cleaner display.
	duration := exactDuration.Round(time.Second)

	// Print the rounded duration to the user.
	fmt.Printf("Clocked OUT! Time Spent: %s\n", duration)

	// Create a Session record with task, today's date, and duration.
	// Date uses YYYY-MM-DD format for consistency when saving.
	newSession := Session{
		Task:     taskName,
		Date:     time.Now().Format("2006-01-02"),
		Duration: duration.String(),
	}

	// Phase 3 — persist the session to the JSON file
	saveSession(newSession)
}

func saveSession(newSession Session) {
	var sessions []Session

	// Try reading existing sessions from disk; if missing, we'll create a new list.
	fileData, err := os.ReadFile(dataFile)
	if err == nil {
		json.Unmarshal(fileData, &sessions)
	}

	// Append the new session to the in-memory list.
	sessions = append(sessions, newSession)

	// Marshal the sessions slice into indented JSON for readability.
	updatedJSON, err := json.MarshalIndent(sessions, "", " ")
	if err != nil {
		fmt.Println("Error Creating JSON:", err)
		return
	}

	// Write the JSON back to disk with standard read/write permissions.
	err = os.WriteFile(dataFile, updatedJSON, 0644)
	if err != nil {
		fmt.Printf("Error saving file: %v\n", err)
	} else {
		fmt.Printf("Session saved to %s!\n", dataFile)
	}
}
