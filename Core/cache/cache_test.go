// Copyright 2017 GoCrab neeke@php.net All Rights Reserved.
//
// Licensed under the GNU General Public License, Version 3.
//
// Powered By CloudWise


package cache

import (
	"os"
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	bm, err := NewCache("memory", `{"interval":20}`)
	if err != nil {
		t.Error("init err")
	}
	if err = bm.Put("Neeke", 1, 10); err != nil {
		t.Error("set Error", err)
	}
	if !bm.IsExist("Neeke") {
		t.Error("check err")
	}

	if v := bm.Get("Neeke"); v.(int) != 1 {
		t.Error("get err")
	}

	time.Sleep(30 * time.Second)

	if bm.IsExist("Neeke") {
		t.Error("check err")
	}

	if err = bm.Put("Neeke", 1, 10); err != nil {
		t.Error("set Error", err)
	}

	if err = bm.Incr("Neeke"); err != nil {
		t.Error("Incr Error", err)
	}

	if v := bm.Get("Neeke"); v.(int) != 2 {
		t.Error("get err")
	}

	if err = bm.Decr("Neeke"); err != nil {
		t.Error("Decr Error", err)
	}

	if v := bm.Get("Neeke"); v.(int) != 1 {
		t.Error("get err")
	}
	bm.Delete("Neeke")
	if bm.IsExist("Neeke") {
		t.Error("delete err")
	}
}

func TestFileCache(t *testing.T) {
	bm, err := NewCache("file", `{"CachePath":"cache","FileSuffix":".bin","DirectoryLevel":2,"EmbedExpiry":0}`)
	if err != nil {
		t.Error("init err")
	}
	if err = bm.Put("Neeke", 1, 10); err != nil {
		t.Error("set Error", err)
	}
	if !bm.IsExist("Neeke") {
		t.Error("check err")
	}

	if v := bm.Get("Neeke"); v.(int) != 1 {
		t.Error("get err")
	}

	if err = bm.Incr("Neeke"); err != nil {
		t.Error("Incr Error", err)
	}

	if v := bm.Get("Neeke"); v.(int) != 2 {
		t.Error("get err")
	}

	if err = bm.Decr("Neeke"); err != nil {
		t.Error("Decr Error", err)
	}

	if v := bm.Get("Neeke"); v.(int) != 1 {
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

	if v := bm.Get("Neeke"); v.(string) != "author" {
		t.Error("get err")
	}
	os.RemoveAll("cache")
}
