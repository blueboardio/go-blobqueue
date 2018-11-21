# go-blobqueue

[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/blueboardio/go-blobqueue)
[![Coverage](https://codecov.io/gh/blueboardio/go-blobqueue/branch/master/graph/badge.svg)](https://codecov.io/gh/blueboardio/go-blobqueue/branch/master)
[![Build](https://travis-ci.org/blueboardio/go-blobqueue.svg?branch=master)](https://travis-ci.org/blueboardio/go-blobqueue)

## About

Blobqueue provides an interface used to manipulate a queue of objects marshalled  into `[]byte` (the blob part).
Two implementations are provided:

* `blobqueue/memory`: a pure go native implementation based on a `[][]byte`
* `blobqueue/redisqueue`: an implementation using a Redis list. (using Redis client [go-redis/redis](https://github.com/go-redis/redis))

## Getting started

### 1. Memory

This example shows how to create a basic and safe memory queue. A safe queue wraps another one to properly handles locks.

```golang
import (
	"fmt"

	"github.com/blueboardio/go-blobqueue"
	"github.com/blueboardio/go-blobqueue/memory"
)

func getStarted() {
	// unsafeQueue implements Queue so you could  use it as it is, but it's not safe against race conditions.
	unsafeQueue := memory.Queue{}

	// This queue handles locks properly.
	queue := blobqueue.SafeQueue(&unsafeQueue)
}
```

### 2. Redis

This example show how to create a Redis based queue.

```golang
import (
	"fmt"

	"github.com/blueboardio/go-blobqueue/queueredis"
	"github.com/go-redis/redis"
)

func getStarted() {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	// Instantiate a redis queue storing its data under the key "hello_world_queue"
	queue := queueredis.New(client, "hello_world_queue")
}

```

### 3. Usage

You can now use theses queues like so:
```golang
	queue.Push([]byte("world"))
	queue.Unshift([]byte("hello"))
	list, _ := queue.List()
	fmt.Println(list)
	// Bytes: [[104 101 108 108 111] [119 111 114 108 100]]
	// Strings: ["hello" "world"]

	first, _ := queue.Shift()
	list, _ = queue.List()
	fmt.Println(first, list)
	// Bytes: [104 101 108 108 111] [[119 111 114 108 100]]
	// Strings: "hello" ["world"]
```
See [GoDoc](https://godoc.org/github.com/blueboardio/go-blobqueue) for all available methods.


## Testing

The package `blobqueue/queuetesting` provides a test suite to validate implementation of `bloqueue.Queue`. 

### Running Redis implementation test

Redis implementation depends on a Redis connection, so you'll need to have a working Redis server to pass tests. In order to connect when running tests you'll have to set these environment variables:

|Name|Usage|Example|
|:---|:---:|------:|
|**TEST_QUEUEREDIS_ADDR**| The redis host used to set connection |`localhost:6379`|
|**TEST_QUEUEREDIS_PWD**| The password to use | |
|**TEST_QUEUEREDIS_DB**| The database index to use | `0` |
