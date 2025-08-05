package main

import (
	"fmt"
	"sync"
)

func main() {
	var mu sync.Mutex
	counter := 0
	var wg sync.WaitGroup

	numGoroutines := 1000000

	wg.Add(numGoroutines)

	increment := func() {
		defer wg.Done()
		mu.Lock()
		counter++
		mu.Unlock()
	}

	for range numGoroutines {
		go increment()
	}

	wg.Wait()
	fmt.Printf("Final counter: %d\n", counter)
}
