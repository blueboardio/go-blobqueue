package queuejson

import (
	"encoding/json"
	"github.com/blueboardio/go-blobqueue"
)

// Push serialize v as JSON and pushes the blob on the queue.
// See json.Marshal for the accepted values.
func Push(q blobqueue.Queue, v interface{}) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return q.Push(b)
}

// Unshift serialize v as JSON and pushes the blob on the queue.
// See json.Marshal for the accepted values.
func Unshift(q blobqueue.Queue, v interface{}) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return q.Unshift(b)
}

// Pop removes the tail of the queue and deserialize it as JSON into v.
// v must be a pointer (see json.Unmarshal).
func Pop(q blobqueue.Queue, v interface{}) error {
	b, err := q.Pop()
	if err != nil {
		return err
	}
	return json.Unmarshal(b, v)
}

// Shift removes the head of the queue and deserialize it as JSON into v.
// v must be a pointer (see json.Unmarshal).
func Shift(q blobqueue.Queue, v interface{}) error {
	b, err := q.Shift()
	if err != nil {
		return err
	}
	return json.Unmarshal(b, v)
}
