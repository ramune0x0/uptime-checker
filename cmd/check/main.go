package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"
)

func checkURL(url string) bool {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		fmt.Printf("DOWN %s -> %v\n", url, err)
		return false
	}

	start := time.Now()
	resp, err := http.DefaultClient.Do(req)
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
