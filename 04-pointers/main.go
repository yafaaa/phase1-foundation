package main

import "fmt"

// doesn't update the orignal balance, update isloated copy of balance.
func updateByValue(val int){
		val += 50
}

//does update balance, use deferencing it can access and modify original balance
func updateByPointer(val *int){
	*val += 50
}

func main() {
	balance := 100
	p1 := &balance
	p2 := p1

	fmt.Println("Initial values (balance, *p1, *p2):", balance, *p1, *p2)

	// Modifying *p2 demonstrates aliasing: balance and *p1 change unexpectedly
	*p2 = 90
	fmt.Println("After *p2 = 90 (balance, *p1, *p2):", balance, *p1, *p2)

	// Part 1 Functions in action
	updateByValue(balance)
	fmt.Println("After updateByValue:", balance) // Still 90

	updateByPointer(p1)
	fmt.Println("After updateByPointer:", balance, *p1, *p2) // All 140
}
