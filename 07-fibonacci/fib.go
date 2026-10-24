package main

import "fmt"

var memo = make(map[int]int)

// SafeFib wraps fib() with defer/recover so an overflow returns -1 instead of crashing
func SafeFib(n int) (result int) {
	// TODO 1: Add a defer func() that calls recover().
	// If recover() != nil, print the panic message and set result = -1.
	defer func(){

		if r := recover(); r != nil {
			fmt.Println("Recovered: ", r)
			result = -1
		}
	}()

	return fib(n)
}

func fib(n int) int {
	if n <= 1 {
		return n
	}

	// TODO 2: Check if 'n' is already in 'memo' using val, ok := memo[n]. If ok, return val.
	if val, ok := memo[n]; ok{
		return val
	}
	a := fib(n - 1)
	b := fib(n - 2)
	sum := a + b

	// TODO 3: If sum < 0 (overflow occurred!), call panic("fibonacci integer overflow")
	if sum < 0{
		panic("fibonacci integer overflow")
	}
	memo[n] = sum
	return sum
}

func main() {
	fmt.Println("Fib(50):", SafeFib(50))
	fmt.Println("Fib(92):", SafeFib(92))
	fmt.Println("Fib(93) [Overflow]:", SafeFib(93))
}