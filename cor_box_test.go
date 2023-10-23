package main

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"sync/atomic"
	"testing"
)

func TestNew(t *testing.T) {
	var box = New(int32(0))

	for i := 0; i < 10000; i++ {
		box.Update(func(waitGroup *sync.WaitGroup, value *int32) {
			go func() {
				*value += 1
				waitGroup.Done()
			}()
		})
	}

	assert.Equal(t, len(box.changes), 10000)
	assert.Equal(t, box.Get(), int32(10000))
}

func TestNew2(t *testing.T) {
	var value = atomic.Int32{}
	var waitGroup = sync.WaitGroup{}

	for i := 0; i < 10000; i++ {
		waitGroup.Add(1)
		go func() {
			value.Add(1)
			waitGroup.Done()
		}()
	}

	waitGroup.Wait()

	assert.Equal(t, value.Load(), int32(10000))
}

func TestNew3(t *testing.T) {
	var value = make(chan int32, 1)
	var waitGroup = sync.WaitGroup{}

	value <- 0

	for i := 0; i < 10000; i++ {
		waitGroup.Add(1)
		go func() {
			var acc = <-value
			acc += 1
			value <- acc
			waitGroup.Done()
		}()
	}

	waitGroup.Wait()

	var result = <-value
	assert.Equal(t, result, int32(10000))
}

// ~10% Better
func BenchmarkNew(b *testing.B) {
	var box = New(0)

	for i := 0; i < b.N; i++ {
		box.Update(func(waitGroup *sync.WaitGroup, value *int) {
			go func() {
				*value += 1
				waitGroup.Done()
			}()
		})
	}

	box.Get()
}

// Worse
func BenchmarkNew2(b *testing.B) {
	var value = atomic.Int32{}
	var waitGroup = sync.WaitGroup{}

	for i := 0; i < b.N; i++ {
		waitGroup.Add(1)
		go func() {
			value.Add(1)
			waitGroup.Done()
		}()
	}

	waitGroup.Wait()
}

func BenchmarkNew3(b *testing.B) {
	var value = make(chan int32, 1)
	var waitGroup = sync.WaitGroup{}

	value <- 0

	for i := 0; i < b.N; i++ {
		waitGroup.Add(1)
		go func() {
			var acc = <-value
			acc += 1
			value <- acc
			waitGroup.Done()
		}()
	}

	waitGroup.Wait()
	<-value
}
