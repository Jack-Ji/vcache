package vcache

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type Cache struct {
	cli      *redis.Client // redis connection
	hasher   Hasher        // default hasher
	expireMs int           // default key expire milliseconds
}

// Create cache handle
func NewCache(addr, user, password string) (*Cache, error) {
	cli := redis.NewClient(&redis.Options{
		Addr:     addr,
		Username: user,
		Password: password,
		DB:       9,
	})

	// detect connection
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	_, err := cli.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return &Cache{
		cli:      cli,
		hasher:   HashFunc(MustMD5),
		expireMs: 1000 * 60,
	}, nil
}

// Write data with default expire time (1 minute)
func (c *Cache) Set(ctx context.Context, key interface{}, value string) error {
	s := c.hasher.ToString(key)
	_, err := c.cli.Set(ctx, s, value, time.Duration(c.expireMs*int(time.Millisecond))).Result()
	if err != nil {
		return err
	}
	return nil
}

// Write data with explicit expire time
func (c *Cache) SetWithLifeTime(ctx context.Context, key interface{}, value string, lifeTime time.Duration) error {
	s := c.hasher.ToString(key)
	_, err := c.cli.Set(ctx, s, value, lifeTime).Result()
	if err != nil {
		return err
	}
	return nil
}

// Write data with default expire time (1 minute), succeed or panic
func (c *Cache) MustSet(key interface{}, value string) {
	s := c.hasher.ToString(key)
	_, err := c.cli.Set(context.Background(), s, value,
		time.Duration(c.expireMs*int(time.Millisecond))).Result()
	if err != nil {
		panic(err)
	}
}

// Write data with explicit expire time, succeed or panic
func (c *Cache) MustSetWithLifeTime(key interface{}, value string, lifeTime time.Duration) error {
	s := c.hasher.ToString(key)
	_, err := c.cli.Set(context.Background(), s, value, lifeTime).Result()
	if err != nil {
		return err
	}
	return nil
}

// Read data
func (c *Cache) Get(ctx context.Context, key interface{}) (val string, exists bool, err error) {
	s := c.hasher.ToString(key)
	val, err = c.cli.Get(ctx, s).Result()
	if err != nil {
		if err == redis.Nil {
			exists = false
			err = nil
			return
		}
		return
	}

	exists = true
	return
}

// Remove data
func (c *Cache) Del(ctx context.Context, key interface{}) (err error) {
	s := c.hasher.ToString(key)
	rt := c.cli.Del(ctx, s)
	return rt.Err()
}
