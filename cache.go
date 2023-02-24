package agollo

import (
	"encoding/gob"
	"os"
	"sync"
)

type namespaceCache struct {
	lock   sync.Mutex
	caches map[string]*cache
}

func newNamespaceCache() *namespaceCache {
	return &namespaceCache{
		caches: map[string]*cache{},
	}
}

func (n *namespaceCache) mustGetCache(namespace string) *cache {
	n.lock.Lock()
	defer n.lock.Unlock()

	if ret, ok := n.caches[namespace]; ok {
		return ret
	}

	cache := newCache()
	n.caches[namespace] = cache
	return cache
}

func (n *namespaceCache) drain() {
	n.lock.Lock()
	defer n.lock.Unlock()

	for namespace := range n.caches {
		delete(n.caches, namespace)
	}
}

func (n *namespaceCache) cacheDump() map[string]map[string]string {
	n.lock.Lock()
	defer n.lock.Unlock()

	var dumps = map[string]map[string]string{}

	for namespace, cache := range n.caches {
		dumps[namespace] = cache.dump()
	}
	return dumps
}

func (n *namespaceCache) dump(name string) error {

	dumps := n.cacheDump()
	tmp := name + "tmp"
	f, err := os.OpenFile(tmp, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}

	if err := gob.NewEncoder(f).Encode(&dumps); err != nil {
		_ = f.Close()
		return err
	}

	if err := f.Close(); err != nil {
		return err
	}

	return os.Rename(tmp, name)
}

func (n *namespaceCache) load(name string) error {
	n.drain()

	f, err := os.OpenFile(name, os.O_RDONLY, 0755)
	if err != nil {
		return err
	}
	defer f.Close()

	var dumps = map[string]map[string]string{}

	if err := gob.NewDecoder(f).Decode(&dumps); err != nil {
		return err
	}

	for namespace, kv := range dumps {
		cache := n.mustGetCache(namespace)
		for k, v := range kv {
			cache.set(k, v)
		}
	}

	return nil
}

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
		k, _ := key.(string)
		v, _ := val.(string)
		ret[k] = v

		return true
	})
	return ret
}
