package redisstore

import (
	"context"

	"github.com/go-redis/redis/v8"
)

type RDStore struct {
	rdb *redis.Client
}

func New() *RDStore {
	c := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return &RDStore{
		rdb: c,
	}
}

func (rd *RDStore) Insert(score float64, user string) error {
	err := rd.rdb.ZIncrBy(context.TODO(), "users", score, user).Err()
	if err != nil {
		panic(err)
	}
	return err
}
