package main

import (
	"fmt"
	"sync"
)

func main() {
	wg := sync.WaitGroup{}
	// Это функция, которая пишет в канал.
	writer := func(wg *sync.WaitGroup) <-chan string {
		out := make(chan string)
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i <= 9; i++ {
				out <- fmt.Sprintf("Hello #%d", i)
			}
			close(out)
		}()
		return out
	}(&wg)

	wg.Add(1)
	// А это функция, которая читает из канала.
	go func(in <-chan string) {
		defer wg.Done()
		for v := range in {
			fmt.Println(v)
		}
		fmt.Println("stop reader")
	}(writer)
	wg.Wait()
}
