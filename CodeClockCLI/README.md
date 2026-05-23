# CodeClock CLI

Simple command-line tool to time coding tasks and save session records to `codeclock.json`.

How it works
- Run the CLI, enter a task name, press Enter to start, then press Enter again to stop.
- Each session is saved as a JSON object with `task`, `date` (YYYY-MM-DD), and `duration`.

Build & Run
```bash
cd CodeClockCLI
go run main.go
# or build a binary
go build -o codeclock .
./codeclock
```

Data file
- Sessions are stored in `codeclock.json` in the same directory. Keep this file in version control at your discretion.

Example
```
What task are you starting? : Write unit tests

Clocked IN at 14:05:23
Press [ENTER] when you are finished...
Clocked OUT! Time Spent: 1m30s
Session saved to codeclock.json!
```

Notes for learners
- The code shows basic use of `bufio`, `time`, and JSON marshalling in Go.
- `Session` uses struct tags to control JSON key names.
