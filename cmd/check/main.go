package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

func checkURL(url string) {
	start := time.Now()
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("DOWN %s -> %v\n", url, err)
		return
	}
	defer resp.Body.Close()
	elapsed := time.Since(start)
	if resp.StatusCode >= http.StatusOK && resp.StatusCode < http.StatusBadRequest {
		fmt.Printf("UP %s -> %d, %dms\n", url, resp.StatusCode, elapsed.Milliseconds())
	} else {
		fmt.Printf("DOWN %s -> %d, %dms\n", url, resp.StatusCode, elapsed.Milliseconds())
	}
}

func main() {
	if len(os.Args) == 1 {
		fmt.Fprintln(os.Stderr, "usage: go run ./cmd/check <URL> ...")
		os.Exit(2)
	}

	for _, v := range os.Args[1:] {
		checkURL(v)
	}
}
