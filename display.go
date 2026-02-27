package main

import (
	"fmt"
	"strings"
)

const (
	// ANSI escape codes
	clearScreen = "\033[2J"
	moveTo00    = "\033[H"
	bold        = "\033[1m"
	reset       = "\033[0m"
	red         = "\033[31m"
	green       = "\033[32m"
	yellow      = "\033[33m"
	cyan        = "\033[36m"
	dim         = "\033[2m"
	hideCursor  = "\033[?25l"
	showCursor  = "\033[?25h"

	barWidth       = 30
	inlineBarWidth = 20
)

// --- Fullscreen display ---

func clearAndReset() {
	fmt.Print(clearScreen + moveTo00)
}

func renderDisplay(label string, remaining int, total int, paused bool) {
	clearAndReset()

	color := green
	if label == "Break" {
		color = cyan
	}

	elapsed := total - remaining
	filled := 0
	if total > 0 {
		filled = (elapsed * barWidth) / total
	}
	if filled > barWidth {
		filled = barWidth
	}
	empty := barWidth - filled

	bar := color + strings.Repeat("█", filled) + dim + strings.Repeat("░", empty) + reset

	mins := remaining / 60
	secs := remaining % 60
	timeStr := fmt.Sprintf("%02d:%02d", mins, secs)

	pct := 0
	if total > 0 {
		pct = (elapsed * 100) / total
	}

	fmt.Println()
	fmt.Printf("  %s%s %s%s\n", bold, color, label, reset)
	fmt.Println()
	fmt.Printf("  %s  %s%s%s  %d%%\n", bar, bold, timeStr, reset, pct)
	fmt.Println()

	if paused {
		fmt.Printf("  %s%s⏸  PAUSED%s\n", bold, yellow, reset)
	} else {
		fmt.Printf("  %s[space]%s pause  %s[q]%s quit\n", dim, reset, dim, reset)
	}
	fmt.Println()
}

func renderDone() {
	clearAndReset()
	fmt.Println()
	fmt.Printf("  %s%s✓  All done! Nice work.%s\n", bold, green, reset)
	fmt.Println()
	fmt.Print("\a")
}

// --- Inline display (single line, no screen clearing) ---

func renderInline(label string, remaining int, total int, paused bool) {
	color := green
	if label == "Break" {
		color = cyan
	}

	elapsed := total - remaining
	filled := 0
	if total > 0 {
		filled = (elapsed * inlineBarWidth) / total
	}
	if filled > inlineBarWidth {
		filled = inlineBarWidth
	}
	empty := inlineBarWidth - filled

	bar := color + strings.Repeat("█", filled) + dim + strings.Repeat("░", empty) + reset

	mins := remaining / 60
	secs := remaining % 60

	status := ""
	if paused {
		status = yellow + " ⏸" + reset
	}

	fmt.Printf("\r  %s%s%s %s %02d:%02d%s  ", color, label, reset, bar, mins, secs, status)
}

func renderInlineDone() {
	fmt.Printf("\r  %s%s✓  All done! Nice work.%s\n", bold, green, reset)
	fmt.Print("\a")
}
