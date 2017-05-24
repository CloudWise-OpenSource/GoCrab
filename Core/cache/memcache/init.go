// Copyright 2017 GoCrab neeke@php.net All Rights Reserved.
//
// Licensed under the GNU General Public License, Version 3.
//
// Powered By CloudWise

// Usage:
// import(
//   _ "github.com/CloudWise-OpenSource/GoCrab/Core/cache/memcache"
//   "github.com/CloudWise-OpenSource/GoCrab/Core/cache"
// )
//
//  bm, err := cache.NewCache("memcache", `{"conn":"127.0.0.1:11211"}`)
//
package memcache

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/CloudWise-OpenSource/GoCrab/Core/cache"
)

// Memcache adapter.
type MemcacheCache struct {
	conn     *Client
	conninfo []string
}

// create new memcache adapter.
func NewMemCache() *MemcacheCache {
	return &MemcacheCache{}
}

// get stats from memcache.
func (rc *MemcacheCache) Stats() map[string]interface{} {
	if rc.conn == nil {
		if err := rc.connectInit(); err != nil {
			stats := make(map[string]interface{})
			stats["stats"] = "faild"
			return stats
		}
	}
	if stats, err := rc.conn.Stats(); err == nil {
		return stats
	}
	return nil
}

// get value from memcache.
func (rc *MemcacheCache) Get(key string) interface{} {
	if rc.conn == nil {
		if err := rc.connectInit(); err != nil {
			return err
		}
	}
	if item, err := rc.conn.Get(key); err == nil {
		return string(item.Value)
	}
	return nil
}

// put value to memcache. only support string.
func (rc *MemcacheCache) Put(key string, val interface{}, timeout int64) error {
	if rc.conn == nil {
		if err := rc.connectInit(); err != nil {
			return err
		}
	}
	v, ok := val.(string)
	if !ok {
		return errors.New("val must string")
	}
	item := Item{Key: key, Value: []byte(v), Expiration: int32(timeout)}
	return rc.conn.Set(&item)
}

// delete value in memcache.
func (rc *MemcacheCache) Delete(key string) error {
	if rc.conn == nil {
		if err := rc.connectInit(); err != nil {
			return err
		}
	}
	return rc.conn.Delete(key)
}

// increase counter.
func (rc *MemcacheCache) Incr(key string) error {
	if rc.conn == nil {
		if err := rc.connectInit(); err != nil {
			return err
		}
	}
	_, err := rc.conn.Increment(key, 1)
	return err
}

// decrease counter.
func (rc *MemcacheCache) Decr(key string) error {
	if rc.conn == nil {
		if err := rc.connectInit(); err != nil {
			return err
		}
	}
	_, err := rc.conn.Decrement(key, 1)
	return err
}

// check value exists in memcache.
func (rc *MemcacheCache) IsExist(key string) bool {
	if rc.conn == nil {
		if err := rc.connectInit(); err != nil {
			return false
		}
	}
	_, err := rc.conn.Get(key)
	if err != nil {
		return false
	}
	return true
}

// clear all cached in memcache.
func (rc *MemcacheCache) ClearAll() error {
	if rc.conn == nil {
		if err := rc.connectInit(); err != nil {
			return err
		}
	}
	return rc.conn.FlushAll()
}

// start memcache adapter.
// config string is like {"conn":"connection info"}.
// if connecting error, return.
func (rc *MemcacheCache) StartAndGC(config string) error {
	var cf map[string]string
	json.Unmarshal([]byte(config), &cf)
	if _, ok := cf["conn"]; !ok {
		return errors.New("config has no conn key")
	}
	rc.conninfo = strings.Split(cf["conn"], ";")
	if rc.conn == nil {
		if err := rc.connectInit(); err != nil {
			return err
		}
	}
	return nil
}

// connect to memcache and keep the connection.
func (rc *MemcacheCache) connectInit() error {
	rc.conn = New(rc.conninfo...)
	return nil
}

func init() {
	cache.Register("memcache", NewMemCache())
}
