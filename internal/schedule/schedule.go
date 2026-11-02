package schedule

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ramune0x0/uptime-checker/internal/check"
	"github.com/ramune0x0/uptime-checker/internal/target"
)

func checkOnce(ctx context.Context, url string, timeout time.Duration) {
	r := check.URL(ctx, url, timeout)
	fmt.Printf("%v\n", r)
}

func runTarget(ctx context.Context, t target.Target, timeout time.Duration) {
	checkOnce(ctx, t.URL, timeout)
	ticker := time.NewTicker(t.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			checkOnce(ctx, t.URL, timeout)
		}
	}
}

func Run(ctx context.Context, targets []target.Target, timeout time.Duration) {
	var wg sync.WaitGroup
	for _, t := range targets {
		wg.Add(1)
		go func() {
			defer wg.Done()
			runTarget(ctx, t, timeout)
		}()
	}
	wg.Wait()
}
