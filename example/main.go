package main

import (
	"fmt"
	"time"

	regulator "github.com/gallir/goregulator"
)

type elem struct {
	n int
}

func main() {
	r := regulator.New(100, 100, 10)

	go producer(r.In)
	go consumer(r.Out)

	go r.Start(10) // 10 ops/sec
	time.Sleep(1 * time.Second)
	r.Stop()
	fmt.Println("one stop")

	go r.Start(2) // 2 ops/secs
	fmt.Println("started again")
	time.Sleep(10 * time.Second)

	r.Stop()
	fmt.Println("bye")

}

func producer(out chan interface{}) {
	for i := 0; i < 100000000; i++ {
		out <- i
	}
}

func consumer(in chan interface{}) {
	for e := range in {
		fmt.Printf("Consumed %d\n", e.(int))
	}
}
