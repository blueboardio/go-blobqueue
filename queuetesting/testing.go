package queuetesting

import (
	"fmt"
	"testing"

	"github.com/blueboardio/go-blobqueue"
	"github.com/stretchr/testify/assert"
)

const (
	listAction = iota
	pushAction
	unshiftAction
	lenAction
	popAction
	shiftAction
	emptyAction
)

type action struct {
	Type       int
	Index      int
	Val        []byte
	ExpectList [][]byte
	ExpectVal  []byte
	ExpectLen  int
	ExpectErr  error
	ShouldErr  bool
}

type testSuite struct {
	Queue   blobqueue.Queue
	Actions []action
}

var actions = []action{
	{Type: listAction, ExpectList: nil, ShouldErr: false},
	{Type: pushAction, Val: []byte("elem 1"), ShouldErr: false},
	{Type: lenAction, ExpectLen: 1, ShouldErr: false},
	{Type: unshiftAction, Val: []byte("elem 0"), ShouldErr: false},
	{Type: lenAction, ExpectLen: 2, ShouldErr: false},
	{Type: pushAction, Val: []byte("elem 2"), ShouldErr: false},
	{Type: lenAction, ExpectLen: 3, ShouldErr: false},
	{Type: pushAction, Val: []byte("elem 3"), ShouldErr: false},
	{Type: lenAction, ExpectLen: 4, ShouldErr: false},
	{Type: listAction, ExpectList: [][]byte{[]byte("elem 0"), []byte("elem 1"), []byte("elem 2"), []byte("elem 3")}},
	{Type: popAction, ExpectVal: []byte("elem 3"), ShouldErr: false},
	{Type: lenAction, ExpectLen: 3, ShouldErr: false},
	{Type: shiftAction, ExpectVal: []byte("elem 0"), ShouldErr: false},
	{Type: lenAction, ExpectLen: 2, ShouldErr: false},
	{Type: listAction, ExpectList: [][]byte{[]byte("elem 1"), []byte("elem 2")}},
	{Type: emptyAction, ShouldErr: false},
	{Type: emptyAction, ShouldErr: false},
	{Type: lenAction, ExpectLen: 0, ShouldErr: false},
	{Type: listAction, ExpectList: nil, ShouldErr: false},
	{Type: popAction, ShouldErr: true, ExpectErr: blobqueue.ErrQueueIsEmtpy, ExpectVal: nil},
	{Type: shiftAction, ShouldErr: true, ExpectErr: blobqueue.ErrQueueIsEmtpy, ExpectVal: nil},
	{Type: lenAction, ExpectLen: 0, ShouldErr: false},
	{Type: listAction, ExpectList: nil, ShouldErr: false},
}

/*
RunTests runs a list of predefined actions on a implementation of blobqueue.Queue.
If you're dependant of an external data storage (redis, db, etc..) and you want to test
client specific error (e.g. a wrong redis connection) you can set shouldAlwaysFail to true.
*/
func RunTests(t *testing.T, queue blobqueue.Queue, shouldAlwaysFail bool) {
	suite := testSuite{
		Queue:   queue,
		Actions: actions,
	}
	for i := 0; i < 2; i++ {
		for i, action := range suite.Actions {
			infoMsg := fmt.Sprintf("assertion error, from action at index %d", i)
			var err error
			switch action.Type {
			case listAction:
				var list [][]byte
				list, err = suite.Queue.List()
				if shouldAlwaysFail {
					action.ExpectList = nil
				}
				assert.Equal(t, action.ExpectList, list, infoMsg)
				break
			case pushAction:
				err = suite.Queue.Push(action.Val)
				break
			case unshiftAction:
				err = suite.Queue.Unshift(action.Val)
				break
			case lenAction:
				var len int
				len, err = suite.Queue.Len()
				if shouldAlwaysFail {
					action.ExpectLen = 0
				}
				assert.Equal(t, action.ExpectLen, len)
				break
			case popAction:
				var val []byte
				val, err = suite.Queue.Pop()
				if shouldAlwaysFail {
					action.ExpectVal = nil
				}
				assert.Equal(t, action.ExpectVal, val, infoMsg)
				break
			case shiftAction:
				var val []byte
				val, err = suite.Queue.Shift()
				if shouldAlwaysFail {
					action.ExpectVal = nil
				}
				assert.Equal(t, action.ExpectVal, val, infoMsg)
				break
			case emptyAction:
				err = suite.Queue.Empty()
				break
			default:
				t.Fatalf("Your test suite has unexpected action type at index %d", i)
			}
			if shouldAlwaysFail {
				assert.NotNil(t, err, infoMsg)
			} else if action.ShouldErr {
				assert.NotNil(t, err, infoMsg)
				if action.ExpectErr != nil {
					assert.Equal(t, action.ExpectErr, err, infoMsg)
				}
			} else {
				assert.Nil(t, err, infoMsg)
			}
		}
	}
}

func RunBenchmarkPush(b *testing.B, queue blobqueue.Queue) {
	for n := 0; n < b.N; n++ {
		queue.Push([]byte("bench"))
	}
}
