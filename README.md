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
pomodoro start                        # 25 min work, 5 min break (defaults)
pomodoro start -w 10 -b 3             # 10 min work, 3 min break
pomodoro start --mode tmux            # force tmux pane mode
pomodoro start --mode window          # force new terminal window
pomodoro start --mode inline          # single line in current terminal
pomodoro --version                    # print version
```

**Controls during a session:**
- `space` — pause / resume
- `q` — quit

### Display Modes

The timer supports multiple display modes, controlled via `--mode`. The default is `auto`.

| Mode     | Behavior |
|----------|----------|
| `auto`   | Detects environment automatically. Uses `tmux` if running inside a tmux session, otherwise opens a new `window`. |
| `tmux`   | Splits a small horizontal pane (8 lines) at the bottom of the current tmux window. The pane closes when the timer finishes or is quit. |
| `window` | Opens the timer in a new terminal window. Detects which terminal you're running in via `$TERM_PROGRAM` and opens the window in that same terminal. Supported: Ghostty, iTerm2, Terminal.app. Falls back to Terminal.app if detection fails, then to `inline` as a last resort. |
| `inline` | Renders a compact single-line progress bar in the current terminal. No screen clearing — stays out of your way. |

**Terminal detection (window mode):**

The launcher reads `$TERM_PROGRAM` to determine which terminal emulator is active and opens the new window there. This means if you run the command from Ghostty, you get a Ghostty window — not a Terminal.app window.

### How It Works

#### `main.go` — Entry Point & CLI Parsing

Handles argument parsing using Go's standard `flag` package. Defines two subcommands:

- `start` — the user-facing launcher. Accepts `-w` (work minutes, default 25), `-b` (break minutes, default 5), and `--mode` (display mode, default `auto`). Delegates to the launcher which decides how and where to display the timer.
- `_run` — internal command used by the launcher to actually run the timer in a spawned window or pane. Not intended to be called directly.

#### `launcher.go` — Environment Detection & Window Spawning

The brains behind display mode resolution. Key responsibilities:

**Mode resolution:** In `auto` mode, checks the `$TMUX` environment variable. If set, the user is inside tmux and the timer opens in a split pane. Otherwise, it opens a new terminal window.

**Terminal detection:** Reads `$TERM_PROGRAM` to identify the current terminal emulator (Ghostty, iTerm2, or Terminal.app) and opens the new window in the same application.

**Launch script:** Since passing complex arguments through `open --args` on macOS is unreliable, the launcher writes a temporary shell script to `/tmp/pomodoro-launch.sh` that runs the timer command and cleans itself up afterward. Each terminal launcher points to this script.

**Terminal-specific launchers:**
- **Ghostty** — uses `open -na Ghostty --args -e <script>`. Window closes automatically when the process exits.
- **iTerm2** — uses AppleScript to create a new window with the default profile running the script.
- **Terminal.app** — uses AppleScript to open a new tab, appends `; exit` to the command, polls until the process finishes, then closes the window automatically.

**Fallback chain:** If the detected terminal fails, falls back to Terminal.app. If that also fails, falls back to inline mode.

#### `timer.go` — Countdown Engine & Terminal Raw Mode

Two main responsibilities:

**Raw terminal mode:** To read keypresses without waiting for Enter, the terminal is switched from cooked mode into raw mode via `ioctl` syscalls that modify the terminal's `termios` settings — disabling `ECHO` (don't print typed characters), `ICANON` (don't buffer by line), and `ISIG` (don't interpret Ctrl+C as a signal). `VTIME` is set to give reads a 100ms timeout so input polling doesn't block the countdown. The original settings are saved and restored on exit.

**Countdown loop:** `countdown()` accepts a render function, allowing the same logic to drive both fullscreen and inline displays. It runs a `time.Ticker` at one-second intervals — each tick decrements remaining time and calls the render function. Between ticks, a non-blocking stdin read checks for `q` (quit) or space (pause toggle). `runTimer()` sequences the work and break phases, firing a terminal bell between them.

#### `display.go` — ANSI-Powered UI

Renders the timer using ANSI escape codes — no TUI framework needed. Provides two render modes:

**Fullscreen (`renderDisplay`):**
- Clears the screen and redraws on each tick
- Progress bar using `█` (filled) and `░` (empty) characters, 30 characters wide
- Color-coded: green for work, cyan for break
- Shows phase label, time remaining, percentage, and controls or pause state
- Hides the cursor during the timer, restores it on exit

**Inline (`renderInline`):**
- Single line updated in place via `\r` carriage return
- Compact 20-character progress bar with the same color coding
- Shows phase, progress, and time remaining — no screen clearing
- Designed to stay out of the way while you work

Both modes ring the terminal bell (`\a`) on phase transitions and completion.
