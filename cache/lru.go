package cache

import (
	"errors"
	"log"
	"time"

	"github.com/creachadair/cityhash"

	lru "github.com/hashicorp/golang-lru"
)

var (
	Nil        = errors.New("cache: nil")
	RemoveNil  = errors.New("remove:nil")
	AddNil     = errors.New("add:nil")
	ErrExpTime = errors.New("args error, expTime must be > 0")
	dbs        map[int]*lru.Cache
	ldb        *lru.Cache
)

func InitLruCache(size int) {
	ldb, _ = lru.New(size)
}

func LSet(key string, value interface{}) {
	ldb.Add(key, value)
}

func LGet(key string) (interface{}, bool) {
	return ldb.Get(key)
}

func Set(key string, value interface{}) error {
	return SetEx(key, value, 86400*time.Second)
}

func SetEx(key string, value interface{}, exp time.Duration) error {
	return setEx(key, value, exp)
}

func Get(key string) (interface{}, error) {
	return get(key)
}

func Del(key string) error {
	return del(key)
}

func init() {
	if ldb == nil {
		dbs = make(map[int]*lru.Cache, 1024)
		for i := 0; i < 1024; i++ {
			dbs[i], _ = lru.New(3000)
		}
	}
	go flush()
}

type data struct {
	d interface{}
	// 过期的时间点
	t time.Time
}

func delByDb(key string, db *lru.Cache) error {
	if db.Remove(key) {
		return nil
	}
	return RemoveNil
}

func getByDB(key string, db *lru.Cache) (interface{}, error) {
	if value, ok := db.Get(key); ok {
		if vv, ok := value.(*data); ok {
			if !vv.t.After(time.Now()) {
				if err := delByDb(key, db); err != nil {
					return nil, err
				}
			}
			return vv.d, nil
		}
	}
	return nil, Nil
}

func flush() {
	i := 0
	for range time.NewTicker(5 * time.Second).C {
		shard := i % 1024
		db, ok := dbs[shard]
		if !ok {
			continue
		}
		for _, k := range db.Keys() {
			if ks, ko := k.(string); ko {
				if _, err := getByDB(ks, db); err != nil {
					log.Printf("flush err:%v", err)
				}
			}
		}
		i++
	}
}

func setEx(key string, value interface{}, exp time.Duration) error {
	// 缓存的最小单位是1ms
	if exp < time.Millisecond {
		return ErrExpTime
	}
	if db, ok := dbs[hash(key)]; ok {
		db.Add(key, &data{d: value, t: time.Now().Add(exp)})
	}
	return nil
}

func hash(key string) int {
	return int(cityhash.Hash32([]byte(key)) % 1024)
}

func get(key string) (interface{}, error) {
	if db, ok := dbs[hash(key)]; ok {
		return getByDB(key, db)
	}
	return nil, Nil
}

func del(key string) error {
	if db, ok := dbs[hash(key)]; ok {
		return delByDb(key, db)
	}
	return Nil
}
