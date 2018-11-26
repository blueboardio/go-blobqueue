package blobqueue_test

import (
	"fmt"

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
