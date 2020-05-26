package agollo

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path"
)

type Client interface {
	Start() error
	Stop() error
	WatchUpdate() <-chan *ChangeEvent

	GetString(key string, opts ...Option) string
	GetContent(opts ...Option) string
	GetAllKeys(opts ...Option) []string
	GetReleaseKey(opts ...Option) string
	SubscribeToNamespaces(namespaces ...string) error
}

type operation struct {
	namespace string
}

func defaultOperation() *operation {
	return &operation{
		namespace: defaultNamespace,
	}
}

// Client for apollo
type client struct {
	conf *Conf

	updateChan chan *ChangeEvent

	caches         *namespaceCache
	releaseKeyRepo *cache

	longPoller poller
	requester  requester

	ctx    context.Context
	cancel context.CancelFunc
}

// result of query config
type result struct {
	// AppID          string            `json:"appId"`
	// Cluster        string            `json:"cluster"`
	NamespaceName  string            `json:"namespaceName"`
	Configurations map[string]string `json:"configurations"`
	ReleaseKey     string            `json:"releaseKey"`
}

// NewClient create client from conf
func NewClient(conf *Conf) Client {
	conf.normalize()
	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: conf.InsecureSkipVerify},
		},
		Timeout: queryTimeout,
	}

	agolloClient := &client{
		conf:           conf,
		caches:         newNamespaceCahce(),
		releaseKeyRepo: newCache(),
	}

	agolloClient.requester = newHTTPRequester(httpClient)
	if conf.AccesskeySecret != "" {
		agolloClient.requester = newHttpSignRequester(
			newSignature(conf.AppID, conf.AccesskeySecret),
			httpClient,
		)
	}

	agolloClient.longPoller = newLongPoller(conf, longPollInterval, agolloClient.handleNamespaceUpdate)
	agolloClient.ctx, agolloClient.cancel = context.WithCancel(context.Background())
	return agolloClient
}

// Start sync config
func (c *client) Start() error {

	// check cache dir
	if err := c.autoCreateCacheDir(); err != nil {
		return err
	}

	// preload all config to local first
	if err := c.preload(); err != nil {
		return err
	}

	// start fetch update
	go c.longPoller.start()

	return nil
}

// handleNamespaceUpdate sync config for namespace, delivery changes to subscriber
func (c *client) handleNamespaceUpdate(namespace string) error {
	change, err := c.sync(namespace)
	if err != nil || change == nil {
		return err
	}

	c.deliveryChangeEvent(change)
	return nil
}

// Stop sync config
func (c *client) Stop() error {
	c.longPoller.stop()
	c.cancel()
	// close(c.updateChan)
	c.updateChan = nil
	return nil
}

// fetchAllCinfig fetch from remote, if failed load from local file
func (c *client) preload() error {
	if err := c.longPoller.preload(); err != nil {
		return c.loadLocal(c.getDumpFileName())
	}
	return nil
}

// loadLocal load caches from local file
func (c *client) loadLocal(name string) error {
	return c.caches.load(name)
}

// dump caches to file
func (c *client) dump(name string) error {
	return c.caches.dump(name)
}

// WatchUpdate get all updates
func (c *client) WatchUpdate() <-chan *ChangeEvent {
	if c.updateChan == nil {
		c.updateChan = make(chan *ChangeEvent, 32)
	}
	return c.updateChan
}

func (c *client) mustGetCache(namespace string) *cache {
	return c.caches.mustGetCache(namespace)
}

// SubscribeToNamespaces fetch namespace config to local and subscribe to updates
func (c *client) SubscribeToNamespaces(namespaces ...string) error {
	return c.longPoller.addNamespaces(namespaces...)
}

// GetStringValueWithNameSpace get value from given namespace
func (c *client) GetString(key string, opts ...Option) string {
	var op = defaultOperation()
	for _, opt := range opts {
		opt(op)
	}

	cache := c.mustGetCache(op.namespace)
	if ret, ok := cache.get(key); ok {
		return ret
	}

	return ""
}

// GetNameSpaceContent get contents of namespace
func (c *client) GetContent(opts ...Option) string {
	return c.GetString("content", opts...)
}

// GetAllKeys return all config keys in given namespace
func (c *client) GetAllKeys(opts ...Option) []string {
	var keys []string
	var op = defaultOperation()
	for _, opt := range opts {
		opt(op)
	}
	cache := c.mustGetCache(op.namespace)
	cache.kv.Range(func(key, value interface{}) bool {
		str, ok := key.(string)
		if ok {
			keys = append(keys, str)
		}
		return true
	})
	return keys
}

// sync namespace config
func (c *client) sync(namespace string) (*ChangeEvent, error) {
	releaseKey := c.GetReleaseKey(WithNamespace(namespace))
	url := configURL(c.conf, namespace, releaseKey)
	bts, err := c.requester.request(url)
	if err != nil || len(bts) == 0 {
		return nil, err
	}
	var result result
	if err := json.Unmarshal(bts, &result); err != nil {
		return nil, err
	}

	return c.handleResult(&result), nil
}

// deliveryChangeEvent push change to subscriber
func (c *client) deliveryChangeEvent(change *ChangeEvent) {
	if c.updateChan == nil {
		return
	}
	select {
	case <-c.ctx.Done():
	case c.updateChan <- change:
	}
}

// handleResult generate changes from query result, and update local cache
func (c *client) handleResult(result *result) *ChangeEvent {
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

	c.setReleaseKey(result.NamespaceName, result.ReleaseKey)

	// dump caches to file
	_ = c.dump(c.getDumpFileName())

	if len(ret.Changes) == 0 {
		return nil
	}

	return &ret
}

func (c *client) getDumpFileName() string {
	cacheDir := c.conf.CacheDir
	fileName := fmt.Sprintf(".%s_%s", c.conf.AppID, c.conf.Cluster)
	return path.Join(cacheDir, fileName)
}

// GetReleaseKey return release key for namespace
func (c *client) GetReleaseKey(opts ...Option) string {
	var op = defaultOperation()
	for _, opt := range opts {
		opt(op)
	}
	releaseKey, _ := c.releaseKeyRepo.get(op.namespace)
	return releaseKey
}

func (c *client) setReleaseKey(namespace, releaseKey string) {
	c.releaseKeyRepo.set(namespace, releaseKey)
}

// autoCreateCacheDir autoCreateCacheDir
func (c *client) autoCreateCacheDir() error {
	fs, err := os.Stat(c.conf.CacheDir)
	if err != nil {
		if os.IsNotExist(err) {
			return os.MkdirAll(c.conf.CacheDir, os.ModePerm)
		}

		return err
	}

	if !fs.IsDir() {
		return fmt.Errorf("conf.CacheDir is not a dir")
	}

	return nil
}
