package main

import (
	"sync/atomic"
	"testing"
)

// ~10% Better
func BenchmarkNew(b *testing.B) {
	var box = New(0)

	for i := 0; i < b.N; i++ {
		box.Update(func(value *int) {
			go func() {
				*value += 1
			}()
		})
	}

	box.Get()
}

// Worse
func BenchmarkNew2(b *testing.B) {
	var value = atomic.Int32{}

	for i := 0; i < b.N; i++ {
		go func() {
			value.Add(1)
		}()
	}
}
