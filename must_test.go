package blobqueue_test

import (
	"fmt"
	"testing"

	"github.com/blueboardio/go-blobqueue"
	"github.com/blueboardio/go-blobqueue/memory"
)

func ExampleMust() {
	q := blobqueue.Must{new(memory.Queue)}
	fmt.Println(q.Len())
	q.Unshift([]byte("AA"))
	q.Push([]byte("BB"))
	fmt.Printf("%s\n", q.List())
	fmt.Println(q.Len())
	fmt.Printf("%s %s\n", q.Pop(), q.Shift())
	fmt.Println(q.Len())

	// Output:
	// 0
	// [AA BB]
	// 2
	// BB AA
	// 0
}

func expectPanic(t *testing.T, msg string) {
	if e := recover(); e != nil {
		errMsg := e.(error).Error()
		if errMsg != msg {
			t.Errorf("got panic %v, expected %q", e, msg)
		}
	} else {
		t.Errorf("expected panic %q, got nothing", msg)
	}
}

func TestMustPop(t *testing.T) {
	q := blobqueue.Must{new(memory.Queue)}
	defer expectPanic(t, blobqueue.ErrQueueIsEmtpy.Error())
	_ = q.Pop()
}

func TestMustShift(t *testing.T) {
	q := blobqueue.Must{new(memory.Queue)}
	defer expectPanic(t, blobqueue.ErrQueueIsEmtpy.Error())
	_ = q.Shift()
}
