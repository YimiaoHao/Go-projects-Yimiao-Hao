package main

import (
	"fmt"
	"math/rand/v2"
	"sync"
	"time"
)

type barrier struct {
	theChan chan bool
	theLock sync.Mutex
	total   int
	count   int
}

func createBarrier(n int) *barrier {
	return &barrier{
		theChan: make(chan bool),
		total:   n,
	}
}

func (b *barrier) wait() {
	b.theLock.Lock()
	b.count++
	last := (b.count == b.total)
	b.theLock.Unlock()

	if last {
		for i := 0; i < b.total-1; i++ {
			<-b.theChan
		}
	} else {
		b.theChan <- true
	}
}

func WorkWithRendezvous(wg *sync.WaitGroup, id int, b *barrier) {
	defer wg.Done()

	time.Sleep(time.Duration(rand.IntN(5)) * time.Second)

	fmt.Println("Part A", id)

	b.wait()

	fmt.Println("Part B", id)
}

func main() {
	threadCount := 5
	var wg sync.WaitGroup
	wg.Add(threadCount)

	b := createBarrier(threadCount)

	for i := 0; i < threadCount; i++ {
		go WorkWithRendezvous(&wg, i, b)
	}

	wg.Wait()
}
