package intcode

import (
	"sync"
)

// SignalQueue gives an implementation of a connection between engines
type SignalQueue struct {
	Queue []int
	ready chan bool
	cond  *sync.Cond
}

func NewSignalQueue() SignalQueue {
	mux := sync.Mutex{}
	q := SignalQueue{ready: make(chan bool), cond: sync.NewCond(&mux)}
	return q
}

// Enqueue adds to end of queue
func (q *SignalQueue) Enqueue(val int) {
	q.cond.L.Lock()
	q.Queue = append(q.Queue, val)
	q.cond.Broadcast()
	q.cond.L.Unlock()
}

// Dequeue removes from beginning of queue
func (q *SignalQueue) Dequeue() int {
	q.cond.L.Lock()
	// Wait for not empty
	for len(q.Queue) == 0 {
		q.cond.Wait()
	}
	retval := q.Queue[0]
	q.Queue[0] = 0
	q.Queue = q.Queue[1:]
	q.cond.L.Unlock()

	return retval
}
