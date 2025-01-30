package main

import "fmt"

func main() {
	myChannel := make(chan string)
	anotherChannel := make(chan string)
	go func() {
		myChannel <- "data"
	}()

	go func() {
		anotherChannel <- "helo"
	}()

	select {
	case msgMyChannnel := <-myChannel:
		fmt.Println(msgMyChannnel)
	case msgAnotherChan := <-anotherChannel:
		fmt.Println(msgAnotherChan)
	}
}
