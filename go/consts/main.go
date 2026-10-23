package main

import (
	"fmt"
)

func main() {
	const Pi = 3.14
	var f float64 = Pi * 2
	var i int = Pi * 2

	fmt.Println(f, i)

	const Big = 1 << 100
	fmt.Println(Big)
	fmt.Println(Big >> 98)
}
