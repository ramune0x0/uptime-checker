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

func URL(url string, timeout time.Duration) Result {
	ctx, cancelFunc := context.WithTimeout(context.Background(), timeout)
	defer cancelFunc()

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

	return Result{resp.StatusCode, elapsed, nil}
}
