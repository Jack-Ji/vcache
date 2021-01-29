package vcache

import (
	"context"
	"strings"
	"testing"
	"time"
)

type Foo struct {
	X int
	Y string
}

func TestCache(t *testing.T) {
	cache, err := NewCache("127.0.0.1:6379", "", "")
	if err != nil {
		t.Fatal(err)
	}

	_, exists, err := cache.Get(context.Background(), "myprefix", "__not_exists")
	if err != nil {
		t.Error("read fail:", err)
	}
	if exists {
		t.Error("key `__not_exists` should not exist")
	}

	err = cache.SetWithLifeTime(context.Background(), "myprefix", "hello", "mydata", time.Second)
	if err != nil {
		t.Error("write cache failed", err)
	}
	val, exists, err := cache.Get(context.Background(), "myprefix", "hello")
	if err != nil {
		t.Error("read cache failed", err)
	}
	if !exists {
		t.Error("cache is gone accidentally")
	}
	if val != "mydata" {
		t.Error("cache's content is invalid:", val)
	}

	time.Sleep(time.Second * 2)
	_, exists, err = cache.Get(context.Background(), "myprefix", "hello")
	if err != nil {
		t.Error("read cache failed", err)
	}
	if exists {
		t.Error("cache is not gone")
	}

	cache.MustSet("myprefix", &Foo{X: 3, Y: "world"}, "mydata")

	_, exists, err = cache.Get(context.Background(), "myprefix", &Foo{X: 4, Y: "world"})
	if err != nil {
		t.Error("read cache failed", err)
	}
	if exists {
		t.Error("cache should not exist")
	}

	val, exists, err = cache.Get(context.Background(), "myprefix", &Foo{X: 3, Y: "world"})
	if err != nil {
		t.Error("read cache failed", err)
	}
	if !exists {
		t.Error("cache is gone accidentally")
	}
	if val != "mydata" {
		t.Error("cache's content is invalid:", val)
	}

	err = cache.Del(context.Background(), "myprefix", &Foo{X: 3, Y: "world"})
	if err != nil {
		t.Error("remove cache failed:", err)
	}

	_, exists, err = cache.Get(context.Background(), "myprefix", &Foo{X: 3, Y: "world"})
	if err != nil {
		t.Error("read cache failed", err)
	}
	if exists {
		t.Error("cache should not exist")
	}
}

func BenchmarkCacheSet(b *testing.B) {
	cache, err := NewCache("127.0.0.1:6379", "", "")
	if err != nil {
		panic(err)
	}

	key := &Foo{X: 3, Y: "world"}
	val := strings.Repeat("a", 64*1000)
	for i := 0; i < b.N; i++ {
		cache.MustSet("myprefix", key, val)
	}
}

func BenchmarkCacheGet(b *testing.B) {
	cache, err := NewCache("127.0.0.1:6379", "", "")
	if err != nil {
		panic(err)
	}

	key := &Foo{X: 3, Y: "world"}
	val := strings.Repeat("b", 64*1000)
	cache.MustSet("myprefix", key, val)

	for i := 0; i < b.N; i++ {
		_, exists, err := cache.Get(context.Background(), "myprefix", key)
		if err != nil || !exists {
			panic(err)
		}
	}
}
