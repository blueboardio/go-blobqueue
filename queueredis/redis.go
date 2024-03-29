package queueredis

import (
	"context"

	"github.com/blueboardio/go-blobqueue"
	"github.com/redis/go-redis/v9"
)

// Queue is the Redis based implementation of blobqueue.Queue interface
type Queue struct {
	client *redis.Client
	key    string
}

// New returns a new Redis queue.
// The `key` parameter gives the key under which the queue will be stored in Redis.
func New(client *redis.Client, key string) *Queue {
	return &Queue{client: client, key: key}
}

// List implements `Queue` interface. It returns all the elements of the queue.
func (q Queue) List() (ret [][]byte, err error) {
	cmd := q.client.LRange(context.Background(), q.key, 0, -1)
	if cmd.Err() != nil {
		return nil, cmd.Err()
	}
	val := cmd.Val()
	for _, v := range val {
		ret = append(ret, []byte(v))
	}
	return ret, err
}

// Push implements `Queue` interface.
// It appends an item at the end of Redis list.
func (q Queue) Push(val []byte) error {
	return q.client.RPush(context.Background(), q.key, val).Err()
}

// Unshift implements `Queue` interface. It inserts an item at the beggining of Redis list.
func (q Queue) Unshift(val []byte) error {
	return q.client.LPush(context.Background(), q.key, val).Err()
}

func removeBase(action func(ctx context.Context, key string) *redis.StringCmd, key string) ([]byte, error) {
	val, err := action(context.Background(), key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, blobqueue.ErrQueueIsEmpty
		}
		return nil, err
	}
	return val, err
}

// Pop implements `Queue` interface. It removes and returns the last elem of the Redis list.
func (q Queue) Pop() ([]byte, error) {
	return removeBase(q.client.RPop, q.key)
}

// Shift implements `Queue` interface. It removes and returns the first elem of the Redis list.
func (q Queue) Shift() ([]byte, error) {
	return removeBase(q.client.LPop, q.key)
}

/*
// Peek implements `Queue` interface. It returns the first elem of the Redis list without removing it.
// It returns ErrQueueIsEmpty if the list is empty.
func (q Queue) Peek() ([]byte, error) {
	cmd := q.client.LRange(context.Background(), q.key, 0, 0)
	if cmd.Err() != nil {
		return nil, cmd.Err()
	}
	val := cmd.Val()
	if len(val) == 0 {
		return nil, blobqueue.ErrQueueIsEmpty
	}
	return []byte(val[0]), nil
}
*/

// Len implements `Queue` interface. It returns the length of the Redis list.
func (q Queue) Len() (int, error) {
	cmd := q.client.LLen(context.Background(), q.key)
	return int(cmd.Val()), cmd.Err()
}

// Empty implements `Queue` interface. It deletes all elements of the queue.
func (q Queue) Empty() error {
	return q.client.Del(context.Background(), q.key).Err()
}
