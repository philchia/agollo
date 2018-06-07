package agollo

import (
	"encoding/json"
	"io"
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

		caches: map[string]*cache{},

		client: http.Client{Timeout: queryTimeout},
	}

	client.longPoller = newLongPoller(conf, longPoolInterval, client.handleNamespaceUpdate)

	return client, nil
}

// Start sync config
func (c *Client) Start() error {

	// fetch all config to local first
	if err := c.fetchAllCinfig(); err != nil {
		return err
	}

	// start fetch update
	go c.longPoller.start()

	return nil
}

func (c *Client) handleNamespaceUpdate(namespace string) error {
	change, err := c.query(namespace)
	if err != nil {
		return err
	}
	if change == nil {
		return nil
	}
	c.deliveryChangeEvent(change)
	return nil
}

// Stop sync config
func (c *Client) Stop() error {
	c.longPoller.stop()
	close(c.updateChan)
	c.updateChan = nil
	return nil
}

// fetchAllCinfig at first run
func (c *Client) fetchAllCinfig() error {
	return c.longPoller.fire()
}

// WatchUpdate get all updates
func (c *Client) WatchUpdate() <-chan *ChangeEvent {
	if c.updateChan == nil {
		c.updateChan = make(chan *ChangeEvent)
	}
	return c.updateChan
}

func (c *Client) mustGetCache(namespace string) *cache {
	c.mutex.RLock()
	if ret, ok := c.caches[namespace]; ok {
		c.mutex.RUnlock()
		return ret
	}
	c.mutex.RUnlock()

	c.mutex.Lock()
	defer c.mutex.Unlock()

	cache := newCache()
	c.caches[namespace] = cache
	return cache
}

// GetStringValueWithNameSapce get value from given namespace
func (c *Client) GetStringValueWithNameSapce(namespace, key, defaultValue string) string {
	cache := c.mustGetCache(namespace)
	if ret, ok := cache.get(key); ok && ret != "" {
		return ret
	}
	return defaultValue
}

// GetStringValue from default namespace
func (c *Client) GetStringValue(key, defaultValue string) string {
	return c.GetStringValueWithNameSapce(defaultNamespace, key, defaultValue)
}

// query updated namespace config
func (c *Client) query(namesapce string) (*ChangeEvent, error) {
	url := configURL(c.ip, c.appID, c.cluster, namesapce)
	bts, err := c.request(url)
	if err != nil || len(bts) == 0 {
		return nil, err
	}
	var result result
	if err := json.Unmarshal(bts, &result); err != nil {
		return nil, err
	}

	return c.handleResult(&result), nil
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

	// Diacard all body if status code is not 200
	io.Copy(ioutil.Discard, resp.Body)
	return nil, nil
}

// deliveryChangeEvent push change to subscriber
func (c *Client) deliveryChangeEvent(change *ChangeEvent) {
	if c.updateChan == nil {
		return
	}
	c.updateChan <- change
}

// handleResult generate changes from query result, and update local cache
func (c *Client) handleResult(result *result) *ChangeEvent {
	var ret = ChangeEvent{
		Namespace: result.NamespaceName,
		Changes:   map[string]*Change{},
	}

	cache := c.mustGetCache(result.NamespaceName)
	kv := cache.dump()

	for k, v := range kv {
		if _, ok := result.Configurations[k]; !ok {
			cache.delete(k)
			ret.Changes[k] = makeDeleteChange(k, v)
		}
	}

	for k, v := range result.Configurations {
		cache.set(k, v)
		old, ok := kv[k]
		if !ok {
			ret.Changes[k] = makeAddChange(k, v)
			continue
		}
		if old != v {
			ret.Changes[k] = makeModifyChange(k, old, v)
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

func makeModifyChange(key, oldValue, newValue string) *Change {
	return &Change{
		ChangeType: MODIFY,
		OldValue:   oldValue,
		NewValue:   newValue,
	}
}

func makeAddChange(key, value string) *Change {
	return &Change{
		ChangeType: ADD,
		NewValue:   value,
	}
}
