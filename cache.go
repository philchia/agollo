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

func (c *cache) set(key, val string) {
	c.kv.Store(key, val)
}

func (c *cache) get(key string) (string, bool) {
	if val, ok := c.kv.Load(key); ok {
		if ret, ok := val.(string); ok {
			return ret, true
		}
	}
	return "", false
}

func (c *cache) delete(key string) {
	c.kv.Delete(key)
}

func (c *cache) dump() map[string]string {
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
