package main

import "fmt"

func main() {
	m := make(map[string]int)
	m["a"] = 1

	v, ok := m["a"]
	fmt.Println(v, ok)

	v2, ok2 := m["b"]
	fmt.Println(v2, ok2)

	m2 := m
	m2["a"] = 100
	fmt.Println(m["a"])
}
