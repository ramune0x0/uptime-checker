package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

func checkURL(url string) bool {
	start := time.Now()
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("DOWN %s -> %v\n", url, err)
		return false
	}
	defer resp.Body.Close()
	elapsed := time.Since(start)
	if resp.StatusCode >= http.StatusOK && resp.StatusCode < http.StatusBadRequest {
		fmt.Printf("UP %s -> %d, %dms\n", url, resp.StatusCode, elapsed.Milliseconds())
		return true
	} else {
		fmt.Printf("DOWN %s -> %d, %dms\n", url, resp.StatusCode, elapsed.Milliseconds())
		return false
	}
}

func main() {
	if len(os.Args) == 1 {
		fmt.Fprintln(os.Stderr, "usage: go run ./cmd/check <URL> ...")
		os.Exit(2)
	}

	hasDown := false
	for _, v := range os.Args[1:] {
		if !checkURL(v) {
			hasDown = true
		}
	}
	if hasDown {
		os.Exit(1)
	}
}
