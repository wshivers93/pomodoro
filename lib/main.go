package main

import (
	"flag"
	"fmt"
	"strconv"
	"time"
)

func main() {
	// init cli flags
	hours := flag.Int("hours", 1, "The hours value for the timer")
	minutes := flag.Int("minutes", 0, "The minutes values for the timer")
	flag.Parse()
	durationStr, _ := time.ParseDuration(strconv.Itoa(*hours) + "h" + strconv.Itoa(*minutes) + "m")

	runPomodoroTimer(durationStr)
}

func runPomodoroTimer(timerDuration time.Duration) {
	endTime := time.Now().Add(timerDuration)
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	done := make(chan bool)
	go func() {
		time.Sleep(timerDuration)
		done <- true
	}()
	for {
		select {
		case <-done:
			fmt.Println("Done!")
			return
		case <-ticker.C:
			remaining := time.Until(endTime)
			fmt.Printf("Focus for: %s\r", remaining.Round(time.Second))
		}
	}

}
