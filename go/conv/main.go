package main

import (
	"fmt"
)

var (
	x int     = 7
	y int     = 2
	f float64 = 3.9
)

func main() {
	fmt.Println(x, y)
	fmt.Println(float64(x / y))
	fmt.Println(float64(x) / float64(y))
	fmt.Println(int(f))
}
