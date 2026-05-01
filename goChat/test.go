package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan string, 1)

	writter := func() {
		for i := 0; i < 10; i++ {
			ch <- "write"
		}
		close(ch)
	}
	reader := func() {
		for i := 0; i < 12; i++ {
			val, done := <-ch
			if !done {
				fmt.Println(val)
				fmt.Print("reader ends\n")
			}
			fmt.Println(val)
		}
	}

	go writter()
	for i := 0; i < 2; i++ {
		go reader()
	}

	time.Sleep(200 * time.Microsecond)
	ch <- "end"
}
