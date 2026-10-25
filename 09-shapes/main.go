package main

import (
	"fmt"
	"math"
	"sort"
)

// 1. Define the Shape interface
type Shape interface {
	Area() float64
	Perimeter() float64
}

// 2. Circle implementation
type Circle struct {
	Radius float64
}

func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

func (c Circle) Perimeter() float64 {
	return 2 * math.Pi * c.Radius
}

// 3. TODO: Implement Rectangle struct (Width, Height float64) with Area() and Perimeter()
type Rectangle struct {
	Width, Height float64
}

func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

func (r Rectangle) Perimeter() float64 {
	return 2 * r.Height + 2* r.Width
}

// 4. TODO: Implement Triangle struct (e.g., Base, Height, SideA, SideB, SideC float64)
//    - Area() = 0.5 * Base * Height
//    - Perimeter() = SideA + SideB + SideC
type Triangle struct{
	Base, Height, SideA, SideB, SideC float64

}

func (t Triangle) Area() float64 {
	return 0.5 * t.Base * t.Height
}

func (t Triangle) Perimeter() float64 {
	return t.SideA + t.SideB + t.SideC
}
// 5. Sort function using the Shape interface in its signature
func sortShapesByArea(shapes []Shape) {
	// TODO: Use sort.Slice(shapes, func(i, j int) bool { ... })
	// to sort from smallest Area() to largest Area()
	sort.Slice(shapes, func(i, j int) bool {
		return shapes[i].Area() < shapes[j].Area()
	})
}

func main() {
	shapes := []Shape{
		Circle{Radius: 5},
		// TODO: Add a Rectangle and a Triangle here
		Rectangle{Width: 5, Height: 6},
		Triangle{Base: 7, Height: 6, SideA: 4, SideB: 4, SideC: 6},
	}

	sortShapesByArea(shapes)

	for _, s := range shapes {
		fmt.Printf("%T -> Area: %.2f, Perimeter: %.2f\n", s, s.Area(), s.Perimeter())
	}
}