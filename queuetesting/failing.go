package queuetesting

import (
	"github.com/blueboardio/go-blobqueue"
)

// Failing wraps a queue to mock errors returned by a blobqueue.Queue.
type Failing struct {
	blobqueue.Queue
	NextError error
}

// List implements Queue interface.
func (q Failing) List() ([][]byte, error) {
	if err := q.NextError; err != nil {
		q.NextError = nil
		return nil, err
	}
	return q.Queue.List()
}

// Push implements Queue interface.
func (q Failing) Push(val []byte) error {
	if err := q.NextError; err != nil {
		q.NextError = nil
		return err
	}
	return q.Queue.Push(val)
}

// Unshift implements Queue interface.
func (q Failing) Unshift(val []byte) error {
	if err := q.NextError; err != nil {
		q.NextError = nil
		return err
	}
	return q.Queue.Unshift(val)
}

// Pop implements Queue interface.
func (q Failing) Pop() ([]byte, error) {
	if err := q.NextError; err != nil {
		q.NextError = nil
		return nil, err
	}
	return q.Queue.Pop()
}

// Shift implements Queue interface.
func (q Failing) Shift() ([]byte, error) {
	if err := q.NextError; err != nil {
		q.NextError = nil
		return nil, err
	}
	return q.Queue.Shift()
}

// Len implements Queue interface.
func (q Failing) Len() (int, error) {
	if err := q.NextError; err != nil {
		q.NextError = nil
		return 0, err
	}
	return q.Queue.Len()
}

// Empty implements Queue interface.
func (q Failing) Empty() error {
	if err := q.NextError; err != nil {
		q.NextError = nil
		return err
	}
	return q.Queue.Empty()
}
