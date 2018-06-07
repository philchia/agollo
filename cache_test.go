package agollo

import "testing"

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
