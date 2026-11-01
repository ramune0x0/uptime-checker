package main

import (
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
		fmt.Fprintln(os.Stderr, "usage: go run ./cmd/check <URL> ...")
		os.Exit(2)
	}

	hasDown := false
	for _, v := range flag.Args() {
		r := check.URL(v, *timeout)
		if r.Err != nil {
			fmt.Printf("DOWN %s -> %v\n", v, r.Err)
			hasDown = true
			continue
		}
		if r.IsUp() {
			fmt.Printf("UP %s -> %d, %dms\n", v, r.StatusCode, r.Elapsed.Milliseconds())
		} else {
			fmt.Printf("DOWN %s -> %d, %dms\n", v, r.StatusCode, r.Elapsed.Milliseconds())
			hasDown = true
		}
	}
	if hasDown {
		os.Exit(1)
	}
}
