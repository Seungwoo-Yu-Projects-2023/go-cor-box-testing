package main

import "sync"

type CoRBox[T any] struct {
	value   T
	changes []func(waitGroup *sync.WaitGroup, value *T)
	mutex   sync.Mutex
}

func New[T any](value T) CoRBox[T] {
	return CoRBox[T]{
		value: value,
		mutex: sync.Mutex{},
	}
}

func (c *CoRBox[T]) Update(f func(waitGroup *sync.WaitGroup, value *T)) {
	c.mutex.Lock()
	c.changes = append(c.changes, f)
	defer c.mutex.Unlock()
}

func (c *CoRBox[T]) Get() T {
	c.mutex.Lock()

	var waitGroup = sync.WaitGroup{}
	var result = c.value
	var changes = c.changes
	c.changes = []func(waitGroup *sync.WaitGroup, value *T){}

	c.mutex.Unlock()

	for len(changes) > 0 {
		waitGroup.Add(1)
		changes[0](&waitGroup, &result)
		changes = changes[1:]
		waitGroup.Wait()
	}

	c.mutex.Lock()

	c.value = result

	defer c.mutex.Unlock()

	return result
}
