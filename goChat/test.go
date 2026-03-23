package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	var rw sync.RWMutex

	observer := func(i int) {
		defer wg.Done()
		rw.RLock()
		fmt.Println("Reading started", i)
		time.Sleep(400 * time.Millisecond)
		fmt.Println("Reading ends", i)
		rw.RUnlock()
	}

	writter := func(i int) {
		defer wg.Done()
		defer rw.Unlock()
		rw.Lock()
		fmt.Println("Writting started", i)
		time.Sleep(400 * time.Millisecond)
		fmt.Println("Writting ended", i)
	}

	for i := 0; i < 100; i++ {
		wg.Add(1)
		if i%25 == 0 {
			go writter(i)
		} else {
			go observer(i)
		}
	}
	wg.Wait()
	fmt.Println("End")
}
