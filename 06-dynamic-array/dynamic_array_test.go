package main

import "testing"

func TestDynamicArray(t *testing.T) {
	arr := &DynamicArray{}

	// Test Append, Len, and Cap growth
	arr.Append(10)
	arr.Append(20)
	arr.Append(30) // Should trigger capacity doubling (cap becomes 4)

	if arr.Len() != 3 {
		t.Errorf("Expected Len 3, got %d", arr.Len())
	}
	if arr.Cap() != 4 {
		t.Errorf("Expected Cap 4, got %d", arr.Cap())
	}

	// Test Get and Set
	if arr.Get(1) != 20 {
		t.Errorf("Expected index 1 to be 20, got %d", arr.Get(1))
	}

	arr.Set(1, 99)
	if arr.Get(1) != 99 {
		t.Errorf("Expected index 1 to be 99 after Set, got %d", arr.Get(1))
	}
}