// Copyright 2017 GoCrab neeke@php.net All Rights Reserved.
//
// Licensed under the GNU General Public License, Version 3.
//
// Powered By CloudWise

package redis

import (
	"testing"
	"time"

	"github.com/garyburd/redigo/redis"

	"github.com/CloudWise-OpenSource/GoCrab/Core/cache"
)

func TestRedisCache(t *testing.T) {
	bm, err := cache.NewCache("redis", `{"conn": "127.0.0.1:6379"}`)
	if err != nil {
		t.Error("init err")
	}
	if err = bm.Put("Neeke", 1, 10); err != nil {
		t.Error("set Error", err)
	}
	if !bm.IsExist("Neeke") {
		t.Error("check err")
	}

	time.Sleep(10 * time.Second)

	if bm.IsExist("Neeke") {
		t.Error("check err")
	}
	if err = bm.Put("Neeke", 1, 10); err != nil {
		t.Error("set Error", err)
	}

	if v, _ := redis.Int(bm.Get("Neeke"), err); v != 1 {
		t.Error("get err")
	}

	if err = bm.Incr("Neeke"); err != nil {
		t.Error("Incr Error", err)
	}

	if v, _ := redis.Int(bm.Get("Neeke"), err); v != 2 {
		t.Error("get err")
	}

	if err = bm.Decr("Neeke"); err != nil {
		t.Error("Decr Error", err)
	}

	if v, _ := redis.Int(bm.Get("Neeke"), err); v != 1 {
		t.Error("get err")
	}
	bm.Delete("Neeke")
	if bm.IsExist("Neeke") {
		t.Error("delete err")
	}
	//test string
	if err = bm.Put("Neeke", "author", 10); err != nil {
		t.Error("set Error", err)
	}
	if !bm.IsExist("Neeke") {
		t.Error("check err")
	}

	if v, _ := redis.String(bm.Get("Neeke"), err); v != "author" {
		t.Error("get err")
	}
	// test clear all
	if err = bm.ClearAll(); err != nil {
		t.Error("clear all err")
	}
}
