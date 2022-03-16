package cache

import (
	"fmt"

	"github.com/alicebob/miniredis/v2"
	goRedis "github.com/go-redis/redis/v8"
)

// NewTest returns a Cache instance over miniredis
func NewTest() Cache {
	_, client := NewMiniredis()
	return client
}

// NewMiniredis returns a miniredis and Cache instance for test purposes
func NewMiniredis() (*miniredis.Miniredis, Cache) {
	mr, _ := miniredis.Run()
	client := goRedis.NewClient(&goRedis.Options{
		Addr:     fmt.Sprintf("%s:%s", mr.Host(), mr.Port()),
		Password: "",
		DB:       0,
	})
	return mr, &cache{
		client:    client,
		keyPrefix: "test",
	}
}
