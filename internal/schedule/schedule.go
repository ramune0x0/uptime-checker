package schedule

import (
	"context"
	"sync"
	"time"

	"github.com/ramune0x0/uptime-checker/internal/check"
	"github.com/ramune0x0/uptime-checker/internal/target"
)

func checkOnce(ctx context.Context, t target.Target, timeout time.Duration, handle func(target.Target, check.Result)) {
	r := check.URL(ctx, t.URL, timeout)
	handle(t, r)
}

func runTarget(ctx context.Context, t target.Target, timeout time.Duration, handle func(target.Target, check.Result)) {
	checkOnce(ctx, t, timeout, handle)
	ticker := time.NewTicker(t.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			checkOnce(ctx, t, timeout, handle)
		}
	}
}

func Run(ctx context.Context, targets []target.Target, timeout time.Duration, handle func(target.Target, check.Result)) {
	var wg sync.WaitGroup
	for _, t := range targets {
		wg.Add(1)
		go func() {
			defer wg.Done()
			runTarget(ctx, t, timeout, handle)
		}()
	}
	wg.Wait()
}
