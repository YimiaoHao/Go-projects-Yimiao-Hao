# Dining Philosophers – Deadlock-Free Go Implementation

## Overview

The Dining Philosophers problem models **five philosophers** sitting around a table, alternating between **thinking** and **eating**. Between each pair of philosophers there is one fork. To eat, a philosopher needs **two forks** (left and right).

The challenges:

- Avoid **deadlock** (everyone holds one fork and waits forever for the other).
- Avoid **starvation** as much as possible.
- Correctly coordinate concurrent access to forks.

In this program:

- Each philosopher is a **goroutine** running in an infinite loop: think → pick forks → eat → put forks down.
- Forks are modelled as **buffered channels of size 1** (token = fork available).
- A **“room”** channel acts like a waiter / semaphore, limiting how many philosophers can try to pick up forks at the same time.
- Philosophers always pick up the **lower-numbered fork first**, then the higher-numbered one (**resource ordering**).

These choices together make the solution **deadlock-free**.

---

## Files

- `dinPhil.go`  
  Contains:
  - The original **template implementation** (commented out) that can deadlock.
  - The new **deadlock-free implementation** using channels and a waiter.
  - `think` and `eat` functions to simulate behaviour.
  - `philosopher` goroutine function.
  - `main` function which starts the philosophers and waits for a signal to exit.

---

## Implementation Details

### Commented-out template (can deadlock)

At the top of the file, inside a block comment `/* ... */`, there is an original version using:

- `forks map[int]chan bool`
- `getForks`, `putForks`, `doPhilStuff` functions
- A `main` that starts 5 philosophers and waits with a `sync.WaitGroup`

In that version, each philosopher simply:

1. Thinks
2. Picks up left fork then right fork
3. Eats
4. Puts down forks

Because all philosophers can try to pick up their left fork at the same time, this approach **can deadlock**. That code is kept for reference only and does **not** run.

---

### Active implementation (deadlock-free)

Below the commented section, the actual running program starts with:

```go
package main

import (
    "fmt"
    "math/rand/v2"
    "os"
    "os/signal"
    "syscall"
    "time"
)
````

#### Thinking and eating

```go
func think(id int) {
    time.Sleep(time.Duration(rand.IntN(5)) * time.Second)
    fmt.Println("Phil", id, "was thinking")
}

func eat(id int) {
    time.Sleep(time.Duration(rand.IntN(5)) * time.Second)
    fmt.Println("Phil", id, "was eating")
}
```

* Each philosopher sleeps for a random duration (0–4 seconds) before printing their action.
* This simulates non-deterministic thinking and eating times.

#### Philosopher goroutine

```go
// Deadlock-free philosopher (infinite loop)
// - forks: each fork is a buffered channel of size 1 (token = fork available)
// - room : waiter/semaphore, at most N-1 philosophers may try to pick forks at the same time
func philosopher(id int, forks []chan struct{}, room chan struct{}) {
    left := id
    right := (id + 1) % len(forks)

    // Resource ordering: always pick the lower-numbered fork first,
    // then the higher-numbered fork → breaks circular wait
    low, high := left, right
    if high < low {
        low, high = high, low
    }

    for {
        think(id)

        // enter the dining room (N-1 limit)
        room <- struct{}{}

        // pick forks in fixed order
        <-forks[low]
        <-forks[high]

        eat(id)

        // put forks back and leave the room
        forks[high] <- struct{}{}
        forks[low] <- struct{}{}
        <-room
    }
}
```

Key ideas:

* **Forks**: each fork is a `chan struct{}` with capacity 1.

  * Having a token in the channel = fork is available.
  * Receiving from the channel = picking up the fork.
  * Sending back into the channel = putting the fork down.
* **Waiter (`room` channel)**:

  * `room` is a buffered channel with capacity `n-1`.
  * Before trying to pick up forks, a philosopher sends a token into `room`.
  * This ensures at most `n-1` philosophers are in the “picking forks” section at once, which helps avoid deadlock.
* **Resource ordering**:

  * Each philosopher calculates `low` and `high` as the lower/higher-numbered fork indices.
  * They *always* pick up `forks[low]` first, then `forks[high]`.
  * This breaks the circular wait condition.

#### `main` function

```go
func main() {
    n := 5

    // init forks: each fork starts with 1 token = available
    forks := make([]chan struct{}, n)
    for i := 0; i < n; i++ {
        forks[i] = make(chan struct{}, 1)
        forks[i] <- struct{}{}
    }

    // waiter: allow at most n-1 philosophers to try to pick forks
    room := make(chan struct{}, n-1)

    // start n philosophers (infinite loop)
    for i := 0; i < n; i++ {
        go philosopher(i, forks, room)
    }

    // wait for Ctrl+C / termination signal to exit manually
    sigCh := make(chan os.Signal, 1)
    signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
    <-sigCh
    fmt.Println("\nShutting down...")
}
```

* Creates `n = 5` forks, each with one initial token.
* Creates the `room` semaphore with capacity `n-1`.
* Starts `n` philosopher goroutines.
* The program runs **until you press Ctrl+C**, at which point it prints a shutdown message and exits.

---

## Requirements

* Go (version supporting `math/rand/v2`, e.g. Go 1.22+).
* Standard library only (no external dependencies).

---

## How to Run

1. Open a terminal and change into the `dinPhil` directory:

   ```bash
   cd dinPhil
   ```

2. Run the program:

   ```bash
   go run dinPhil.go
   ```

3. You should see continuous output similar to:

   ```text
   Phil 0 was thinking
   Phil 2 was thinking
   Phil 3 was eating
   Phil 1 was eating
   Phil 4 was thinking
   Phil 0 was eating
   ...
   ```

   Philosophers will keep thinking and eating until you stop the program with **Ctrl+C**.

---


