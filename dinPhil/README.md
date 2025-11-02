This project is based on the classroom example of “5 philosophers + 5 forks.”
The original (commented-out) template acquires forks in the order “pick my own fork first, then the one on the right,” which introduces a potential deadlock.
Building on that template, this implementation applies two classic deadlock-avoidance techniques:
1. a waiter channel that limits how many philosophers may enter the critical section at the same time (at most 4 may attempt to pick up forks);
2. a fixed acquisition order (“always pick the lower-numbered fork first, then the higher-numbered fork”), which breaks possible circular wait conditions.

In addition, this version adds explanatory comments and a Ctrl+C graceful shutdown to make it easier to demonstrate in class and to submit as coursework. Overall, it preserves the original problem setting but provides a runnable, explainable, and deadlock-free solution.
