package rdb

import (
	"fmt"
	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
	"time"
)

type RDB struct {
	Cli       *redis.Client
	Cache     *cache.Cache
	keyPrefix string
}

func New(opt *redis.Options, opts ...Option) *RDB {
	cli := redis.NewClient(opt)
	ca := cache.New(&cache.Options{
		Redis:      cli,
		LocalCache: cache.NewTinyLFU(256, time.Minute),
	})
	r := &RDB{
		Cli:       cli,
		Cache:     ca,
		keyPrefix: "rdb",
	}
	for _, v := range opts {
		v(r)
	}
	return r
}

func (r *RDB) Key(suffix string) string {
	return fmt.Sprintf("%s_%s", r.keyPrefix, suffix)
}
