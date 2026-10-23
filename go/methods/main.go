package main

import "fmt"

type Fib struct {
	a, b int
}

// ここに Next メソッドを書く
func (f *Fib) Next() int {
	tmp := f.a
	f.a, f.b = f.b, f.a+f.b
	return tmp
}

func main() {
	f := Fib{a: 0, b: 1}
	for i := 0; i < 10; i++ {
		fmt.Println(f.Next()) // 0 1 1 2 3 5 8 13 21 34
	}
}
