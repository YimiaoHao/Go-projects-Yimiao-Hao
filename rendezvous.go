/*package main

import (
	"fmt"
	"math/rand/v2"
	"sync"
	"time"
)

//Global variables shared between functions --A BAD IDEA

func WorkWithRendezvous(wg *sync.WaitGroup, Num int) bool {
	var X time.Duration
	X = time.Duration(rand.IntN(5))
	time.Sleep(X * time.Second) //wait random time amount
	fmt.Println("Part A", Num)
	//Rendezvous here

	fmt.Println("PartB", Num)
	wg.Done()
	return true
}

func main() {
	var wg sync.WaitGroup
	//barrier := make(chan bool)
	threadCount := 5

	wg.Add(threadCount)
	for N := range threadCount {
		go WorkWithRendezvous(&wg, N)
	}
	wg.Wait() //wait here until everyone (10 go routines) is done

} */

package main

import (
	"fmt"
	"math/rand/v2"
	"sync"
	"time"
)

func worker(id int, startB <-chan struct{}, aCh chan<- int, bCh chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	time.Sleep(time.Duration(rand.IntN(5)) * time.Second)
	aCh <- id // Send the signal indicating completion of A to the main thread
	<-startB  // Block here, waiting for the main thread to broadcast ‘Start B’.
	bCh <- id // Upon receiving the signal, transmit the ID completed by B.
}

func main() {
	n := 5
	aCh := make(chan int, n*2)
	bCh := make(chan int, n*2)
	startB := make(chan struct{})

	var wg sync.WaitGroup
	wg.Add(n)
	for i := 0; i < n; i++ {
		go worker(i, startB, aCh, bCh, &wg)
	}

	seenA := make(map[int]bool)
	for printed := 0; printed < n; {
		id := <-aCh
		if !seenA[id] {
			fmt.Printf("Part A %d\n", id)
			seenA[id] = true
			printed++
		}
	}

	close(startB)

	seenB := make(map[int]bool)
	for printed := 0; printed < n; {
		id := <-bCh
		if !seenB[id] {
			fmt.Printf("Part B %d\n", id)
			seenB[id] = true
			printed++
		}
	}

	wg.Wait()
}
