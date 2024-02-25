package main

import (
	"fmt"
	"time"
)

// В нашем примере мы будем выбирать между двумя каналами.
func main() {
	c1 := make(chan string)
	c2 := make(chan string)
	done := make(chan struct{})

	// запишем в каждый канал сообщение в разное время
	go func() {
		time.Sleep(1 * time.Second)
		c1 <- "one"
	}()
	go func() {
		time.Sleep(1500 * time.Millisecond)
		c2 <- "two"
	}()

	// Это нужно, чтобы выйти из бесконечного цикла for ниже.
	go func() {
		defer close(done)
		time.Sleep(2 * time.Second)
	}()

	// Select это как switch, только для каналов.
	// Мы запускаем его в цикле, так как первое пришедшее сообщение из канала распечатается и мы покинем select.
	for {
		select {
		case msg1 := <-c1:
			fmt.Println("received", msg1)
		case msg2 := <-c2:
			fmt.Println("received", msg2)
		case <-done:
			return
		}
	}
}
