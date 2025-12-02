\# Go Reusable Barrier Synchronisation Demo



\## Overview



A \*\*barrier\*\* is a synchronisation primitive that makes a group of concurrent tasks (here: goroutines) all wait until \*every\* member of the group has reached the same point in the code. Only then can they all continue.



In this program:



\- We start `n` worker goroutines.

\- Each goroutine executes \*\*Part A\*\*, then waits at the \*\*barrier\*\*.

\- Once \*all\* goroutines have reached the barrier, they are released to execute \*\*Part B\*\*.

\- This process is done \*\*twice\*\*, to show that the barrier is reusable, not one-shot.



Key properties:



\- \*\*No goroutine prints “Part B” until all of them have printed “Part A”\*\*.

\- The same barrier instance is used for \*\*two rounds\*\* of goroutines.



This program is a simple demonstration of how to implement a \*\*reusable barrier\*\* using Go’s `sync/atomic` package and an unbuffered channel.



---



\## Files



\- `barrier2.go`  

&nbsp; Main Go source file containing:

&nbsp; - License header (GPL-3.0).

&nbsp; - The `ReusableBarrier` type and its methods.

&nbsp; - `doStuff` function that simulates work in two phases (A and B).

&nbsp; - `main` function that runs two rounds of goroutines using the same barrier.



\- `go.mod`  

&nbsp; Module definition for this small project (`module barrier2`).



---



\## Implementation Details



\### `ReusableBarrier` structure



The barrier is implemented as:



```go

type ReusableBarrier struct {

&nbsp;   max        int32

&nbsp;   arrived    int32

&nbsp;   waitChan   chan bool

&nbsp;   generation int32

}

````



\* `max`: total number of goroutines that must reach the barrier.

\* `arrived`: how many goroutines have reached the barrier in the current round.

\* `waitChan`: unbuffered channel used to block and release goroutines.

\* `generation`: a counter to track how many times the barrier has completed (a simple way to support reuse).



\### `NewReusableBarrier`



```go

func NewReusableBarrier(max int) \*ReusableBarrier {

&nbsp;   return \&ReusableBarrier{

&nbsp;       max:        int32(max),

&nbsp;       arrived:    0,

&nbsp;       waitChan:   make(chan bool),

&nbsp;       generation: 0,

&nbsp;   }

}

```



Creates and initialises a reusable barrier for `max` goroutines.



\### `(\*ReusableBarrier) Wait`



```go

func (b \*ReusableBarrier) Wait() {

&nbsp;   newArrived := atomic.AddInt32(\&b.arrived, 1)

&nbsp;   atomic.LoadInt32(\&b.generation)



&nbsp;   if newArrived == b.max {

&nbsp;       atomic.StoreInt32(\&b.arrived, 0)

&nbsp;       atomic.AddInt32(\&b.generation, 1)



&nbsp;       for i := int32(0); i < b.max-1; i++ {

&nbsp;           b.waitChan <- true

&nbsp;       }

&nbsp;   } else {

&nbsp;       <-b.waitChan

&nbsp;   }

}

```



Logic:



1\. Each goroutine calls `Wait`, which atomically increments `arrived`.

2\. If `newArrived == max`, this goroutine is the \*\*last\*\* to reach the barrier:



&nbsp;  \* It resets `arrived` back to `0` for the \*\*next round\*\*.

&nbsp;  \* It increments `generation` to mark a new barrier cycle.

&nbsp;  \* It sends `max - 1` signals into `waitChan`, releasing all waiting goroutines.

3\. If `newArrived < max`, the goroutine blocks on `<-b.waitChan` until the last goroutine arrives and sends a signal.



Because `arrived` and `generation` are manipulated using `sync/atomic`, we avoid using an explicit mutex.



---



\## `doStuff` function



```go

func doStuff(goNum int, barrier \*ReusableBarrier, wg \*sync.WaitGroup) {

&nbsp;   defer wg.Done()



&nbsp;   time.Sleep(time.Second)

&nbsp;   fmt.Println("Part A", goNum)



&nbsp;   barrier.Wait()



&nbsp;   fmt.Println("Part B", goNum)

}

```



\* Simulates some work in \*\*Part A\*\* using `time.Sleep`.

\* Prints `Part A <id>`.

\* Calls `barrier.Wait()` to synchronise with the other goroutines.

\* After the barrier is released, prints `Part B <id>`.



Because of the barrier, the output will always show \*\*all the “Part A” lines for that round first\*\*, and only after that will the “Part B” lines appear.

The exact order of IDs within Part A or Part B is not guaranteed.



---



\## `main` function



```go

func main() {

&nbsp;   totalRoutines := 10



&nbsp;   barrier := NewReusableBarrier(totalRoutines)

&nbsp;   var wg sync.WaitGroup

&nbsp;   wg.Add(totalRoutines)



&nbsp;   for i := 0; i < totalRoutines; i++ {

&nbsp;       go doStuff(i, barrier, \&wg)

&nbsp;   }



&nbsp;   wg.Wait()

&nbsp;   fmt.Println("All goroutines have completed their first round of execution.")



&nbsp;   fmt.Println("\\nCommencing the second round of implementation...")

&nbsp;   wg.Add(totalRoutines)

&nbsp;   for i := 0; i < totalRoutines; i++ {

&nbsp;       go doStuff(i+10, barrier, \&wg)

&nbsp;   }

&nbsp;   wg.Wait()

&nbsp;   fmt.Println("All goroutines have completed their second round of execution.")

}

```



Two rounds:



1\. \*\*First round\*\*



&nbsp;  \* Starts 10 goroutines with IDs `0` to `9`.

&nbsp;  \* Uses the barrier to ensure all `Part A` prints happen before any `Part B` prints.

&nbsp;  \* Prints a message when all have completed.



2\. \*\*Second round\*\*



&nbsp;  \* Reuses the \*\*same\*\* `ReusableBarrier`.

&nbsp;  \* Starts another 10 goroutines with IDs `10` to `19`.

&nbsp;  \* Again, the barrier ensures proper synchronisation.

&nbsp;  \* Prints a message when the second round is complete.



This shows that the barrier is \*\*reusable\*\*, not just a one-time synchronisation point.



---



\## Requirements



\* Go (any recent version should work).

&nbsp; The `go.mod` file currently specifies:



&nbsp; ```text

&nbsp; go 1.25.1

&nbsp; ```



&nbsp; If your installed Go version is different, you can update this line.



---



\## How to Run



1\. Open a terminal and change into the folder containing `go.mod` and `barrier2.go`, for example:



&nbsp;  ```bash

&nbsp;  cd barrier2

&nbsp;  ```



2\. Run the program:



&nbsp;  ```bash

&nbsp;  go run .

&nbsp;  ```



&nbsp;  or



&nbsp;  ```bash

&nbsp;  go run barrier2.go

&nbsp;  ```



3\. Observe the console output.

&nbsp;  You should see two rounds of output, something like:



&nbsp;  ```text

&nbsp;  Part A 3

&nbsp;  Part A 0

&nbsp;  ...

&nbsp;  Part A 9

&nbsp;  Part B 2

&nbsp;  ...

&nbsp;  Part B 9

&nbsp;  All goroutines have completed their first round of execution.



&nbsp;  Commencing the second round of implementation...

&nbsp;  Part A 13

&nbsp;  ...

&nbsp;  Part B 19

&nbsp;  All goroutines have completed their second round of execution.

&nbsp;  ```



Within each round:



\* All `Part A` messages appear before any `Part B` messages for that round.

\* The order of IDs is nondeterministic and depends on the scheduler.



