package agollo

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/philchia/agollo/v4/internal/properties"
)

type Client interface {
	// Start fetch all config to local cache and run period poll to keep update to remote server
	Start() error
	// Stop period poll
	Stop() error

	OnUpdate(func(*ChangeEvent))

	// Get string value for key
	GetString(key string, opts ...OpOption) string
	// GetContent for namespace
	GetContent(opts ...OpOption) string
	// GetPropertiesContent for properties namespace
	GetPropertiesContent(opts ...OpOption) string
	// GetAllKeys return all keys
	GetAllKeys(opts ...OpOption) []string
	// GetReleaseKey return release key for namespace
	GetReleaseKey(opts ...OpOption) string
	// SubscribeToNamespaces will subscribe to new namespace and keep update
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
	conf           *Conf
	skipLocalCache bool

	logger Logger

	caches         *namespaceCache
	releaseKeyRepo *cache

	longPoller poller
	requester  requester

	ctx    context.Context
	cancel context.CancelFunc

	// onUpdateMtx guards for onUpdate
	onUpdateMtx sync.Mutex
	onUpdate    func(event *ChangeEvent)
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
func NewClient(conf *Conf, opts ...ClientOption) Client {
	conf.normalize()
	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: conf.InsecureSkipVerify},
		},
		Timeout: time.Millisecond * time.Duration(conf.SyncTimeout),
	}

	agolloClient := &client{
		conf:           conf,
		logger:         newLogger(),
		caches:         newNamespaceCache(),
		releaseKeyRepo: newCache(),
	}

	agolloClient.requester = newHTTPRequester(httpClient, conf.Retry)
	if conf.AccesskeySecret != "" {
		agolloClient.requester = newHttpSignRequester(
			newSignature(conf.AppID, conf.AccesskeySecret),
			httpClient,
			conf.Retry,
		)
	}

	agolloClient.longPoller = newLongPoller(conf, longPollInterval, agolloClient.handleNamespaceUpdate)
	agolloClient.ctx, agolloClient.cancel = context.WithCancel(context.Background())

	for _, opt := range opts {
		opt(agolloClient)
	}
	return agolloClient
}

// Start sync config
func (c *client) Start() error {
	c.logger.Infof("start agollo client...")
	if !c.skipLocalCache {
		// check cache dir
		if err := c.autoCreateCacheDir(); err != nil {
			c.logger.Errorf("fail to create cache dir: %v", err)
			return err
		}
	}

	// preload all config to local first
	if err := c.preload(); err != nil {
		c.logger.Errorf("fail to preload %v", err)
		return err
	}

	// start fetch update
	go c.longPoller.start()

	return nil
}

// handleNamespaceUpdate sync config for namespace, delivery changes to subscriber
func (c *client) handleNamespaceUpdate(namespace string, notificationId int) error {
	c.logger.Infof("handle namespace %s update", namespace)
	change, err := c.sync(namespace, notificationId)
	if err != nil || change == nil {
		return err
	}

	c.deliveryChangeEvent(change)
	return nil
}

// Stop sync config
func (c *client) Stop() error {
	c.logger.Infof("stop agollo ...")
	c.longPoller.stop()
	c.cancel()
	c.OnUpdate(nil)
	return nil
}

func (c *client) OnUpdate(handler func(*ChangeEvent)) {
	c.onUpdateMtx.Lock()
	defer c.onUpdateMtx.Unlock()

	c.onUpdate = handler
}

// fetchAllConfig fetch from remote, if failed load from local file
func (c *client) preload() error {
	if err := c.longPoller.preload(); err != nil {
		if c.skipLocalCache {
			return err
		}
		c.logger.Infof("preload from remote error : %v", err)
		err2 := c.loadLocal(c.getDumpFileName())
		if err2 != nil {
			return fmt.Errorf("preload from server error [%s], then load local error [%s]", err, err2)
		}
	}
	return nil
}

// loadLocal load caches from local file
func (c *client) loadLocal(name string) error {
	return c.caches.load(name)
}

// dump caches to file
func (c *client) dump(name string) error {
	if c.skipLocalCache {
		return nil
	}
	c.logger.Infof("dump config to local file:%s", name)
	return c.caches.dump(name)
}

func (c *client) mustGetCache(namespace string) *cache {
	return c.caches.mustGetCache(nomalizeNamespace(namespace))
}

// SubscribeToNamespaces fetch namespace config to local and subscribe to updates
func (c *client) SubscribeToNamespaces(namespaces ...string) error {
	c.logger.Infof("subscribe to namespace %#v", namespaces)
	return c.longPoller.addNamespaces(namespaces...)
}

// GetStringValueWithNameSpace get value from given namespace
func (c *client) GetString(key string, opts ...OpOption) string {
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
func (c *client) GetContent(opts ...OpOption) string {
	var op = defaultOperation()
	for _, opt := range opts {
		opt(op)
	}
	if strings.HasSuffix(op.namespace, propertiesSuffix) {
		return c.getPropertiesNamespaceContent(opts...)
	}

	return c.GetString("content", opts...)
}

// GetNameSpaceContent get contents of namespace
func (c *client) GetPropertiesContent(opts ...OpOption) string {
	return c.getPropertiesNamespaceContent(opts...)
}

func (c *client) getPropertiesNamespaceContent(opts ...OpOption) string {
	var op = defaultOperation()
	for _, opt := range opts {
		opt(op)
	}

	doc := properties.New()
	cache := c.mustGetCache(op.namespace).dump()
	for key, value := range cache {
		doc.Set(key, value)
	}

	var buf bytes.Buffer
	_ = properties.Save(doc, &buf)
	return buf.String()
}

// GetAllKeys return all config keys in given namespace
func (c *client) GetAllKeys(opts ...OpOption) []string {
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
func (c *client) sync(namespace string, notificationId int) (*ChangeEvent, error) {
	c.logger.Infof("sync namespace %s with remote config server", namespace)
	url := configURL(c.conf, namespace, notificationId)
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
	c.logger.Infof("delivery update for namespace:%s", change.Namespace)
	c.onUpdateMtx.Lock()
	defer c.onUpdateMtx.Unlock()

	if c.onUpdate != nil {
		c.onUpdate(change)
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
func (c *client) GetReleaseKey(opts ...OpOption) string {
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
	if c.conf.CacheDir == "" {
		return nil
	}

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
