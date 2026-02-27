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

	if len(os.Args) < 2 {
		fmt.Println("Usage: pomodoro start [-w minutes] [-b minutes]")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "start":
		startCmd.Parse(os.Args[2:])
		runTimer(*work, *brk)
	case "--version", "-v":
		fmt.Println("pomodoro v0.2.0")
	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		fmt.Println("Usage: pomodoro start [-w minutes] [-b minutes]")
		os.Exit(1)
	}
}
