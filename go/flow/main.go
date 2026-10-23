package main

import "fmt"

func classify(code int) string {
	if code >= 200 && code < 300 {
		return "ok"
	}
	return "ng"
}

func describe(code int) string {
	switch code {
	case 200:
		return "OK"
	case 301, 302:
		return "Redirect"
	case 404:
		return "Not Found"
	default:
		return "Unknown"
	}
}

func level(latencyMs int) string {
	switch {
	case latencyMs < 100:
		return "fast"
	case latencyMs < 500:
		return "slow"
	default:
		return "timeout"
	}
}

func forLoops() {
	for i := 0; i < 3; i++ {
		fmt.Println("a")
	}

	retries := 0
	for retries < 3 {
		fmt.Println("retry")
		retries++
	}

	count := 0
	for {
		count++
		fmt.Println("try", count)
		if count == 5 {
			break
		}
	}
}

func main() {
	if result := classify(503); result == "ng" {
		fmt.Println("alert")
	}
}
