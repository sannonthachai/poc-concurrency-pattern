package main

import (
	"fmt"
	"sync"
)

func main() {
	in := gen(2, 3, 4)
	concurrency(2, in)
}

func concurrency(numGoroutine int, in <-chan int) {
	test := make([]<-chan int, numGoroutine)

	for i := 0; i < numGoroutine; i++ {
		test[i] = sq(in)
	}

	for n := range merge(test...) {
		fmt.Println(n)
	}
}

func gen(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		for _, n := range nums {
			out <- n
		}
	}()
	return out
}

func sq(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			out <- n * n
		}
	}()
	return out
}

func merge(cs ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	out := make(chan int)

	output := func(c <-chan int) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}
	wg.Add(len(cs))

	for _, c := range cs {
		go output(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}
