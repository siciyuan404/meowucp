package redis

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

type Client struct {
	*redis.Client
}

func NewClient(host string, port int, password string, db int, poolSize int) (*Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: password,
		DB:       db,
		PoolSize: poolSize,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	log.Println("Redis connected successfully")

	return &Client{client}, nil
}

func (c *Client) Close() error {
	return c.Client.Close()
}

func (c *Client) Get(ctx context.Context, key string) (string, error) {
	return c.Client.Get(ctx, key).Result()
}

func (c *Client) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return c.Client.Set(ctx, key, value, expiration).Err()
}

func (c *Client) Del(ctx context.Context, keys ...string) error {
	return c.Client.Del(ctx, keys...).Err()
}

func (c *Client) Exists(ctx context.Context, keys ...string) (int64, error) {
	return c.Client.Exists(ctx, keys...).Result()
}

func (c *Client) AddToStream(ctx context.Context, stream string, values map[string]interface{}) (string, error) {
	return c.Client.XAdd(ctx, &redis.XAddArgs{
		Stream: stream,
		Values: values,
	}).Result()
}

func (c *Client) ReadFromStream(ctx context.Context, stream, group, consumer string, count int64, block time.Duration) ([]redis.XStream, error) {
	return c.Client.XReadGroup(ctx, &redis.XReadGroupArgs{
		Group:    group,
		Consumer: consumer,
		Streams:  []string{stream, ">"},
		Count:    count,
		Block:    block,
	}).Result()
}

func (c *Client) AckStream(ctx context.Context, stream, group string, ids ...string) error {
	return c.Client.XAck(ctx, stream, group, ids...).Err()
}
