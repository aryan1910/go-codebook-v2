package main

import (
	"fmt"
	"math/rand"
	"runtime"
	// "sync"
	"time"
)

func repeatFunc[T any, K any](done <-chan K, fn func() T) <-chan T {
	stream := make(chan T)
	go func() {
		defer close(stream)
		for {
			select {
			case <-done:
				return
			case stream <- fn():
			}
		}
	}()
	return stream
}

func take[T any, K any](done <-chan K, stream <-chan T, n int) <-chan T {
	taken := make(chan T)
	go func() {
		defer close(taken)
		for i := 0; i < n; i++ {
			select {
			case <-done:
				return
			// Put value on stream to taken chan
			case taken <- <-stream:
			}
		}
	}()
	return taken
}

func primeFinder(done <-chan int, randIntStream <-chan int) <-chan int {
	isPrime := func(randomInt int) bool {
		for i := randomInt - 1; i > 1; i-- {
			if randomInt%i == 0 {
				return false
			}
		}
		return true
	}

	primes := make(chan int)

	go func() {
		defer close(primes)
		for {
			select {
			case <-done:
				return
			case randomInt := <-randIntStream:
				if isPrime(randomInt) {
					primes <- randomInt
				}
			}
		}
	}()
	return primes
}

func fanIn[T any](done <-chan int, channels ...<-chan T) <-chan T {
	// var wg sync.WaitGroup

	fannedInStream := make(chan T)

	transfer := func(c <-chan T) {
		// defer wg.Done()
		for i := range c {
			select {
			case <-done:
				return
			case fannedInStream <- i:
			}
		}
	}

	for _, val := range channels {
		// wg.Add(1)
		go transfer(val)
	}

	// go func() {
	// 	wg.Wait()
	// 	close(fannedInStream)
	// }()
	return fannedInStream
}

func main() {
	done := make(chan int)

	start := time.Now()
	fn := func() int {
		return rand.Intn(200000000)
	}

	randIntStream := repeatFunc(done, fn)
	// primeStream := primeFinder(done, randIntStream)
	// // naive slow way
	// for rn := range take(done, primeStream, 10) {
	// 	fmt.Println(rn)
	// }
	// fmt.Println(time.Since(start).Seconds())

	// start = time.Now()
	cpuCount := runtime.NumCPU()
	fmt.Println("Number of cpus", cpuCount)
	primeStreams := make([]<-chan int, cpuCount)
	for i := 0; i < cpuCount; i++ {
		primeStreams[i] = primeFinder(done, randIntStream)
	}
	fannedInStream := fanIn(done, primeStreams...)
	for rn := range take(done, fannedInStream, 10) {
		fmt.Println(rn)
	}
	fmt.Println(time.Since(start).Seconds())
	close(done)
}
