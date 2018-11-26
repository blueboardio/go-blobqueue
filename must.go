package blobqueue

// Must wraps a Queue with methods that panic instead of returning error.
type Must struct {
	Queue Queue
}

// List transforms errors from Queue.List into panic.
func (q Must) List() [][]byte {
	list, err := q.Queue.List()
	if err != nil {
		panic(err)
	}
	return list
}

// Push transforms errors from Queue.Push into panic.
func (q Must) Push(val []byte) {
	if err := q.Queue.Push(val); err != nil {
		panic(err)
	}
}

// Unshift transforms errors from Queue.Unshift into panic.
func (q Must) Unshift(val []byte) {
	if err := q.Queue.Unshift(val); err != nil {
		panic(err)
	}
}

// Unshift transforms errors from Queue.Unshift into panic.
func (q Must) Pop() []byte {
	val, err := q.Queue.Pop()
	if err != nil {
		panic(err)
	}
	return val
}

// Shift transforms errors from Queue.Shift into panic.
func (q Must) Shift() []byte {
	val, err := q.Queue.Shift()
	if err != nil {
		panic(err)
	}
	return val
}

// Len transforms errors from Queue.Len into panic.
func (q Must) Len() int {
	l, err := q.Queue.Len()
	if err != nil {
		panic(err)
	}
	return l
}

// Empty transforms errors from Queue.Empty into panic.
func (q Must) Empty() {
	if err := q.Queue.Empty(); err != nil {
		panic(err)
	}
}
