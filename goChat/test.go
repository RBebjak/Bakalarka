package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	ch := make(chan int, 1)
	go func() {
		for {
			i := rand.Intn(5)
			ch <- i
			time.Sleep(time.Duration(i) * time.Second)
		}
	}()

	for {
		select {
		case num := <-ch:
			fmt.Println(num)

		case <-time.After(3 * time.Second):
			return
		}
	}
}
