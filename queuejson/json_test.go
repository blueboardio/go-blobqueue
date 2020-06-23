package queuejson_test

import (
	"reflect"
	"testing"

	"github.com/blueboardio/go-blobqueue"
	"github.com/blueboardio/go-blobqueue/memory"
	"github.com/blueboardio/go-blobqueue/queuejson"
)

func TestPushPop(t *testing.T) {
	var q = new(memory.Queue)

	checkEmpty := func() {
		var out interface{}
		err := queuejson.Pop(q, &out)
		if err != blobqueue.ErrQueueIsEmpty {
			t.Fatal("Pop should return ErrQueueIsEmpty")
		}
		if out != nil {
			t.Fatal("Failing Pop should not modify the target")
		}
	}

	checkEmpty()

	in := []interface{}{1.0, "abc", false}
	err := queuejson.Push(q, in)
	if err != nil {
		t.Fatal("Push should work")
	}

	if l, _ := q.Len(); l != 1 {
		t.Fatal("Len should be 1")
	}

	var out interface{}
	err = queuejson.Pop(q, &out)
	if err != nil {
		t.Fatal("Pop should work")
	}
	if !reflect.DeepEqual(in, out) {
		t.Fatalf("Pop: Bad value: %#v vs %#v", in, out)
	}

	checkEmpty()
}
