package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Fprintln(os.Stderr, "usage: go run ./cmd/check <URL> ...")
		os.Exit(2)
	}
	for i, v := range os.Args[1:] {
		fmt.Printf("args[%d] -> %s\n", i+1, v)
	}
}
