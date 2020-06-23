package blobqueue

import (
	"errors"
	"sync"
)

/*
Queue is an interface providing operations to handle queues.
All methods may return an error as its implementation may have runtime failures.
*/
type Queue interface {
	// List returns all the items of the queue as [][]byte.
	List() ([][]byte, error)

	// Push appends the queue with a new elem (val).
	Push(val []byte) error
	// Unshift add a new elem at the beggining of the queue.
	Unshift(val []byte) error

	// Pop deletes the last element of the queue, it returns this precise element.
	// If the queue is empty it returns error ErrQueueIsEmpty.
	Pop() ([]byte, error)

	// Shift deletes the first element of the queue, it returns this precise element.
	// If the queue is empty it returns error ErrQueueIsEmpty.
	Shift() ([]byte, error)

	// Len returns the length of the queue.
	Len() (int, error)
	// Empty clears the queue.
	Empty() error
}

// These variables defines typed errors.
var (
	ErrIndexOutOfRange = errors.New("index out of range")
	ErrQueueIsEmpty    = errors.New("empty queue")
)

type safeQueue struct {
	mu sync.RWMutex
	q  Queue
}

// SafeQueue takes a Queue and wraps it to lock properly all operations.
func SafeQueue(q Queue) Queue {
	return &safeQueue{
		q: q,
	}
}

// List implements `Queue` interface.
func (q *safeQueue) List() ([][]byte, error) {
	q.mu.RLock()
	defer q.mu.RUnlock()
	list, err := q.q.List()
	return append([][]byte(nil), list...), err
}

// Push implements `Queue` interface.
func (q *safeQueue) Push(val []byte) error {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.q.Push(val)
}

// Unshift implements Queue interface
func (q *safeQueue) Unshift(val []byte) error {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.q.Unshift(val)
}

// Pop implements Queue interface
func (q *safeQueue) Pop() ([]byte, error) {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.q.Pop()
}

// Shift implements Queue interface
func (q *safeQueue) Shift() ([]byte, error) {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.q.Shift()
}

// Len implements Queue interface
func (q *safeQueue) Len() (int, error) {
	q.mu.RLock()
	defer q.mu.RUnlock()
	return q.q.Len()
}

// Empty implements Queue interface
func (q *safeQueue) Empty() error {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.q.Empty()
}
