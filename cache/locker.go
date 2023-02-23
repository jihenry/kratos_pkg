package cache

import (
	"context"
	"math/rand"
	"sync/atomic"
	"time"

	"github.com/go-redis/redis/v8"
)

const (
	tolerance       = 500
	millisPerSecond = 1000
	lockCommand     = `if redis.call("GET", KEYS[1]) == ARGV[1] then
    redis.call("SET", KEYS[1], ARGV[1], "PX", ARGV[2])
    return "OK"
else
    return redis.call("SET", KEYS[1], ARGV[1], "NX", "PX", ARGV[2])
end`
	delCommand = `if redis.call("GET", KEYS[1]) == ARGV[1] then
    return redis.call("DEL", KEYS[1])
else
    return 0
end`
	randomLen = 16
)

type RedisLock struct {
	store   *redis.Client
	seconds uint32
	count   int32
	key     string
	id      string
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func NewRedisLock(store *redis.Client, key string) *RedisLock {
	return &RedisLock{
		store: store,
		key:   key,
		id:    Randn(randomLen),
	}
}

// Acquire 设置分布式锁
func (rl *RedisLock) Acquire(ctx context.Context) (bool, error) {
	seconds := atomic.LoadUint32(&rl.seconds)
	//优化代码
	args := make([]interface{}, 0)
	args = append(args, rl.id, int(seconds)*millisPerSecond+tolerance)
	res, err := rl.store.Eval(ctx, lockCommand, []string{rl.key}, args...).Result()
	if err == redis.Nil {
		return false, nil
	} else if err != nil {
		return false, err
	} else if res == nil {
		return false, nil
	}
	reply, ok := res.(string)
	if ok && reply == "OK" {
		return true, nil
	}
	return false, nil
}

// Release 释放锁
func (rl *RedisLock) Release(ctx context.Context) (bool, error) {
	resp, err := rl.store.Eval(ctx, delCommand, []string{rl.key}, []string{rl.id}).Result()
	if err != nil {
		return false, err
	}
	if reply, ok := resp.(int64); ok {
		return reply == 1, nil
	}
	return false, nil
}

// SetExpire 设置过期时间
func (rl *RedisLock) SetExpire(seconds int) {
	atomic.StoreUint32(&rl.seconds, uint32(seconds))
}

const (
	letterBytes    = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	letterIdxBits  = 6
	defaultRandLen = 8
	letterIdxMask  = 1<<letterIdxBits - 1
	letterIdxMax   = 63 / letterIdxBits
)

func Randn(n int) string {
	b := make([]byte, n)
	for i, cache, remain := n-1, rand.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rand.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}
