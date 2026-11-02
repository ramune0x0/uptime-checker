package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/ramune0x0/uptime-checker/internal/check"
)

func main() {
	timeout := flag.Duration("timeout", 5*time.Second, "timeout per request (e.g. 5s, 500ms)")
	flag.Parse()
	if len(flag.Args()) == 0 {
		fmt.Fprintf(os.Stderr, "usage: %s [-timeout 5s] <URL> ...\n", os.Args[0])
		os.Exit(2)
	}

	hasDown := false
	for _, url := range flag.Args() {
		r := check.URL(context.Background(), url, *timeout)
		if r.Err != nil {
			fmt.Printf("DOWN %s -> %v\n", url, r.Err)
			hasDown = true
			continue
		}
		if r.IsUp() {
			fmt.Printf("UP %s -> %d, %dms\n", url, r.StatusCode, r.Elapsed.Milliseconds())
		} else {
			fmt.Printf("DOWN %s -> %d, %dms\n", url, r.StatusCode, r.Elapsed.Milliseconds())
			hasDown = true
		}
	}
	if hasDown {
		os.Exit(1)
	}
}
