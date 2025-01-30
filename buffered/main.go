package main

import "fmt"

func main() {
	bufChannel := make(chan string, 3)
	chars := []string{"a", "b", "c"}
	for _, c := range chars {
		bufChannel <- c
	}

	close(bufChannel)

	// Iterate even after closing the channel. This is async.
	for res := range bufChannel {
		fmt.Println(res)
	}
}
