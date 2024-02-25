package main

import (
	"fmt"
	"time"
)

func main() {
	var ch chan string

	// Из nil-канала можно читать и в него можно писать.
	// Но горутины зависнут и ни deadlock detector ни race detector не помогут нам узнать об этом.
	go func() {
		fmt.Println("write to nil-channel")
		ch <- "he"
		fmt.Println("write done")
	}()

	go func() {
		fmt.Println("read from nil-channel")
		<-ch
		fmt.Println("read done")
	}()

	time.Sleep(5 * time.Second)
	fmt.Println("done")
}
