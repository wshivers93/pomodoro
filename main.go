package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	startCmd := flag.NewFlagSet("start", flag.ExitOnError)
	work := startCmd.Int("w", 25, "work duration in minutes")
	brk := startCmd.Int("b", 5, "break duration in minutes")
	mode := startCmd.String("mode", "auto", "display mode: auto, tmux, window, inline")

	runCmd := flag.NewFlagSet("_run", flag.ExitOnError)
	runWork := runCmd.Int("w", 25, "work duration in minutes")
	runBreak := runCmd.Int("b", 5, "break duration in minutes")
	runMode := runCmd.String("mode", "fullscreen", "render mode: fullscreen, inline")

	if len(os.Args) < 2 {
		fmt.Println("Usage: pomodoro start [-w minutes] [-b minutes] [--mode auto|tmux|window|inline]")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "start":
		startCmd.Parse(os.Args[2:])
		launch(*work, *brk, *mode)
	case "_run":
		runCmd.Parse(os.Args[2:])
		runTimer(*runWork, *runBreak, *runMode == "inline")
	case "--version", "-v":
		fmt.Println("pomodoro v0.2.0")
	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		fmt.Println("Usage: pomodoro start [-w minutes] [-b minutes] [--mode auto|tmux|window|inline]")
		os.Exit(1)
	}
}
