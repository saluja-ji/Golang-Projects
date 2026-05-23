
# CodeClock CLI — Learner Guide

Small CLI that times a single coding task and appends the session to `codeclock.json`.

**Why this is good for learners**
- **Hands-on**: Shows practical use of `bufio`, `time`, and `encoding/json`.
- **Small surface area**: Easy to read and modify a few functions.

**Files**
- **Source**: [CodeClockCLI/main.go](CodeClockCLI/main.go) — the program entry and core logic.
- **Data**: `codeclock.json` — where sessions are stored after each run.
- **Guide**: [CodeClockCLI/README.md](CodeClockCLI/README.md) — this learner-focused walkthrough.

**Quick Start**
```bash
cd CodeClockCLI
go run main.go
```

**What to watch for (quick walkthrough)**
- **`main()`**: Prompts for a task, records `startTime`, waits for Enter, then computes duration.
- **Reading input**: Uses `bufio.NewReader(os.Stdin)` and `ReadString('\n')` to capture the task name.
- **Timing**: `time.Now()` marks the start; `time.Since(startTime)` gives the elapsed time.
- **Rounding**: `duration.Round(time.Second)` simplifies the display to whole seconds.
- **Persistence**: `saveSession()` reads `codeclock.json` (if present), appends the new `Session`, and writes indented JSON with `json.MarshalIndent`.

**Key concepts (1–2 lines each)**n+- **Struct tags**: `Session` uses ``json:"task"`` to map Go fields to JSON keys when marshalling.
- **Error handling**: The program prints errors and exits early when critical IO fails — a common Go pattern.
- **File IO**: `os.ReadFile`/`os.WriteFile` make file operations concise for small files.

**Exercises for learners**
- **Add validation**: Prevent empty task names before starting the timer.
- **CSV export**: Write a small command to export `codeclock.json` to CSV.
- **Duration format**: Show duration in `mm:ss` instead of the default `1m30s` string.

**Next steps**
- Try the exercises above and open [CodeClockCLI/main.go](CodeClockCLI/main.go) while you edit.
- Want me to implement one of the exercises now? Tell me which one.

