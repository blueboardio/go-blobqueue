package queuebinary_test

import (
	"testing"

	"github.com/blueboardio/go-blobqueue"
	"github.com/blueboardio/go-blobqueue/memory"
	"github.com/blueboardio/go-blobqueue/queuebinary"
)

type S string

func (s S) MarshalBinary() ([]byte, error) {
	return []byte(s), nil
}

func (s *S) UnmarshalBinary(b []byte) error {
	*s = S(b)
	return nil
}

func TestPushPop(t *testing.T) {
	var q = new(memory.Queue)

	checkEmpty := func() {
		var out S
		err := queuebinary.Pop(q, &out)
		if err != blobqueue.ErrQueueIsEmpty {
			t.Fatal("Pop should return ErrQueueIsEmpty")
		}
		if out != "" {
			t.Fatal("Failing Pop should not modify the target")
		}
	}

	checkEmpty()

	err := queuebinary.Push(q, S("abc"))
	if err != nil {
		t.Fatal("Push should work")
	}

	if l, _ := q.Len(); l != 1 {
		t.Fatal("Len should be 1")
	}

	var out S
	err = queuebinary.Pop(q, &out)
	if err != nil {
		t.Fatal("Pop should work")
	}
	if out != "abc" {
		t.Fatal("Pop: Bad value")
	}

	checkEmpty()
}
