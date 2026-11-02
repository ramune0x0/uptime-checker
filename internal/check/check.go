package check

import (
	"context"
	"net/http"
	"time"
)

type Result struct {
	StatusCode int
	Elapsed    time.Duration
	Err        error
}

func (r Result) IsUp() bool {
	return r.StatusCode >= http.StatusOK && r.StatusCode < http.StatusBadRequest
}

func URL(ctx context.Context, url string, timeout time.Duration) Result {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return Result{Err: err}
	}

	start := time.Now()
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return Result{Err: err}
	}

	defer resp.Body.Close()
	elapsed := time.Since(start)

	return Result{StatusCode: resp.StatusCode, Elapsed: elapsed}
}
