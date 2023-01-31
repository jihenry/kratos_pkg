package util

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"gitlab.yeahka.com/gaas/pkg/cache"
)

func TestRedisLock(t *testing.T) {
	key := RandString(defaultRandLen)
	Rdb, _ := cache.NewRedisClient(
		cache.WithAddr("127.0.0.1:6379"),
		cache.WithDb(1))
	ctx := context.Background()

	firstLock := NewRedisLock(Rdb, key)
	firstLock.SetExpire(5)

	fLock, err := firstLock.Acquire(ctx)
	assert.Nil(t, err)
	assert.True(t, fLock)

	secondLock := NewRedisLock(Rdb, key)
	secondLock.SetExpire(5)
	againAcquire, err := secondLock.Acquire(ctx)
	assert.Nil(t, err)
	assert.False(t, againAcquire)

}
