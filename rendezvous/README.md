# Go Rendezvous Synchronisation Demo

## Overview

A **rendezvous** is a synchronisation point where multiple concurrent tasks must all arrive before they are allowed to continue.

In this program:

- We start `n` worker goroutines.
- Each worker:
  - Does some work for **Part A** (simulated by a random sleep).
  - Reports completion of Part A to the main goroutine.
  - Waits for a **start Part B** signal.
  - Then performs **Part B** and reports completion.
- The main goroutine:
  - Waits until it has seen **all** workers finish Part A and prints their IDs.
  - Then broadcasts a signal so that all workers may start Part B.
  - Finally waits until it has seen **all** workers finish Part B and prints their IDs.

Key property:  
**No “Part B” line is printed before all “Part A” lines have been printed.**

---

## Files

- `rendezvous.go`  
  Contains:
  - A commented-out **first attempt** `WorkWithRendezvous` and its `main`, which do not implement a true rendezvous.
  - The active implementation with:
    - `worker` goroutine that performs Part A and Part B.
    - `main` function that coordinates a proper rendezvous using channels.

---

## Implementation Details

### Commented-out template (naive version)

The block at the top:

```go
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
````

This version:

* Prints `Part A` and then `PartB` **inside each goroutine** with no real global synchronisation point.
* The comment `//Rendezvous here` marks where a rendezvous *should* happen, but there is no barrier logic.
* It is kept as a reference but does not run.

---

### Active implementation – proper rendezvous

The actual running code starts after the comment block:

```go
package main

import (
    "fmt"
    "math/rand/v2"
    "sync"
    "time"
)
```

#### `worker` goroutine

```go
func worker(
    id int,
    startB <-chan struct{},
    aCh chan<- int,
    bCh chan<- int,
    wg *sync.WaitGroup,
) {
    defer wg.Done()
    time.Sleep(time.Duration(rand.IntN(5)) * time.Second)
    aCh <- id  // signal: Part A completed by this worker
    <-startB   // block here until main broadcasts “start B”
    bCh <- id  // signal: Part B completed by this worker
}
```

For each worker:

1. Sleep for a random duration (simulating work for Part A).
2. Send its `id` on `aCh` to tell the main goroutine that Part A is done.
3. Block on `startB` waiting for the broadcast signal.
4. Once unblocked, send its `id` on `bCh` to indicate Part B is complete.

#### `main` function

```go
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

    // Phase 1: wait for all Part A completions
    seenA := make(map[int]bool)
    for printed := 0; printed < n; {
        id := <-aCh
        if !seenA[id] {
            fmt.Printf("Part A %d\n", id)
            seenA[id] = true
            printed++
        }
    }

    // Broadcast: allow all workers to start Part B
    close(startB)

    // Phase 2: wait for all Part B completions
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
```

Important details:

* `n := 5` – number of worker goroutines.
* `aCh` and `bCh` are buffered `chan int`:

  * `aCh` carries IDs of workers that finished Part A.
  * `bCh` carries IDs of workers that finished Part B.
* `startB` is an **unbuffered** `chan struct{}` used as a one-to-many broadcast.

**Phase 1: Part A**

* `seenA` tracks which IDs we have already printed.
* The loop:

  ```go
  for printed := 0; printed < n; {
      id := <-aCh
      if !seenA[id] {
          fmt.Printf("Part A %d\n", id)
          seenA[id] = true
          printed++
      }
  }
  ```

  waits until we have printed `n` distinct `Part A` lines (one for each worker).

**Rendezvous / Broadcast**

* After all Part A completions are seen, `close(startB)` is called.

  * Closing a channel unblocks **all** waiting receives.
  * As a result, every worker blocked on `<-startB` wakes up and continues to Part B.

**Phase 2: Part B**

* Symmetric code for B:

  ```go
  seenB := make(map[int]bool)
  for printed := 0; printed < n; {
      id := <-bCh
      if !seenB[id] {
          fmt.Printf("Part B %d\n", id)
          seenB[id] = true
          printed++
      }
  }
  ```

* Again we wait until we have printed `n` distinct `Part B` lines.

Finally, `wg.Wait()` ensures `main` doesn’t exit before all workers have finished.

---

## Requirements

* Go (version that supports `math/rand/v2`, typically Go 1.22+).
* Standard library only (no external dependencies).

---

## How to Run

1. Open a terminal and change into the directory containing `rendezvous.go`:

   ```bash
   cd rendezvous
   ```

2. Run the program:

   ```bash
   go run rendezvous.go
   ```

3. You should see output similar to:

   ```text
   Part A 3
   Part A 0
   Part A 4
   Part A 1
   Part A 2
   Part B 0
   Part B 2
   Part B 1
   Part B 3
   Part B 4
   ```

   The exact order of IDs within Part A or Part B will vary due to scheduling and random sleeps,
   but **all “Part A …” lines will appear before any “Part B …” lines**.

---

## Notes

* The commented-out `WorkWithRendezvous` shows a common mistake: printing A then B inside each goroutine without a true rendezvous.
* The active implementation moves the rendezvous logic into `main`, using channels to coordinate the phases cleanly.

