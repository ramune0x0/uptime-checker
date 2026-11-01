package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"
)

func checkURL(url string, timeout time.Duration) bool {
	ctx, cancelFunc := context.WithTimeout(context.Background(), timeout)
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
	timeout := flag.Duration("timeout", 5*time.Second, "timeout per request (e.g. 5s, 500ms)")
	flag.Parse()
	if len(flag.Args()) == 0 {
		fmt.Fprintln(os.Stderr, "usage: go run ./cmd/check <URL> ...")
		os.Exit(2)
	}
	hasDown := false
	for _, v := range flag.Args() {
		if !checkURL(v, *timeout) {
			hasDown = true
		}
	}
	if hasDown {
		os.Exit(1)
	}
}
