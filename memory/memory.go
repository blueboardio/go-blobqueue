package memory

import (
	"github.com/blueboardio/go-blobqueue"
)

// Queue is the in memory based implementation of the blobqueue.Queue interface.
type Queue [][]byte

// List implements Queue interface
func (q Queue) List() ([][]byte, error) {
	return q, nil
}

// Push implements Queue interface
func (q *Queue) Push(val []byte) error {
	*q = append(*q, val)
	return nil
}

// Unshift implements Queue interface
func (q *Queue) Unshift(val []byte) error {
	*q = append([][]byte{val}, *q...)
	return nil
}

// Pop implements Queue interface
func (q *Queue) Pop() ([]byte, error) {
	if len(*q) == 0 {
		return nil, blobqueue.ErrQueueIsEmpty
	}
	save := (*q)[len(*q)-1]
	*q = (*q)[0 : len(*q)-1]
	return save, nil
}

// Shift implements Queue interface
func (q *Queue) Shift() ([]byte, error) {
	if len(*q) == 0 {
		return nil, blobqueue.ErrQueueIsEmpty
	}
	save := (*q)[0]
	*q = (*q)[1:]
	return save, nil
}

// Len implements Queue interface
func (q Queue) Len() (int, error) {
	return len(q), nil
}

// Empty implements Queue interface
func (q *Queue) Empty() error {
	*q = (*q)[:0]
	return nil
}
