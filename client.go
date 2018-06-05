package agollo

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sync"
)

// Client for apollo
type Client struct {
	appID      string
	cluster    string
	ip         string
	namespaces []string

	updateChan chan *ChangeEvent

	mutex  sync.RWMutex
	caches map[string]*cache

	longPoller poller
	client     http.Client
}

// result of query config
type result struct {
	AppID          string            `json:"appId"`
	Cluster        string            `json:"cluster"`
	NamespaceName  string            `json:"namespaceName"`
	Configurations map[string]string `json:"configurations"`
	ReleaseKey     string            `json:"releaseKey"`
}

// NewClient create client from conf
func NewClient(conf *Conf) (*Client, error) {
	client := &Client{
		appID:      conf.AppID,
		cluster:    conf.Cluster,
		ip:         conf.IP,
		namespaces: conf.NameSpaceNames,

		updateChan: make(chan *ChangeEvent),
		caches:     map[string]*cache{},
		longPoller: newLongPoller(conf, longPoolInterval),
		client:     http.Client{Timeout: queryTimeout},
	}

	return client, nil
}

// Start sync config
func (c *Client) Start() error {

	// fetch all config to local first
	if err := c.FetchAllCinfig(); err != nil {
		return err
	}

	go func() {
		c.longPoller.Start(c.handleNamespaceUpdate)
	}()
	return nil
}

func (c *Client) handleNamespaceUpdate(notification *notification) {
	change, err := c.Query(notification.NamespaceName)
	if err != nil || change == nil {
		return
	}
	c.DeliveryChangeEvent(change)
	c.longPoller.UpdateNotification(notification)
}

// Stop sync config
func (c *Client) Stop() error {
	c.longPoller.Stop()
	close(c.updateChan)
	return nil
}

// FetchAllCinfig at first run
func (c *Client) FetchAllCinfig() error {
	for _, namespace := range c.namespaces {
		if _, err := c.Query(namespace); err != nil {
			continue
		}
	}
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
	c.mutex.RUnlock()
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

// Query updated namespace config
func (c *Client) Query(namesapce string) (*ChangeEvent, error) {
	url := configURL(c.ip, c.appID, c.cluster, namesapce)
	bts, err := c.request(url)
	if err != nil || len(bts) == 0 {
		return nil, err
	}
	var result result
	if err := json.Unmarshal(bts, &result); err != nil {
		return nil, err
	}

	return c.HandleResult(&result), nil
}

func (c *Client) request(url string) ([]byte, error) {
	resp, err := c.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return ioutil.ReadAll(resp.Body)
	}
	return nil, nil
}

// DeliveryChangeEvent push change to subscriber
func (c *Client) DeliveryChangeEvent(change *ChangeEvent) {
	select {
	case c.updateChan <- change:
	}
}

// HandleResult generate changes from query result, and update local cache
func (c *Client) HandleResult(result *result) *ChangeEvent {
	var ret = ChangeEvent{
		Namespace: result.NamespaceName,
		Changes:   map[string]*Change{},
	}

	cache := c.mustGetCache(result.NamespaceName)
	kv := cache.All()

	for k, v := range kv {
		if _, ok := result.Configurations[k]; !ok {
			cache.Delete(k)
			ret.Changes[k] = makeDeleteChange(k, v)
		}
	}

	for k, v := range result.Configurations {
		cache.Set(k, v)
		old, ok := kv[k]
		if !ok {
			ret.Changes[k] = makeAddChange(k, v)
			continue
		}
		if old != v {
			ret.Changes[k] = makeModifyChange(k, v)
		}
	}

	if len(ret.Changes) == 0 {
		return nil
	}

	return &ret
}

func makeDeleteChange(key, value string) *Change {
	return &Change{
		ChangeType: DELETE,
		OldValue:   value,
	}
}

func makeModifyChange(key, value string) *Change {
	return &Change{
		ChangeType: MODIFY,
		NewValue:   value,
	}
}

func makeAddChange(key, value string) *Change {
	return &Change{
		ChangeType: ADD,
		NewValue:   value,
	}
}
