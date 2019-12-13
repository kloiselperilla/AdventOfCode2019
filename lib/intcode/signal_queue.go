package intcode

import (
	//"fmt"
	"sync"
)

// SignalQueue gives an implementation of a connection between engines
type SignalQueue struct {
	Queue []int
	Cond  *sync.Cond
}

func NewSignalQueue() SignalQueue {
	mux := sync.Mutex{}
	q := SignalQueue{Cond: sync.NewCond(&mux)}
	return q
}

// Enqueue adds to end of queue
func (q *SignalQueue) Enqueue(val int) {
	q.Cond.L.Lock()
	q.Queue = append(q.Queue, val)
	q.Cond.Broadcast()
	q.Cond.L.Unlock()
}

// Dequeue removes from beginning of queue
func (q *SignalQueue) Dequeue() int {
	q.Cond.L.Lock()
	// Wait for not empty
	for len(q.Queue) == 0 {
		q.Cond.Wait()
	}
	retval := q.Queue[0]
	q.Queue[0] = 0
	q.Queue = q.Queue[1:]
	q.Cond.L.Unlock()

	return retval
}

// IsEmpty checks if empty
func (q *SignalQueue) IsEmpty() bool {
	return len(q.Queue) == 0
}
