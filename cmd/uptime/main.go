package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
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
	schedule.Run(ctx, targets, *timeout, logResult)
}

func logResult(t target.Target, r check.Result) {
	if r.Err != nil {
		slog.Info("DOWN", "name", t.Name, "url", t.URL, "err", r.Err, "ms", r.Elapsed.Milliseconds())
		return
	}
	if r.IsUp() {
		slog.Info("UP", "name", t.Name, "url", t.URL, "status", r.StatusCode, "ms", r.Elapsed.Milliseconds())
	} else {
		slog.Info("DOWN", "name", t.Name, "url", t.URL, "status", r.StatusCode, "ms", r.Elapsed.Milliseconds())
	}
}
