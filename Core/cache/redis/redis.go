// Copyright 2017 GoCrab neeke@php.net All Rights Reserved.
//
// Licensed under the GNU General Public License, Version 3.
//
// Powered By CloudWise


// package redis for cache provider
//
// depend on github.com/garyburd/redigo/redis
//
// go install github.com/garyburd/redigo/redis
//
// Usage:
// import(
//   _ "github.com/CloudWise-OpenSource/GoCrab/Core/cache/redis"
//   "github.com/CloudWise-OpenSource/GoCrab/Core/cache"
// )
//
//  bm, err := cache.NewCache("redis", `{"conn":"127.0.0.1:11211"}`)
//
package redis

import (
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/CloudWise-OpenSource/GoCrab/Libs/redigo/redis"

	"github.com/CloudWise-OpenSource/GoCrab/Core/cache"
)

var (
	// the collection name of redis for cache adapter.
	DefaultKey string = "NeekeGao"
)

// Redis cache adapter.
type RedisCache struct {
	p        *redis.Pool // redis connection pool
	conninfo string
	dbNum    int
	key      string
}

// create new redis cache with default collection name.
func NewRedisCache() *RedisCache {
	return &RedisCache{key: DefaultKey}
}

// actually do the redis cmds
func (rc *RedisCache) do(commandName string, args ...interface{}) (reply interface{}, err error) {
	c := rc.p.Get()
	defer c.Close()

	return c.Do(commandName, args...)
}

// Get cache from redis.
func (rc *RedisCache) Get(key string) interface{} {
	if v, err := rc.do("GET", key); err == nil {
		return v
	}
	return nil
}

// put cache to redis.
func (rc *RedisCache) Put(key string, val interface{}, timeout int64) error {
	var err error
	if _, err = rc.do("SETEX", key, timeout, val); err != nil {
		return err
	}

	if _, err = rc.do("HSET", rc.key, key, true); err != nil {
		return err
	}
	return err
}

// delete cache in redis.
func (rc *RedisCache) Delete(key string) error {
	var err error
	if _, err = rc.do("DEL", key); err != nil {
		return err
	}
	_, err = rc.do("HDEL", rc.key, key)
	return err
}

// check cache's existence in redis.
func (rc *RedisCache) IsExist(key string) bool {
	v, err := redis.Bool(rc.do("EXISTS", key))
	if err != nil {
		return false
	}
	if v == false {
		if _, err = rc.do("HDEL", rc.key, key); err != nil {
			return false
		}
	}
	return v
}

// increase counter in redis.
func (rc *RedisCache) Incr(key string) error {
	_, err := redis.Bool(rc.do("INCRBY", key, 1))
	return err
}

// decrease counter in redis.
func (rc *RedisCache) Decr(key string) error {
	_, err := redis.Bool(rc.do("INCRBY", key, -1))
	return err
}

// clean all cache in redis. delete this redis collection.
func (rc *RedisCache) ClearAll() error {
	cachedKeys, err := redis.Strings(rc.do("HKEYS", rc.key))
	if err != nil {
		return err
	}
	for _, str := range cachedKeys {
		if _, err = rc.do("DEL", str); err != nil {
			return err
		}
	}
	_, err = rc.do("DEL", rc.key)
	return err
}

// start redis cache adapter.
// config is like {"key":"collection key","conn":"connection info","dbNum":"0"}
// the cache item in redis are stored forever,
// so no gc operation.
func (rc *RedisCache) StartAndGC(config string) error {
	var cf map[string]string
	json.Unmarshal([]byte(config), &cf)

	if _, ok := cf["key"]; !ok {
		cf["key"] = DefaultKey
	}

	if _, ok := cf["conn"]; !ok {
		return errors.New("config has no conn key")
	}
	if _, ok := cf["dbNum"]; !ok {
		cf["dbNum"] = "0"
	}
	rc.key = cf["key"]
	rc.conninfo = cf["conn"]
	rc.dbNum, _ = strconv.Atoi(cf["dbNum"])
	rc.connectInit()

	c := rc.p.Get()
	defer c.Close()

	return c.Err()
}

// connect to redis.
func (rc *RedisCache) connectInit() {
	dialFunc := func() (c redis.Conn, err error) {
		c, err = redis.Dial("tcp", rc.conninfo)
		_, selecterr := c.Do("SELECT", rc.dbNum)
		if selecterr != nil {
			c.Close()
			return nil, selecterr
		}
		return
	}
	// initialize a new pool
	rc.p = &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 180 * time.Second,
		Dial:        dialFunc,
	}
}

func init() {
	cache.Register("redis", NewRedisCache())
}
