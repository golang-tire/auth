package testutils

import (
	"fmt"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
	"github.com/golang-tire/pkg/pubsub"
)

var testRedisServer *miniredis.Miniredis

func TestUp() {
	var err error
	testRedisServer, err = miniredis.Run()
	if err != nil {
		panic(err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", testRedisServer.Host(), testRedisServer.Port()),
	})

	pubsub.New(rdb)
}

func TestDown() {
	testRedisServer.Close()
}
