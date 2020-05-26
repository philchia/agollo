package agollo

import (
	"sync"
)

var (
	once          sync.Once
	defaultClient Client
)

// Start agollo client with Conf, start will block until fetch all config to local
func Start(conf *Conf) error {
	once.Do(func() { defaultClient = NewClient(conf) })
	return defaultClient.Start()
}

// Stop sync config
func Stop() error {
	return defaultClient.Stop()
}

// WatchUpdate get all updates
func WatchUpdate() <-chan *ChangeEvent {
	return defaultClient.WatchUpdate()
}

// SubscribeToNamespaces fetch namespace config to local and subscribe to updates
func SubscribeToNamespaces(namespaces ...string) error {
	return defaultClient.SubscribeToNamespaces(namespaces...)
}

// GetString get value from given namespace
func GetString(key string, opts ...Option) string {
	return defaultClient.GetString(key, opts...)
}

// GetNameSpaceContent get contents of namespace
func GetContent(opts ...Option) string {
	return defaultClient.GetContent(opts...)
}

// GetAllKeys return all config keys in given namespace
func GetAllKeys(opts ...Option) []string {
	return defaultClient.GetAllKeys(opts...)
}

// GetReleaseKey return release key for namespace
func GetReleaseKey(opts ...Option) string {
	return defaultClient.GetReleaseKey(opts...)
}
