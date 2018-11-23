package typedqueue_test

import (
	"encoding/binary"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/blueboardio/go-blobqueue/memory"
	"github.com/blueboardio/go-blobqueue/typedqueue"
)

type Uint64 uint64

func (u Uint64) MarshalBinary() ([]byte, error) {
	b := make([]byte, binary.MaxVarintLen64)
	binary.PutUvarint(b, uint64(u))
	return b, nil
}

func (u *Uint64) UnmarshalBinary(b []byte) error {
	tmp, n := binary.Uvarint(b)
	if n <= 0 {
		return errors.New("invalid data")
	}
	*u = Uint64(tmp)
	return nil
}

func Example() {
	q := typedqueue.New(&memory.Queue{}, Uint64(0))

	_ = q.Push(Uint64(1))
	_ = q.Push(Uint64(2))
	_ = q.Push(Uint64(3))

	li, _ := q.List()
	fmt.Printf("%#v\n", li)

	// Output:
	// []typedqueue_test.Uint64{0x1, 0x2, 0x3}
}

func TestQueue(t *testing.T) {
	assert := assert.New(t)
	q := typedqueue.New(&memory.Queue{}, Uint64(0))

	li, err := q.List()
	assert.NoError(err)
	assert.Equal(li, []Uint64(nil))

	l, err := q.Len()
	assert.NoError(err)
	assert.Equal(l, 0)

	q.Push(Uint64(1))
	l, err = q.Len()
	assert.NoError(err)
	assert.Equal(l, 1)

	q.Push(Uint64(2))
	l, err = q.Len()
	assert.NoError(err)
	assert.Equal(l, 2)

	q.Push(Uint64(3))
	l, err = q.Len()
	assert.NoError(err)
	assert.Equal(l, 3)

	li, err = q.List()
	assert.NoError(err)
	assert.Equal(li, []Uint64{1, 2, 3})

	n, err := q.Pop()
	assert.NoError(err)
	assert.Equal(n, Uint64(3))

	li, err = q.List()
	assert.NoError(err)
	assert.Equal(li, []Uint64{1, 2})

	n, err = q.Shift()
	assert.NoError(err)
	assert.Equal(n, Uint64(1))

	li, err = q.List()
	assert.NoError(err)
	assert.Equal(li, []Uint64{2})
	l, err = q.Len()
	assert.NoError(err)
	assert.Equal(l, 1)
}
