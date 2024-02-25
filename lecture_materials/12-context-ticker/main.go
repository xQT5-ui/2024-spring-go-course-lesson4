package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	parentCtx := context.Background()
	ctx, cancel := context.WithCancel(parentCtx)
	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop() // не завываем закрывать ticker, иначе он не будет удален из памяти и будет утечка
		for {
			select {
			case <-ctx.Done():
				fmt.Println("ctx done")
				return
			case <-ticker.C:
				fmt.Println(time.Now().Format(time.RFC1123))
			}
		}
	}()

	time.Sleep(time.Second * 10)
	cancel()
}
