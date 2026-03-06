package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	console "thathost.com/golang/gotest/console/internal"
)

func handleCtrlC() {
	cleanup()
}

func setupCtrlCHandler() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		handleCtrlC()
		os.Exit(0)
	}()
}

func cleanup() {
	console.Clear()
	console.SetCursorVisible(true)
}

func bounce() {
	setupCtrlCHandler()
	console.SetCursorVisible(false)
	console.SetBackground(console.Blue)
	console.SetForeground(console.Yellow)
	console.Clear()
	w, h := console.GetSizeWH()
	x := 7
	y := 1
	dx := 1
	dy := 1
	for true {
		start := time.Now()
		console.MoveToXY(x, y)
		fmt.Print(" ")
		x += dx
		if x < 0 || x >= w {
			dx = -dx
			x += dx
		}
		y += dy
		if y < 0 || y >= h {
			dy = -dy
			y += dy
		}
		console.MoveToXY(x, y)
		fmt.Print("\u2b24")
		timeSpent := time.Since(start)
		time.Sleep(30*time.Millisecond - timeSpent)
	}
}

func main() {
	bounce()
}
