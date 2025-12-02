# Go Reusable Barrier Synchronisation Demo

## Overview

A **barrier** is a synchronisation primitive that makes a group of concurrent tasks (here: goroutines) all wait until *every* member of the group has reached the same point in the code. Only then can they all continue.

In this program:

- We start `n` worker goroutines.

- Each goroutine executes **Part A**, then waits at the **barrier**.

- Once *all* goroutines have reached the barrier, they are released to execute **Part B**.

- This process is done **twice**, to show that the barrier is reusable, not one-shot.

Key properties:

- **No goroutine prints “Part B” until all of them have printed “Part A”**.

- The same barrier instance is used for **two rounds** of goroutines.

This program is a simple demonstration of how to implement a **reusable barrier** using Go’s `sync/atomic` package and an unbuffered channel.

---

## Files

- `barrier2.go`  

  Main Go source file containing:

  - License header (GPL-3.0).

  - The `ReusableBarrier` type and its methods.

  - `doStuff` function that simulates work in two phases (A and B).

  - `main` function that runs two rounds of goroutines using the same barrier.



- `go.mod`  

  Module definition for this small project (`module barrier2`).
  
---

## Implementation Details

### `ReusableBarrier` structure

The barrier is implemented as:

```go

type ReusableBarrier struct {

  max        int32

  arrived    int32

  waitChan   chan bool

  generation int32

}

```

* `max`: total number of goroutines that must reach the barrier.

* `arrived`: how many goroutines have reached the barrier in the current round.

* `waitChan`: unbuffered channel used to block and release goroutines.

* `generation`: a counter to track how many times the barrier has completed (a simple way to support reuse).


### `NewReusableBarrier`

```go

func NewReusableBarrier(max int) *ReusableBarrier {

  return &ReusableBarrier{

    max:        int32(max),

    arrived:    0,

    waitChan:   make(chan bool),

    generation: 0,

  }

}

```

Creates and initialises a reusable barrier for `max` goroutines.

### `(*ReusableBarrier) Wait`

```go

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

```

Logic:

1. Each goroutine calls `Wait`, which atomically increments `arrived`.

2. If `newArrived == max`, this goroutine is the **last** to reach the barrier:



   * It resets `arrived` back to `0` for the **next round**.

   * It increments `generation` to mark a new barrier cycle.

   * It sends `max - 1` signals into `waitChan`, releasing all waiting goroutines.

3. If `newArrived < max`, the goroutine blocks on `<-b.waitChan` until the last goroutine arrives and sends a signal.

Because `arrived` and `generation` are manipulated using `sync/atomic`, we avoid using an explicit mutex.

---

## `doStuff` function

```go

func doStuff(goNum int, barrier *ReusableBarrier, wg *sync.WaitGroup) {

  defer wg.Done()

  time.Sleep(time.Second)

  fmt.Println("Part A", goNum)

  barrier.Wait()

  fmt.Println("Part B", goNum)

}

```

* Simulates some work in **Part A** using `time.Sleep`.

* Prints `Part A <id>`.

* Calls `barrier.Wait()` to synchronise with the other goroutines.

* After the barrier is released, prints `Part B <id>`.

Because of the barrier, the output will always show **all the “Part A” lines for that round first**, and only after that will the “Part B” lines appear.

The exact order of IDs within Part A or Part B is not guaranteed.

---

## `main` function

```go

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

```

Two rounds:

1. **First round**

   * Starts 10 goroutines with IDs `0` to `9`.

   * Uses the barrier to ensure all `Part A` prints happen before any `Part B` prints.

   * Prints a message when all have completed.

2. **Second round**

   * Reuses the **same** `ReusableBarrier`.

   * Starts another 10 goroutines with IDs `10` to `19`.

   * Again, the barrier ensures proper synchronisation.

   * Prints a message when the second round is complete.

This shows that the barrier is **reusable**, not just a one-time synchronisation point.

---

## Requirements

* Go (any recent version should work).

  The `go.mod` file currently specifies:

  ```text

  go 1.25.1

  ```

  If your installed Go version is different, you can update this line.

---

## How to Run

1. Open a terminal and change into the folder containing `go.mod` and `barrier2.go`, for example:

   ```bash

   cd barrier2

   ```

2. Run the program:

   ```bash

   go run .

   ```

   or

   ```bash

   go run barrier2.go

   ```

3. Observe the console output.

   You should see two rounds of output, something like:

   ```text

   Part A 3

   Part A 0

   ...

   Part A 9

   Part B 2

   ...

   Part B 9

   All goroutines have completed their first round of execution.


   Commencing the second round of implementation...

   Part A 13

   ...

   Part B 19

   All goroutines have completed their second round of execution.

   ```


Within each round:


* All `Part A` messages appear before any `Part B` messages for that round.

* The order of IDs is nondeterministic and depends on the scheduler.
