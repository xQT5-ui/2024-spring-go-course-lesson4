package main

import (
	"fmt"
	"runtime"
	"sync"
)

func main() {
	runtime.GOMAXPROCS(1) // это нужно, чтобы код выполнялся только на одном логическом процессоре
	var wg sync.WaitGroup
	wg.Add(2)
	fmt.Println("Starting...")
	go func() {
		defer wg.Done()
		for char := 'a'; char < 'a'+26; char++ {
			//runtime.Gosched()
			fmt.Printf("%c ", char)
		}
	}()
	go func() {
		defer wg.Done()
		for char := 'A'; char < 'A'+26; char++ {
			//runtime.Gosched()
			fmt.Printf("%c ", char)
		}
	}()
	wg.Wait()
	fmt.Println("\nFinished")
}
