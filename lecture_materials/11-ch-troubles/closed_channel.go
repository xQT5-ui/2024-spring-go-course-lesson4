package main

import "fmt"

// Нельзя писать в закрытый канал, но читать из закрытого канала - можно.
func main() {
	ch := make(chan string, 5)
	close(ch)

	ch <- "ok"

	fmt.Printf("ok is: %v\n", <-ch)
}
