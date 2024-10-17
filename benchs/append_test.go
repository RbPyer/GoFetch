package main

import "testing"


func BenchmarkAppend1(b *testing.B) {
	var sl []int
	for i := 0; i < b.N; i++ {
		sl = append(sl, i)
	}
}


func BenchmarkAppend2(b *testing.B) {
	sl := make([]int, 0, b.N)
	for i := 0; i < b.N; i++ {
		sl = append(sl, i)
	}
}