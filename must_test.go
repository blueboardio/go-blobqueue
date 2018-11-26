package blobqueue_test

import (
	"fmt"

	"github.com/blueboardio/go-blobqueue"
	"github.com/blueboardio/go-blobqueue/memory"
)

func ExampleMust() {
	q := blobqueue.Must{new(memory.Queue)}
	fmt.Println(q.Len())
	q.Push([]byte("AA"))
	fmt.Println(q.Len())
	fmt.Printf("%s\n", q.Pop())
	fmt.Println(q.Len())

	// Output:
	// 0
	// 1
	// AA
	// 0
}
