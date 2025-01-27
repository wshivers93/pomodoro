package cmd

import (
	"fmt"
	"time"
	//"strconv"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts the pomodoro timer",
	Run: func(cmd *cobra.Command, args []string) {
	      fmt.Printf("Inside start Run with args: %v\n", args)
    },
}

func init() {
	startCmd.Flags().IntP("hours", "h", 1, "The hours value for the timer")
	startCmd.Flags().IntP("minutes", "m", 0, "The minutes value for the timer")
	//durationStr, _ := time.ParseDuration(strconv.Itoa(*hours) + "h" + strconv.Itoa(*minutes) + "m")

	rootCmd.AddCommand(startCmd)
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
