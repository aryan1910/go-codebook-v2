package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

func generate(done <-chan int, fn func() int) <-chan int {
	channel := make(chan int)
	go func() {
		defer close(channel)
		for {
			select {
			case <-done:
				return
			case channel <- fn():
			}
		}

	}()
	return channel
}

func take(done <-chan int, stream <-chan int, n int) <-chan int {
	taken := make(chan int)
	go func() {
		defer close(taken)
		for i := 0; i < n; i++ {
			select {
			case <-done:
				return
			case taken <- <-stream:
			}
		}
	}()
	return taken
}

func primeFinder(done <-chan int, randIntStream <-chan int) <-chan int {
	isPrime := func(num int) bool {
		if num < 2 {
			return false // Numbers less than 2 are not prime.
		}
		for i := 2; i*i <= num; i++ { // Check divisors up to the square root of num.
			if num%i == 0 {
				return false // Divisible by i, so it's not prime.
			}
		}
		return true // No divisors found, so it's prime.
	}
	primes := make(chan int)
	go func() {
		defer close(primes)
		for {
			select {
			case <-done:
				return
			case num := <-randIntStream:
				if isPrime(num) {
					primes <- num
				}
			}
		}
	}()
	return primes
}

func slow() {
	done := make(chan int)
	defer close(done)
	randNumFetcher := func() int { return rand.Int() }
	randNumStream := generate(done, randNumFetcher)
	primeStream := primeFinder(done, randNumStream)
	for i := range take(done, primeStream, 10) {
		fmt.Println(i)
	}
}

func fanIn(done <-chan int, channels ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	fannedInStream := make(chan int)

	transfer := func(c <-chan int) {
		defer wg.Done()
		for i := range c {
			select {
			case <-done:
				return
			case fannedInStream <- i:
			}
		}
	}

	for _, c := range channels {
		wg.Add(1)
		go transfer(c)
	}

	go func() {
		wg.Wait()
		close(fannedInStream)
	}()

	return fannedInStream
}

func fast() {
	done := make(chan int)
	defer close(done)
	randNumFetcher := func() int { return rand.Int() }
	randNumStream := generate(done, randNumFetcher)

	cpuCount := runtime.NumCPU()
	fmt.Println("Num CPUs", cpuCount)

	// Fan out
	primeFinderChannels := make([]<-chan int, cpuCount)
	for i := 0; i < cpuCount; i++ {
		primeFinderChannels[i] = primeFinder(done, randNumStream)
	}

	fanInStream := fanIn(done, primeFinderChannels...)
	for i := range take(done, fanInStream, 10) {
		fmt.Println(i)
	}
}
func main() {
	start := time.Now()
	fast()
	fmt.Println(time.Since(start))
}
