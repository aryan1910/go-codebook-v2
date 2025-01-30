package main

import (
	"bytes"
	"fmt"
	"runtime"
	"runtime/pprof"
	"time"
)

func reportGoroutines() {
	var buf bytes.Buffer
	pprof.Lookup("goroutine").WriteTo(&buf, 1)
	fmt.Println(buf.String())
}

func main() {

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	// Run the reporting loop for 30 seconds
	go func() {
		for i := 0; i < 50; i++ {
			<-ticker.C
			fmt.Printf("=== Goroutine Report at %s ===\n", time.Now().Format(time.RFC3339))
			fmt.Printf("Number of goroutines: %d\n", runtime.NumGoroutine())
			// reportGoroutines()
		}
	}()

	var or func(channels ...<-chan interface{}) <-chan interface{}
	or = func(channels ...<-chan interface{}) <-chan interface{} {
		fmt.Println("Input to or", channels, len(channels))
		switch len(channels) {
		case 0:
			fmt.Println("Or 0")
			return nil
		case 1:
			fmt.Println("Or 1")
			return channels[0]
		}
		orDone := make(chan interface{})
		go func() {
			defer close(orDone)
			switch len(channels) {
			case 2:
				select {
				case <-channels[0]:
				case <-channels[1]:
				}
			default:
				select {
				case <-channels[0]:
				case <-channels[1]:
				case <-channels[2]:
				case <-or(append(channels[3:], orDone)...):
				}
			}
		}()
		return orDone
	}

	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}

	start := time.Now()
	go func() {
		<-or(
			sig(2*time.Hour),
			sig(5*time.Minute),
			sig(3*time.Second),
			sig(1*time.Hour),
			sig(1*time.Minute),
			sig(10*time.Second),
			sig(20*time.Second),
		)
	}()
	fmt.Println("done after", time.Since(start))
	time.Sleep(2 * time.Minute)
}
