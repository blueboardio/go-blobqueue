# go-blobqueue

[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/blueboardio/go-blobqueue)
[![Coverage](https://codecov.io/gh/blueboardio/go-blobqueue/branch/master/graph/badge.svg)](https://codecov.io/gh/blueboardio/go-blobqueue/branch/master)
[![Build](https://travis-ci.org/blueboardio/go-blobqueue.svg?branch=master)](https://travis-ci.org/blueboardio/go-blobqueue)

## About

Blobqueue provides an interface used to manipulate a queue of objects marshalled  into `[]byte` (the blob part).
Two implementations are provided:

* `blobqueue/memory`: a pure go native implementation based on a `[][]byte`
* `blobqueue/redisqueue`: an implementation using a Redis list. (using Redis client [go-redis/redis](https://github.com/go-redis/redis))

## Testing

The package `blobqueue/queuetesting` provides a test suite to validate implementation of `bloqueue.Queue`. 

### Running Redis implementation test

Redis implementation depends on a Redis connection, so you'll need to have a working Redis server to pass tests. In order to connect when running tests you'll have to set these environment variables:

|Name|Usage|Example|
|:---|:---:|------:|
|**TEST_QUEUEREDIS_ADDR**| The redis host used to set connection |`localhost:6379`|
|**TEST_QUEUEREDIS_PWD**| The password to use | |
|**TEST_QUEUEREDIS_DB**| The database index to use | `0` |

## SafeQueue usage

Blobqueue provides a constructor:
```golang
func SafeQueue(q Queue) Queue
```

It wraps `q` with an higher level queue handling locks to avoid race condition. To instantiate thread safe memory and/or Redis implementations you should:

```golang
// Memory
queue := blobqueue.SafeQueue(&memory.Queue{})

// Redis
qRedis := queueredis.New(client, "test_key")
queue := blobqueue.SafeQueue(qRedis)
```
