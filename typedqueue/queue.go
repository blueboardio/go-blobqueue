// Package typedqueue wraps a blobqueue, with serialization and runtime type checking of values.
package typedqueue

import (
	"encoding"
	"fmt"
	"reflect"

	blobqueue "github.com/blueboardio/go-blobqueue"
)

type Value struct {
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler
}

type Queue struct {
	q    blobqueue.Queue
	t    reflect.Type
	zero encoding.BinaryMarshaler
}

// New wraps a queue to store any value of the same type that implements
// encoding.BinaryMarshaler and encoding.BinaryUnmarshaler.
//
// Any marshaling/unmarshaling error is unexpected and raises a panic.
// So any error returned by the queue methods come from the original queue.
func New(q blobqueue.Queue, zeroValue encoding.BinaryMarshaler) *Queue {
	t := reflect.TypeOf(zeroValue)
	// Checks of the contract for values
	_ = decodeValue(t, encodeValue(t, zeroValue))

	return &Queue{q: q, t: t, zero: zeroValue}
}

// Len returns the length of the queue.
func (q *Queue) Len() (int, error) {
	return q.q.Len()
}

// Empty clears the queue.
func (q *Queue) Empty() error {
	return q.q.Empty()
}

func encodeValue(t reflect.Type, v encoding.BinaryMarshaler) []byte {
	// Runtime check of type
	// This could be removed from production builds
	if reflect.TypeOf(v) != t {
		panic(fmt.Errorf("unexpected value of type %T", v))
	}

	b, err := v.MarshalBinary()
	if err != nil {
		panic(err)
	}
	return b
}

func decodeValue(t reflect.Type, b []byte) interface{} {
	v := reflect.New(t)
	if err := v.Interface().(encoding.BinaryUnmarshaler).UnmarshalBinary(b); err != nil {
		panic(err)
	}
	return v.Elem().Interface()
}

func decodeSlice(t reflect.Type, src [][]byte) interface{} {
	if len(src) == 0 {
		// nil slice
		return reflect.New(reflect.SliceOf(t)).Elem().Interface()
	}
	sl := reflect.MakeSlice(reflect.SliceOf(t), len(src), len(src))
	for i, b := range src {
		err := sl.Index(i).Addr().Interface().(encoding.BinaryUnmarshaler).UnmarshalBinary(b)
		if err != nil {
			panic(fmt.Errorf("index %d: %v", i, err))
		}
	}
	return sl.Interface()
}

// Push appends the queue with a new elem (val).
func (q *Queue) Push(val encoding.BinaryMarshaler) error {
	return q.q.Push(encodeValue(q.t, val))
}

// Unshift adds a new elem at the begining of the queue.
func (q *Queue) Unshift(val encoding.BinaryMarshaler) error {
	return q.q.Unshift(encodeValue(q.t, val))
}

// Pop deletes the last element of the queue and returns this precise element.
// If the queue is empty it returns error blobqueue.ErrQueueIsEmtpy.
func (q *Queue) Pop() (encoding.BinaryMarshaler, error) {
	b, err := q.q.Pop()
	if err != nil {
		return q.zero, err
	}
	return decodeValue(q.t, b).(encoding.BinaryMarshaler), nil
}

// Shift deletes the first element of the queue and returns this precise element.
// If the queue is empty it returns error blobqueue.ErrQueueIsEmtpy.
func (q *Queue) Shift() (encoding.BinaryMarshaler, error) {
	b, err := q.q.Shift()
	if err != nil {
		return q.zero, err
	}
	return decodeValue(q.t, b).(encoding.BinaryMarshaler), nil
}

// List returns all the items of the queue as a slice of values.
func (q *Queue) List() (interface{}, error) {
	src, err := q.q.List()
	if err != nil {
		return decodeSlice(q.t, nil), err
	}
	return decodeSlice(q.t, src), nil
}
