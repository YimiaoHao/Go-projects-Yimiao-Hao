# Go Barrier Synchronisation Demo

## Overview

A **barrier** is a synchronisation primitive that makes a group of concurrent tasks (here: goroutines) all wait until *every* member of the group has reached the same point in the code. Only then can they all continue.

In this program:

- We start `n` worker goroutines.
- Each goroutine executes **Part A**, reaches a **barrier**, and waits.
- Once *all* goroutines have arrived at the barrier, they are released and can execute **Part B**.
- The key property is that **no goroutine prints “Part B” until all of them have printed “Part A”**.

This program is a simple demonstration of how to implement such a barrier using Go’s concurrency primitives.

---

## Files

- `barrier(2).go`  
  A single Go source file containing:
  - License header (GPL-3.0).
  - A `barrier` type that implements the barrier logic.
  - A `doStuff` function that simulates work in two phases (A and B).
  - A `main` function that starts multiple goroutines and waits for them to finish.

---

## Implementation Details

### `barrier` structure

The barrier is implemented as a small struct:

```go
type barrier struct {
    total int           // total number of goroutines that must reach the barrier
    count int           // number of goroutines that have arrived so far
    mu    sync.Mutex    // protects 'count'
    sem   chan struct{} // used to block / release waiting goroutines
}
````

* `total` is the number of goroutines that must reach the barrier.
* `count` is incremented as each goroutine arrives.
* `mu` protects `count` from data races.
* `sem` is used to block goroutines that arrive early and to release them when the last goroutine arrives.

### `newBarrier`

```go
func newBarrier(n int) *barrier {
    return &barrier{
        total: n,
        sem:   make(chan struct{}),
    }
}
```

Creates and initialises a barrier for `n` goroutines.

### `(*barrier) wait`

```go
func (b *barrier) wait() {
    b.mu.Lock()
    b.count++
    last := (b.count == b.total)
    b.mu.Unlock()

    if last {
        // Last goroutine to arrive: release all the others
        for i := 0; i < b.total-1; i++ {
            b.sem <- struct{}{}
        }
    } else {
        // Arrived early: wait to be released
        <-b.sem
    }
}
```

Logic:

1. Each goroutine increments `count` under the mutex.
2. The last one to arrive (`count == total`) will:

   * Send `total - 1` tokens into the channel, releasing all waiting goroutines.
3. All earlier goroutines block on `<-b.sem` until they receive a token.

---

## `doStuff` function

```go
func doStuff(id int, wg *sync.WaitGroup, br *barrier) {
    defer wg.Done()

    time.Sleep(time.Second)
    fmt.Println("Part A", id)

    br.wait() // barrier: wait for everyone to complete Part A

    fmt.Println("Part B", id)
}
```

* Simulates some work in **Part A** using `time.Sleep`.
* Prints `Part A <id>`.
* Calls `br.wait()` to synchronise with the other goroutines.
* After the barrier is released, prints `Part B <id>`.

Because of the barrier, the output will always show **all the “Part A” lines first**, and only after that will the “Part B” lines appear.
The *order of IDs* within Part A or Part B is not guaranteed and depends on scheduling.

---

## `main` function

```go
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
```

* Creates a barrier for `n = 10` goroutines.
* Starts 10 goroutines, each running `doStuff`.
* Uses a `sync.WaitGroup` to ensure `main` waits until all goroutines have finished.

---

## Requirements

* Go (tested with Go 1.21+, but any recent Go version should work).

---

## How to Run

1. Make sure Go is installed (`go version`).

2. In the directory containing the file, run:

   ```bash
   go run "barrier(2).go"
   ```

   > If you rename the file to `barrier.go`, you can also run:
   >
   > ```bash
   > go run barrier.go
   > ```

3. Observe the console output.
   You should see all the `Part A` lines first, followed by all the `Part B` lines, e.g.:

   ```text
   Part A 3
   Part A 0
   Part A 7
   Part A 2
   Part A 5
   Part A 1
   Part A 4
   Part A 8
   Part A 6
   Part A 9
   Part B 2
   Part B 7
   Part B 0
   Part B 3
   Part B 5
   Part B 1
   Part B 4
   Part B 8
   Part B 6
   Part B 9
   ```

   (The exact order of IDs may vary, but **no “Part B” line should appear before all “Part A” lines**.)

```
```
