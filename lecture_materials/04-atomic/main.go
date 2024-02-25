package main

import (
	"fmt"
	"sync"
)

var counter int64

func main() {
	var wg sync.WaitGroup
	wg.Add(2)
	go incCounter(&wg) //goroutine #1
	go incCounter(&wg) //goroutine #2
	wg.Wait()
	fmt.Println("Final counter: ", counter)
}

func incCounter(wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < 10; i++ {
		counter++
	}
}
