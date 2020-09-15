// Package queuemsgpack provides utilities to inject/remove elements
// from a blobqueue using MessagePack encoding.
package queuemsgpack

import (
	"github.com/blueboardio/go-blobqueue"

	// 2020-06-24 v4 is the stable branch
	"github.com/vmihailenco/msgpack/v4"
)

// Push serialize v as MessagePack and pushes the blob on the queue.
// See msgpack.Marshal for the accepted values.
func Push(q blobqueue.Queue, v interface{}) error {
	b, err := msgpack.Marshal(v)
	if err != nil {
		return err
	}
	return q.Push(b)
}

// Unshift serialize v as MessagePack and pushes the blob on the queue.
// See msgpack.Marshal for the accepted values.
func Unshift(q blobqueue.Queue, v interface{}) error {
	b, err := msgpack.Marshal(v)
	if err != nil {
		return err
	}
	return q.Unshift(b)
}

// Pop removes the tail of the queue and deserialize it as MessagePack into v.
// v must be a pointer (see msgpack.Unmarshal).
func Pop(q blobqueue.Queue, v interface{}) error {
	b, err := q.Pop()
	if err != nil {
		return err
	}
	return msgpack.Unmarshal(b, v)
}

// Shift removes the head of the queue and deserialize it as MessagePack into v.
// v must be a pointer (see msgpack.Unmarshal).
func Shift(q blobqueue.Queue, v interface{}) error {
	b, err := q.Shift()
	if err != nil {
		return err
	}
	return msgpack.Unmarshal(b, v)
}

// Peek returns the head of the queue without removing it and deserialize it as MessagePack into v.
// v must be a pointer (see msgpack.Unmarshal).
func Peek(q blobqueue.Queue, v interface{}) error {
	b, err := q.Peek()
	if err != nil {
		return err
	}
	return msgpack.Unmarshal(b, v)
}
