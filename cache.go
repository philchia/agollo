package agollo

import (
	"sync"
)

type cache struct {
	kv sync.Map
}

func newCache() *cache {
	return &cache{
		kv: sync.Map{},
	}
}

func (c *cache) Set(key, val string) {
	c.kv.Store(key, val)
}

func (c *cache) Get(key string) (string, bool) {
	if val, ok := c.kv.Load(key); ok {
		if ret, ok := val.(string); ok {
			return ret, true
		}
	}
	return "", false
}

func (c *cache) Delete(key string) {
	c.kv.Delete(key)
}

func (c *cache) All() map[string]string {
	var ret = map[string]string{}
	c.kv.Range(func(key, val interface{}) bool {
		if key, ok := key.(string); ok {
			if val, ok := val.(string); ok {
				ret[key] = val
			}
		}
		return true
	})
	return ret
}
