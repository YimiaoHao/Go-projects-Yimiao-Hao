# Go-projects-Yimiao-Hao

## Project Overview
Learn several typical synchronisation methods for goroutines in the ‘meet/block’ problem through four simplest feasible examples. Understand the synchronisation semantics in concurrency where ‘all A must complete before B collectively begins’, along with the trade-offs between one-time barriers and reusable barriers in implementation and performance. Grasp the minimal implementation of a barrier allowing N processes to synchronise once: each worker prints part A, waits at the barrier, and only prints part B after all participants have arrived.

## Example Output
rendezvous.go：
Part A 3
Part A 4
Part A 2
Part A 1
Part A 0
Part B 3
Part B 1
Part B 0
Part B 2
Part B 4

barrier2.go：
Part A 9
Part A 3
Part A 7
Part A 6
Part A 1
Part A 4
Part A 0
Part A 5
Part A 2
Part A 8
Part B 8
Part B 6
Part B 9
Part B 3
Part B 7
Part B 4
Part B 1
Part B 0
Part B 5
Part B 2
All goroutines have completed their first round of execution.

Commencing the second round of implementation...
Part A 10
Part A 15
Part A 19
Part A 18
Part A 17
Part A 14
Part A 12
Part A 11
Part B 11
Part B 18
Part B 10
Part B 15
Part B 19
Part B 16
Part B 17
Part B 13
Part B 14
Part B 12
All goroutines have completed their second round of execution.

barrierstruct.go：
Part A 0
Part A 1
Part A 4
Part A 2
Part A 3
Part B 3
Part B 1
Part B 0
Part B 2
Part B 4

barrier(2).go：
Part A 2
Part A 4
Part A 0
Part A 5
Part A 3
Part A 1
Part A 9
Part A 8
Part A 7
Part A 6
Part B 6
Part B 5
Part B 2
Part B 4
Part B 0
Part B 1
Part B 3
Part B 9
Part B 8
Part B 7
