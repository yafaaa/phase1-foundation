package main

import "testing"

func BenchmarkFib(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// Clear or reuse memo to benchmark memoized lookup speed
		SafeFib(90)
	}
}