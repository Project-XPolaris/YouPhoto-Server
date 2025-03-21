package utils

import "time"

func DelayRunAsync(delay time.Duration, f func()) {
	go func() {
		time.Sleep(delay)
		f()
	}()
}
