package main

import (
	"fmt"
	"sync"
	"time"

	_ "github.com/lib/pq"
)

var counter *Counter

// Counter implements thread safe counter
type Counter struct {
	sync.RWMutex
	count int
}

func (c *Counter) Increment() {
	c.Lock()
	c.count++
	c.Unlock()
}

func (c *Counter) Decrement() {
	c.Lock()
	c.count--
	c.Unlock()
}

func (c *Counter) Count() int {
	c.RLock()
	defer c.RUnlock()
	return c.count
}

func main() {
	fmt.Println("Singleton v1: Using pointer")
	counter = getCounter()

	// Looks not as good as `for i in range(10)` in Python
	// still it's fun to know that there is similar possibility)
	for range [10]struct{}{} {
		go func() {
			counter.Increment()
		}()
	}

	newCounter := getCounter()

	for range [5]struct{}{} {
		go func() {
			newCounter.Decrement()
		}()
	}

	time.Sleep(time.Second * 3)
	// In ideal world WaitGroup shoud always be used for such cases
	// I mean to wait till all the goroutines return
	fmt.Println("Count 1:", counter.Count())
	fmt.Println("Count 2:", newCounter.Count())
	// Well we achieved our goal - a single instance of counter
	fmt.Println("Adresses are equal:", counter == newCounter)

	fmt.Println("Singleton v2: Sync package")
	fmt.Println("Singleton v3: Using channels")
}

func getCounter() *Counter {
	if counter != nil {
		return counter
	}

	return new(Counter)
}
