package main

import (
	"fmt"
	"os"
	"syscall"
	"time"
	"unsafe"
)

// Raw terminal mode handling via ioctl
type termios syscall.Termios

func tcgetattr(fd uintptr, t *termios) error {
	_, _, err := syscall.Syscall(syscall.SYS_IOCTL, fd, uintptr(syscall.TIOCGETA), uintptr(unsafe.Pointer(t)))
	if err != 0 {
		return err
	}
	return nil
}

func tcsetattr(fd uintptr, t *termios) error {
	_, _, err := syscall.Syscall(syscall.SYS_IOCTL, fd, uintptr(syscall.TIOCSETA), uintptr(unsafe.Pointer(t)))
	if err != 0 {
		return err
	}
	return nil
}

func enableRawMode(fd uintptr) (*termios, error) {
	var orig termios
	if err := tcgetattr(fd, &orig); err != nil {
		return nil, err
	}

	raw := orig
	raw.Lflag &^= syscall.ECHO | syscall.ICANON | syscall.ISIG
	raw.Cc[syscall.VMIN] = 0
	raw.Cc[syscall.VTIME] = 1 // 100ms timeout on reads

	if err := tcsetattr(fd, &raw); err != nil {
		return nil, err
	}
	return &orig, nil
}

func restoreMode(fd uintptr, orig *termios) {
	tcsetattr(fd, orig)
}

func runTimer(workMins, breakMins int) {
	fd := os.Stdin.Fd()
	orig, err := enableRawMode(fd)
	if err != nil {
		panic(err)
	}
	defer func() {
		restoreMode(fd, orig)
		fmt.Print(showCursor)
	}()

	fmt.Print(hideCursor)

	// Work phase
	if quit := countdown("Work", workMins*60); quit {
		return
	}

	// Bell between phases
	fmt.Print("\a")

	// Break phase
	if quit := countdown("Break", breakMins*60); quit {
		return
	}

	renderDone()
}

func countdown(label string, totalSecs int) bool {
	remaining := totalSecs
	paused := false

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	buf := make([]byte, 1)
	renderDisplay(label, remaining, totalSecs, paused)

	for remaining > 0 {
		// Check for keypress (non-blocking due to VTIME)
		n, _ := os.Stdin.Read(buf)
		if n > 0 {
			switch buf[0] {
			case 'q', 'Q':
				clearAndReset()
				return true
			case ' ':
				paused = !paused
				renderDisplay(label, remaining, totalSecs, paused)
			}
		}

		select {
		case <-ticker.C:
			if !paused {
				remaining--
				renderDisplay(label, remaining, totalSecs, paused)
			}
		default:
		}
	}
	return false
}
