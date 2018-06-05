package agollo

import (
	"sync"
)

// Client for apollo
type Client struct {
	updateChan chan *ChangeEvent

	mutex  sync.RWMutex
	caches map[string]*cache
}

// NewClient create client from conf
func NewClient(conf *Conf) (*Client, error) {
	client := &Client{
		updateChan: make(chan *ChangeEvent),
		caches:     map[string]*cache{},
	}
	return client, nil
}

// Start sync config
func (c *Client) Start() error {
	return nil
}

// Stop sync config
func (c *Client) Stop() error {
	close(c.updateChan)
	return nil
}

// WatchUpdate get all updates
func (c *Client) WatchUpdate() <-chan *ChangeEvent {
	return c.updateChan
}

func (c *Client) mustGetCache(namespace string) *cache {
	c.mutex.RLock()
	if ret, ok := c.caches[namespace]; ok {
		c.mutex.RUnlock()
		return ret
	}

	cache := newCache()
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.caches[namespace] = cache
	return cache
}

// GetStringValueWithNameSapce get value from given namespace
func (c *Client) GetStringValueWithNameSapce(namespace, key, defaultValue string) string {
	cache := c.mustGetCache(namespace)
	if ret, ok := cache.Get(key); ok {
		return ret
	}
	return defaultValue
}

// GetStringValue from default namespace
func (c *Client) GetStringValue(key, defaultValue string) string {
	return c.GetStringValueWithNameSapce(defaultNamespace, key, defaultValue)
}
