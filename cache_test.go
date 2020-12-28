package agollo

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestCache(t *testing.T) {
	cache := newCache()

	cache.set("key", "val")
	if val, ok := cache.get("key"); !ok || val != "val" {
		t.FailNow()
	}

	cache.set("key", "val2")
	if val, ok := cache.get("key"); !ok || val != "val2" {
		t.FailNow()
	}

	kv := cache.dump()
	if len(kv) != 1 || kv["key"] != "val2" {
		t.FailNow()
	}

	cache.delete("key")
	if _, ok := cache.get("key"); ok {
		t.FailNow()
	}
}

func TestCacheDump(t *testing.T) {
	var caches = newNamespaceCache()
	defer caches.drain()
	caches.mustGetCache("namespace").set("key", "val")

	f, err := ioutil.TempFile(".", "agollo")
	if err != nil {
		t.Error(err)
	}
	f.Close()
	defer os.Remove(f.Name())

	if err := caches.dump(f.Name()); err != nil {
		t.Error(err)
	}

	var restore = newNamespaceCache()
	defer restore.drain()
	if err := restore.load(f.Name()); err != nil {
		t.Error(err)
	}

	if val, _ := restore.mustGetCache("namespace").get("key"); val != "val" {
		t.FailNow()
	}

	if err := restore.load("null"); err == nil {
		t.FailNow()
	}

	if err := restore.load("./testdata/app.properties"); err == nil {
		t.FailNow()
	}
}
