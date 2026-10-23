package main

import (
	"fmt"
)

func backoff(attempts int) float64 {
	waitSec := 1.0
	// for ; attempts > 0; attempts-- {
	for waitSec < 60 {
		fmt.Println(waitSec)
		waitSec = waitSec * 2
	}
	return waitSec
}

func main() {
	backoff(10)
}
