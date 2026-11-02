package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/ramune0x0/uptime-checker/internal/check"
	"github.com/ramune0x0/uptime-checker/internal/schedule"
	"github.com/ramune0x0/uptime-checker/internal/target"
)

func main() {
	timeout := flag.Duration("timeout", 5*time.Second, "timeout per request (e.g. 5s, 500ms)")
	config := flag.String("config", "targets.json", "config file")
	flag.Parse()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	targets, err := target.Load(*config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		os.Exit(1)
	}
	schedule.Run(ctx, targets, *timeout, func(t target.Target, r check.Result) { fmt.Println(t.Name, r) })
}
