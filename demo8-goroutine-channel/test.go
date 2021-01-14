package main

import (
	"fmt"
	"runtime"
	"time"
)

func producer(a int, ch chan<- int) {
	for i := 1; ; i++ {
		ch <- i * a
		// time.Sleep(1 * time.Second)
	}
}

func consumer(ch <-chan int) {
	for v := range ch {
		fmt.Println(v)
	}
}

func main() {
	ch := make(chan int, 10)
	runtime.GOMAXPROCS(1)

	go producer(3, ch)
	go producer(4, ch)
	go consumer(ch)

	time.Sleep(1 * time.Second)
}
