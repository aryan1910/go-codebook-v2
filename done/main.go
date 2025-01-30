package main

import (
	"fmt"
	"time"
)

// done as read only channel
func doWork(done <-chan bool) {
	for {
		select {
		case <-done:
			time.Sleep(1 * time.Second)
			fmt.Println("closing")
			time.Sleep(1 * time.Second)
			return
		default:
			fmt.Println("hello")
		}
	}
}
func main() {
	done := make(chan bool)

	go doWork(done)

	time.Sleep(2 * time.Second)
	close(done)
	time.Sleep(10 * time.Second)
}
