package main

import (
	"fmt"
)

func main() {
	ch := make(chan string)
	fmt.Println("write to ch: ")
	ch <- "a" // deadlock тут
	res := <-ch
	fmt.Println("get res: ", res)
}
