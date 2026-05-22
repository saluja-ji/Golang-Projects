# Terminal Typing Speed Test

A command-line typing speed test application written in Go that measures your typing speed and accuracy in real-time.

## Overview

This project is a learning exercise demonstrating fundamental Go concepts while building a practical typing test utility. The program generates random typing prompts, measures your typing performance, and calculates detailed statistics including WPM (Words Per Minute) and accuracy.

## Features

- **Random Prompt Selection**: Displays a random prompt from a predefined list of Go-related sentences
- **Real-time Measurement**: Accurately measures typing time using `time.Now()` and `time.Since()`
- **Character-Level Accuracy**: Compares character-by-character for precise accuracy measurement
- **WPM Calculation**: Calculates Net WPM using the standard typing formula: `(Correct Characters / 5) / Minutes`
- **Detailed Results**: Shows:
  - Time elapsed (in seconds)
  - Net WPM (words per minute)
  - Accuracy percentage
  - Correct characters typed
  - Correct words matched
- **Colored Output**: Uses ANSI escape codes for enhanced terminal UI

## How to Run

### Prerequisites
- Go 1.11 or later installed on your system

### Instructions

1. Navigate to the project directory:
   ```bash
   cd /home/pushpit-saluja/go/Projects/Terminal_Typing_Speed_Test
   ```

2. Run the program:
   ```bash
   go run main.go
   ```

3. Follow the on-screen prompts:
   - Read the target text displayed in cyan
   - Press ENTER to start the test
   - Type the text exactly as shown
   - Press ENTER to finish

4. View your results displayed in a formatted table

## Learning Concepts

This project demonstrates several key Go concepts:

### 1. **Packages and Imports**
   - `bufio` - Buffered input/output for efficient reading
   - `fmt` - Formatted input/output
   - `strings` - String manipulation utilities
   - `time` - Time measurement and operations
   - `math/rand` - Random number generation
   - `os` - Operating system operations

### 2. **Random Number Generation**
   - `rand.Seed()` - Initializes the random generator with a time-based seed
   - `rand.Intn()` - Generates a random integer within range
   - Good practice: Manual seeding is useful for learning even though Go 1.20+ auto-seeds

### 3. **Input/Output**
   - `bufio.NewReader()` - Creates a buffered reader
   - `reader.ReadString()` - Reads until delimiter (newline)
   - `fmt.Println()` and `fmt.Printf()` - Output formatting
   - ANSI escape codes (`\033[36m...` for colored output)

### 4. **String Processing**
   - `strings.TrimSpace()` - Removes leading/trailing whitespace
   - `strings.ToLower()` - Converts to lowercase for case-insensitive comparison
   - `strings.Fields()` - Splits string by whitespace efficiently

### 5. **Time Measurement**
   - `time.Now()` - Gets current time
   - `time.Since()` - Calculates elapsed duration
   - `duration.Seconds()` - Extracts seconds from duration
   - `duration.Minutes()` - Extracts minutes from duration

### 6. **Control Flow**
   - If-else statements for error handling
   - For loops for character and word comparison
   - Boundary checking to prevent index out of range errors

### 7. **Error Handling**
   - Checking `err != nil` after I/O operations
   - Using `os.Exit()` to terminate gracefully on errors

### 8. **Slices and Arrays**
   - String slices for storing multiple prompts
   - Random selection using slice indexing
   - Iteration through slices

## Code Structure

The main program flow is divided into 10 logical sections:

1. **Random Seeding** - Initialize random number generator
2. **Prompt Selection** - Choose a random prompt from the list
3. **UI Display** - Show welcome message and target text
4. **Input Reader** - Create buffered reader for user input
5. **Wait for Start** - Wait for user to press ENTER
6. **Timer and Input** - Start timer and read user input
7. **Clean Input** - Normalize strings for comparison
8. **Character Accuracy** - Calculate character-level accuracy
9. **WPM Calculation** - Compute words per minute
10. **Word Statistics** - Count correct words
11. **Display Results** - Show final statistics

## Example Output

```
========================================
      TERMINAL TYPING SPEED TEST        
========================================
Type the following text exactly as shown:

Go is an open source programming language designed to make it easy to build simple reliable and efficient software.

Press [ENTER] when you are ready to begin...
🚀 GO! Start typing... (Press [ENTER] when finished)

----------------------------------------
             TEST RESULTS               
----------------------------------------
Time Elapsed : 12.45 seconds
Net WPM      : 58 words/minute
Accuracy     : 96.5%
Correct Chars: 110 / 114
Correct Words: 18 / 19
========================================
```

## Algorithm Details

### Accuracy Calculation
Character-level comparison is used instead of word-level for more accurate results:
- Each character in target and user input is compared
- Only matching characters count as correct
- Prevents distortion from missing words

**Formula**: `(Correct Characters / Total Target Characters) × 100`

### WPM Calculation
Uses the standard typing test formula:
- Word count = `Correct Characters / 5`
- Time = measured duration in minutes
- **Formula**: `(Correct Characters / 5) / Minutes`

### Word Matching
Uses `strings.Fields()` which automatically handles multiple spaces better than `strings.Split()`.

## Possible Enhancements

- Add command-line flags for difficulty levels
- Load prompts from external file
- Add leaderboard/score history
- Add typing practice modes (specific topics)
- Calculate error rate (errors per minute)
- Add optional timeout for test
- Generate random prompts from dictionary
- Add color highlighting for correct/incorrect characters

## Key Takeaways

This project reinforces:
- Go's efficient standard library
- Proper error handling patterns
- Time measurement for performance analysis
- String manipulation techniques
- Clean code organization with comments
- Formula implementation for real-world calculations

---

**Created for Go Learning** 📚
