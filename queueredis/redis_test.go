package queueredis_test

import (
	"os"
	"strconv"
	"testing"

	"github.com/blueboardio/go-blobqueue"

	"github.com/go-redis/redis"
	"github.com/stretchr/testify/assert"

	"github.com/blueboardio/go-blobqueue/queueredis"
	"github.com/blueboardio/go-blobqueue/queuetesting"
)

const (
	testKey     = "testQueue"
	benchKey    = "benchmarkQueue"
	testFailKey = "testFailQueue"
)

func instantiateRedisClient() (*redis.Client, error) {
	db, err := strconv.Atoi(os.Getenv("TEST_QUEUEREDIS_DB"))
	if err != nil {
		return nil, err
	}
	return redis.NewClient(&redis.Options{
		Addr:     os.Getenv("TEST_QUEUEREDIS_ADDR"),
		Password: os.Getenv("TEST_QUEUEREDIS_PWD"),
		DB:       db,
	}), nil
}

func TestRedisImplem(t *testing.T) {
	client, err := instantiateRedisClient()
	assert.Nil(t, err)
	defer client.Close()
	qRedis := queueredis.New(client, testKey)
	err = client.Del(testKey).Err()
	assert.Nil(t, err, "unable to clear queue before redis implem test")
	queue := blobqueue.SafeQueue(qRedis)
	queuetesting.RunTests(t, queue, false)
}

func TestRedisBadConnection(t *testing.T) {
	client, err := instantiateRedisClient()
	assert.Nil(t, err)
	defer client.Close()
	qRedis := queueredis.New(client, testFailKey)
	client.Close()
	queue := blobqueue.SafeQueue(qRedis)
	queuetesting.RunTests(t, queue, true)
}

func BenchmarkRedisPush(b *testing.B) {
	client, err := instantiateRedisClient()
	if err != nil {
		b.Fail()
	}
	defer client.Close()
	qRedis := queueredis.New(client, benchKey)
	err = client.Del(benchKey).Err()
	if err != nil {
		b.Fail()
	}
	queue := blobqueue.SafeQueue(qRedis)
	queuetesting.RunBenchmarkPush(b, queue)
}
