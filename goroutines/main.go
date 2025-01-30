package main

import (
	"fmt"
	"time"
)

func someFunc(x int) {
	fmt.Println(x)
}
func main() {
	go someFunc(1)
	go someFunc(2)
	go someFunc(3)

	// Follows fork join model
	// Child process is forked from the main process
	// At some point, child process joins the main process
	time.Sleep(2 * time.Second)

	fmt.Println("Hi")
}
