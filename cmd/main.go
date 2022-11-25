package main

import (
	"sync"
	"time"
)

func main() {
	var mutex sync.Mutex

	go func() {
		mutex.Lock()

		defer mutex.Unlock()

		mutex.Lock()
	}()

	time.Sleep(time.Second * 5)
}
