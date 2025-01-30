package main

import "fmt"

func sliceToChannel(s []int) <-chan int {
	out := make(chan int, 2)
	go func() {
		// Out is blocked until it is received
		for _, i := range s {
			out <- i
		}
		close(out)
	}()
	return out
}

func square(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			out <- n * n
		}
		close(out)
	}()
	return out
}

func main() {
	a := []int{1, 2, 3, 4}

	sliceChan := sliceToChannel(a)

	sqChan := square(sliceChan)

	for sq := range sqChan {
		fmt.Println(sq)
	}
}
