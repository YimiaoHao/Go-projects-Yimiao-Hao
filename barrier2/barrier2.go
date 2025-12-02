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
// Description:
// A simple barrier implemented using mutex and unbuffered channel
// Issues:
// None I hope
//1. Change mutex to atomic variable
//2. Make it a reusable barrier
//--------------------------------------------

package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type ReusableBarrier struct {
	max        int32
	arrived    int32
	waitChan   chan bool
	generation int32
}

func NewReusableBarrier(max int) *ReusableBarrier {
	return &ReusableBarrier{
		max:        int32(max),
		arrived:    0,
		waitChan:   make(chan bool),
		generation: 0,
	}
}

func (b *ReusableBarrier) Wait() {
	newArrived := atomic.AddInt32(&b.arrived, 1)
	atomic.LoadInt32(&b.generation)

	if newArrived == b.max {
		atomic.StoreInt32(&b.arrived, 0)
		atomic.AddInt32(&b.generation, 1)

		for i := int32(0); i < b.max-1; i++ {
			b.waitChan <- true
		}
	} else {
		<-b.waitChan
	}
}

func doStuff(goNum int, barrier *ReusableBarrier, wg *sync.WaitGroup) {
	defer wg.Done()

	time.Sleep(time.Second)
	fmt.Println("Part A", goNum)

	barrier.Wait()

	fmt.Println("Part B", goNum)
}

func main() {
	totalRoutines := 10

	barrier := NewReusableBarrier(totalRoutines)
	var wg sync.WaitGroup
	wg.Add(totalRoutines)

	for i := 0; i < totalRoutines; i++ {
		go doStuff(i, barrier, &wg)
	}

	wg.Wait()
	fmt.Println("All goroutines have completed their first round of execution.")

	fmt.Println("\nCommencing the second round of implementation...")
	wg.Add(totalRoutines)
	for i := 0; i < totalRoutines; i++ {
		go doStuff(i+10, barrier, &wg)
	}
	wg.Wait()
	fmt.Println("All goroutines have completed their second round of execution.")
}
