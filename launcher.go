package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
)

func inTmux() bool {
	return os.Getenv("TMUX") != ""
}

func resolveMode(mode string) string {
	if mode != "auto" {
		return mode
	}
	if inTmux() {
		return "tmux"
	}
	return "window"
}

func selfPath() string {
	path, err := os.Executable()
	if err != nil {
		path, err = exec.LookPath("pomodoro")
		if err != nil {
			return os.Args[0]
		}
	}
	return path
}

func launch(workMins, breakMins int, mode string) {
	resolved := resolveMode(mode)

	switch resolved {
	case "inline":
		runTimer(workMins, breakMins, true)
	case "tmux":
		launchTmux(workMins, breakMins)
	case "window":
		launchWindow(workMins, breakMins)
	default:
		fmt.Printf("Unknown mode: %s\n", resolved)
		os.Exit(1)
	}
}

func launchTmux(workMins, breakMins int) {
	bin := selfPath()
	runCmd := fmt.Sprintf("%s _run -w %d -b %d -mode fullscreen",
		bin, workMins, breakMins)

	cmd := exec.Command("tmux", "split-window", "-v", "-l", "8", runCmd)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open tmux pane: %v\n", err)
		fmt.Fprintln(os.Stderr, "Falling back to inline mode.")
		runTimer(workMins, breakMins, true)
	}
}

func launchWindow(workMins, breakMins int) {
	bin := selfPath()
	w := strconv.Itoa(workMins)
	b := strconv.Itoa(breakMins)

	term := os.Getenv("TERM_PROGRAM")

	switch term {
	case "ghostty":
		if launchGhostty(bin, w, b) {
			return
		}
	case "iTerm.app":
		if launchITerm(bin, w, b) {
			return
		}
	case "Apple_Terminal":
		if launchTerminalApp(bin, w, b) {
			return
		}
	}

	if term != "Apple_Terminal" {
		if launchTerminalApp(bin, w, b) {
			return
		}
	}

	fmt.Fprintln(os.Stderr, "Could not open a new terminal window.")
	fmt.Fprintln(os.Stderr, "Falling back to inline mode.")
	runTimer(workMins, breakMins, true)
}

// writeLaunchScript creates a temp shell script that runs the timer
// and cleans itself up afterward.
func writeLaunchScript(bin, w, b string) (string, error) {
	tmp := filepath.Join(os.TempDir(), "pomodoro-launch.sh")
	content := fmt.Sprintf("#!/bin/sh\n%s _run -w %s -b %s -mode fullscreen\nrm -f %s\n", bin, w, b, tmp)
	if err := os.WriteFile(tmp, []byte(content), 0755); err != nil {
		return "", err
	}
	return tmp, nil
}

func launchGhostty(bin, w, b string) bool {
	script, err := writeLaunchScript(bin, w, b)
	if err != nil {
		return false
	}

	cmd := exec.Command("open", "-na", "Ghostty", "--args", "-e", script)
	if err := cmd.Start(); err != nil {
		os.Remove(script)
		return false
	}
	return true
}

func launchITerm(bin, w, b string) bool {
	script, err := writeLaunchScript(bin, w, b)
	if err != nil {
		return false
	}

	appleScript := fmt.Sprintf(`tell application "iTerm"
	activate
	create window with default profile command "%s"
end tell`, script)

	cmd := exec.Command("osascript", "-e", appleScript)
	if err := cmd.Run(); err != nil {
		os.Remove(script)
		return false
	}
	return true
}

func launchTerminalApp(bin, w, b string) bool {
	script, err := writeLaunchScript(bin, w, b)
	if err != nil {
		return false
	}

	// Terminal.app keeps windows open after process exits.
	// We use AppleScript to run the command, wait for it to finish,
	// then close the specific tab/window.
	appleScript := fmt.Sprintf(`tell application "Terminal"
	activate
	set newTab to do script "%s; exit"
	-- wait for the process to finish
	repeat
		delay 1
		if not busy of newTab then exit repeat
	end repeat
	close window 1 of (every window whose tabs contains newTab)
end tell`, script)

	cmd := exec.Command("osascript", "-e", appleScript)
	if err := cmd.Start(); err != nil {
		os.Remove(script)
		return false
	}
	return true
}
