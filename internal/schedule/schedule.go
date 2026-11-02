package schedule

import (
	"context"
	"fmt"
	"time"

	"github.com/ramune0x0/uptime-checker/internal/check"
	"github.com/ramune0x0/uptime-checker/internal/target"
)

func checkOnce(ctx context.Context, url string, timeout time.Duration) {
	r := check.URL(ctx, url, timeout)
	fmt.Printf("%v\n", r)
}

func RunTarget(ctx context.Context, t target.Target, timeout time.Duration) {
	checkOnce(ctx, t.URL, timeout)
	ticker := time.NewTicker(t.Interval)
	defer ticker.Stop()

	for {
		<-ticker.C
		checkOnce(ctx, t.URL, timeout)
	}
}
