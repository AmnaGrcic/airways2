package redis

import (
	"context"

	"github.com/go-redis/redis/v8"
)

var rdb *redis.Client

func StartRedis() {

	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

}

func SetRedisKey(key string, value string, expiration int) error {
	var ctx = context.Background()
	err := rdb.Set(ctx, key, value, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

func GetRedisKey(key string) (string, error) {
	var ctx = context.Background()
	value, err := rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return value, nil
	}
	if err != nil {
		return value, err
	}
	return value, nil
}

func DeleteRedisKey(keys ...string) error {
	var ctx = context.Background()

	for _, k := range keys {

		err := rdb.Del(ctx, k).Err()
		if err != nil {
			return err
		}
	}
	return nil
}
