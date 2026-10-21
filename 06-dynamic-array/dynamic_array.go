package main

type DynamicArray struct {
	data []int
}

func (d *DynamicArray) Len() int {
	// TODO: return current length of d.data
	return len(d.data)
}

func (d *DynamicArray) Cap() int {
	// TODO: return current capacity of d.data
	return cap(d.data)
}

func (d *DynamicArray) Get(index int) int {
	// TODO: return element at index
	return d.data[index]
}

func (d *DynamicArray) Set(index int, val int) {
	// TODO: set element at index to val
	d.data[index] = val
}

func (d *DynamicArray) Append(val int) {
	// TODO: Grow d.data by 1 without built-in append (just like appendInt!):
	// 1. If len(d.data)+1 <= cap(d.data), stretch d.data = d.data[:len(d.data)+1]
	// 2. Else, double capacity (or set to 1 if cap is 0), make() a new slice, copy(newSlice, d.data), and set d.data = newSlice
	// 3. Put 'val' into the last index of d.data
	needed := len(d.data) + 1

	if needed <= cap(d.data){
		d.data = d.data[:needed]
	} else {
		var z []int
		x := needed
		if needed < 2*cap(d.data){
			
			x = 2*cap(d.data)
		}

		z = make([]int, needed, x)
		
		copy(z, d.data)
		d.data = z
	}
	d.data[len(d.data)-1] = val
}