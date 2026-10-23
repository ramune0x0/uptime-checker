package main

import (
	"fmt"
)

func splitSeconds(sec int) (second, minutes int) {
	minutes = sec / 60
	second = sec % 60
	return second, minutes
}

func main() {
	sec := 125
	fmt.Println(splitSeconds(sec))
}
