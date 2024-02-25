package main

import (
	"fmt"
	"sync"
)

func main() {
	unbuffered := make(chan string)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() { // читатель
		defer wg.Done()
		for {
			fmt.Println("write to channel")
			v, ok := <-unbuffered
			if !ok {
				fmt.Println("stop reader")
				return
			}
			fmt.Printf("reader: %s\n", v)
		}
	}()
	wg.Add(1)
	go func() { // писатель
		defer wg.Done()
		for i := 0; i <= 9; i++ {
			unbuffered <- fmt.Sprintf("Hello #%d", i)
		}
		// Не забываем закрывать канал.
		// Если этого не сделать возможны два исхода:
		// 1) канал будет уничтожен сборщиком мусора (нам повезло и в этом примере это наш сценарий)
		// 2) утечет память
		close(unbuffered)
	}()
	wg.Wait()
}
