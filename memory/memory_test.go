package memory_test

import (
	"testing"

	"github.com/blueboardio/go-blobqueue"

	"github.com/blueboardio/go-blobqueue/memory"
	"github.com/blueboardio/go-blobqueue/queuetesting"
)

func TestInMemoryImplem(t *testing.T) {
	queue := blobqueue.SafeQueue(&memory.Queue{})
	queuetesting.RunTests(t, queue, false)
}

func BenchmarkMemoryPush(b *testing.B) {
	queue := blobqueue.SafeQueue(&memory.Queue{})
	queuetesting.RunBenchmarkPush(b, queue)
}
