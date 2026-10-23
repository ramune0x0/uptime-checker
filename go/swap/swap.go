package main

import (
	"fmt"
)

func swap(a, b string) (string, string) {
	return b, a
}

func main() {
	a := "World"
	b := "Hello"
	fmt.Println("before: ", a, b)
	fmt.Println(swap(a, b))
}
