package cmd

import (
	"fmt"
	"time"
	"strconv"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts the pomodoro timer",
	Run: func(cmd *cobra.Command, args []string) {
		runPomodoroTimer()
    },
}
var Hours int
var Minutes int

func init() {
	startCmd.Flags().IntVarP(&Hours, "hours", "r", 1, "The hours value for the timer")
	startCmd.Flags().IntVarP(&Minutes, "minutes", "m", 0, "The minutes value for the timer")

	rootCmd.AddCommand(startCmd)
}

func runPomodoroTimer() {
	timerDuration, _ := time.ParseDuration(strconv.Itoa(Hours) + "h" + strconv.Itoa(Minutes) + "m")
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
