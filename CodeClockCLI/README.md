# CodeClock CLI

Lightweight command-line app that times a coding task and persists session records to `codeclock.json`.

What you will learn
- How to read user input from the terminal using `bufio.NewReader` and `ReadString`.
- How to work with time in Go: `time.Now()`, `time.Since()`, and `Round` for human-friendly durations.
- How to define structs with JSON tags and marshal/unmarshal JSON (`encoding/json`).
- Basic file I/O with `os.ReadFile` and `os.WriteFile` and simple error handling.

Quick overview
- Run the CLI, type a task name and press Enter to start the timer.
- Press Enter again to stop; the program prints the duration and appends a session to `codeclock.json`.

Build & run
```bash
cd CodeClockCLI
go run main.go
# or build a binary
go build -o codeclock .
./codeclock
```

Example run
```
What task are you starting? : Write unit tests

Clocked IN at 14:05:23
Press [ENTER] when you are finished...
Clocked OUT! Time Spent: 1m30s
Session saved to codeclock.json!
```

File format
- Each entry in `codeclock.json` is an object with `task`, `date` (YYYY-MM-DD), and `duration` (string).

Walkthrough (learner-friendly)
- `Session` struct: shows how to annotate fields with JSON tags so keys match expected names.
- Input: `bufio.NewReader(os.Stdin)` and `ReadString('\n')` capture typed task names.
- Timing: `start := time.Now()` and `elapsed := time.Since(start).Round(time.Second)` produce neat durations.
- Persistence: existing data is read with `os.ReadFile`, parsed with `json.Unmarshal`, appended, and written back with `json.MarshalIndent` + `os.WriteFile`.

Suggested exercises
- Add a flag to show the last N sessions (use `flag` package).
- Save sessions in CSV instead of JSON to practice different encodings.
- Add a start/stop command mode so you can run `codeclock start` and `codeclock stop` separately.
- Add validation to prevent empty task names and handle corrupted `codeclock.json` gracefully.

Debugging tips
- If runs fail, print errors and inspect `codeclock.json` for valid JSON.
- Use `fmt.Printf("%#v\n", variable)` to inspect Go values during development.

Further reading
- Go by Example: IO and JSON examples (gobyexample.com)
- The `time` package docs: https://pkg.go.dev/time
