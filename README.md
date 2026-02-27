## Pomodoro

A zero-dependency CLI pomodoro timer built with Go's standard library. Work timer, break timer, progress bar, keyboard controls — nothing fancy, nothing imported.

### Install from Source

Clone the repo

```
git clone https://github.com/wshivers93/pomodoro.git
```

Install the executable

```
go install
```

Verify installation
```
pomodoro --version
```

*If you receive a `command not found: pomodoro` error, you can either update GOBIN to be `usr/local/bin` or update your path to include the `go/bin` directory*

### Usage

```
pomodoro start              # 25 min work, 5 min break (defaults)
pomodoro start -w 10 -b 3   # 10 min work, 3 min break
pomodoro --version           # print version
```

**Controls during a session:**
- `space` — pause / resume
- `q` — quit

### How It Works

#### `main.go` — Entry Point & CLI Parsing

Handles argument parsing using Go's standard `flag` package. Uses `flag.NewFlagSet` to define the `start` subcommand with two optional flags: `-w` (work minutes, default 25) and `-b` (break minutes, default 5). Routes to `runTimer()` for the `start` command or prints version/usage info otherwise.

#### `timer.go` — Countdown Engine & Terminal Raw Mode

The core of the app. Two main responsibilities:

**Raw terminal mode:** To read keypresses without waiting for Enter, the terminal needs to be switched from its default (cooked) mode into raw mode. This is done via `ioctl` syscalls that modify the terminal's `termios` settings — specifically disabling `ECHO` (don't print typed characters), `ICANON` (don't buffer by line), and `ISIG` (don't interpret Ctrl+C as a signal). `VTIME` is set to give reads a 100ms timeout so input polling doesn't block the countdown. The original settings are saved and restored on exit.

**Countdown loop:** `countdown()` runs a `time.Ticker` at one-second intervals. Each tick decrements the remaining time and re-renders the display. Between ticks, it does a non-blocking read on stdin to check for `q` (quit) or space (pause toggle). The `runTimer()` function sequences the work and break phases, firing a terminal bell (`\a`) between them.

#### `display.go` — ANSI-Powered UI

Renders the timer display using ANSI escape codes — no TUI framework needed. Key pieces:

- **Screen control:** `\033[2J` clears the screen, `\033[H` moves the cursor home, `\033[?25l` / `\033[?25h` hide/show the cursor.
- **Progress bar:** Calculated from elapsed vs total seconds. Filled portion uses `█` characters, empty uses `░`, with the filled section colored green (work) or cyan (break).
- **Layout:** Each render clears and redraws the full display — phase label, progress bar with time and percentage, and a status line showing either pause state or available controls.
- **Completion:** `renderDone()` shows a checkmark message and rings the terminal bell.
