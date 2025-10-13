//Barrier.go Template Code
//Copyright (C) 2024 Dr. Joseph Kehoe

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

//--------------------------------------------
// Author: Joseph Kehoe (Joseph.Kehoe@setu.ie)
// Created on 30/9/2024
// Modified by:
// Issues:
// The barrier is not implemented!
//--------------------------------------------

/*
package main

import (
	"context"
	"fmt"
	"sync"
	"time"
	"golang.org/x/sync/semaphore"
)


// Place a barrier in this function --use Mutex's and Semaphores
func doStuff(goNum int, wg *sync.WaitGroup) bool {
	time.Sleep(time.Second)
	fmt.Println("Part A",goNum)
	//we wait here until everyone has completed part A
	fmt.Println("PartB",goNum)
	wg.Done()
	return true
}


func main() {
	totalRoutines:=10
	var wg sync.WaitGroup
	wg.Add(totalRoutines)
	//we will need some of these
	ctx := context.TODO()
	var theLock sync.Mutex
	sem := semaphore.NewWeighted(int64(totalRoutines))
	theLock.Lock()
	sem.Acquire(ctx, 1)
	for i := range totalRoutines {//create the go Routines here
		go doStuff(i, &wg)
	}
	sem.Release(1)
	theLock.Unlock()

	wg.Wait() //wait for everyone to finish before exiting
}
*/

package main

import (
	"fmt"
	"sync"
	"time"
)

type barrier struct {
	total int
	count int
	mu    sync.Mutex
	sem   chan struct{}
}

func newBarrier(n int) *barrier {
	return &barrier{
		total: n,
		sem:   make(chan struct{}),
	}
}

func (b *barrier) wait() {
	b.mu.Lock()
	b.count++
	last := (b.count == b.total)
	b.mu.Unlock()

	if last {
		for i := 0; i < b.total-1; i++ {
			b.sem <- struct{}{}
		}
	} else {
		<-b.sem
	}
}

func doStuff(id int, wg *sync.WaitGroup, br *barrier) {
	defer wg.Done()

	time.Sleep(time.Second)
	fmt.Println("Part A", id)

	br.wait()

	fmt.Println("Part B", id)
}

func main() {
	n := 10
	var wg sync.WaitGroup
	wg.Add(n)

	br := newBarrier(n)

	for i := 0; i < n; i++ {
		go doStuff(i, &wg, br)
	}

	wg.Wait()
}
