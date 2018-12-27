package agollo

import (
	"fmt"
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
	var caches = newNamespaceCahce()
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

	var restore = newNamespaceCahce()
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

func TestDecode(t *testing.T) {
	// decode
	type Client struct {
		APP struct {
			Key string `mapstructure:"key"`
		} `mapstructure:"application"`
		Client1 struct {
			Name string `mapstructure:"name"`
		} `mapstructure:"client.json"`
		Client2 struct {
			Name string `mapstructure:"name"`
		} `mapstructure:"client.yaml"`
	}
	var mc = newNamespaceCahce()
	defer mc.drain()
	mc.mustGetCache("application").set("key", "val")
	mc.mustGetCache("client.json").set("content", `{"name":"json"}`)
	mc.mustGetCache("client.yaml").set("content", "name: yaml")
	var c Client
	if err := mc.decode(&c); err != nil {
		t.Error(err)
	}

	fmt.Printf("%+v", c)
	if c.APP.Key != "val" || c.Client1.Name != "json" || c.Client2.Name != "yaml" {
		t.FailNow()
	}

}
