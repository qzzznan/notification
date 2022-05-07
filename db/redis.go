package db

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"reflect"
)

var rdb *redis.Client
var ctx context.Context

func InitRedis() error {
	ctx = context.Background()
	rdb = redis.NewClient(&redis.Options{
		Addr:     viper.GetString("redis.addr"),
		Password: "",
		DB:       0,
	})

	return rdb.Ping(context.Background()).Err()
}

func SetStruct(key string, v interface{}) error {
	rv := reflect.Indirect(reflect.ValueOf(v))
	rt := rv.Type()
	n := rt.NumField()
	m := make(map[string]interface{})
	for i := 0; i < n; i++ {
		fv := rv.Field(i)
		tg := rt.Field(i).Tag.Get("redis")
		if tg != "" && fv.CanInterface() {
			m[tg] = fv.Interface()
		}
	}

	return rdb.HSet(ctx, key, m).Err()
}

func GetStruct(key string, v interface{}) error {
	return rdb.HGetAll(ctx, key).Scan(v)
}
