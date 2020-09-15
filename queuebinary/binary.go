package queuebinary

import (
	"encoding"
	"github.com/blueboardio/go-blobqueue"
)

func Push(q blobqueue.Queue, v encoding.BinaryMarshaler) error {
	b, err := v.MarshalBinary()
	if err != nil {
		return err
	}
	return q.Push(b)
}

func Unshift(q blobqueue.Queue, v encoding.BinaryMarshaler) error {
	b, err := v.MarshalBinary()
	if err != nil {
		return err
	}
	return q.Unshift(b)
}

func Pop(q blobqueue.Queue, v encoding.BinaryUnmarshaler) error {
	b, err := q.Pop()
	if err != nil {
		return err
	}
	return v.UnmarshalBinary(b)
}

func Shift(q blobqueue.Queue, v encoding.BinaryUnmarshaler) error {
	b, err := q.Shift()
	if err != nil {
		return err
	}
	return v.UnmarshalBinary(b)
}

func Peek(q blobqueue.Queue, v encoding.BinaryUnmarshaler) error {
	b, err := q.Peek()
	if err != nil {
		return err
	}
	return v.UnmarshalBinary(b)
}
