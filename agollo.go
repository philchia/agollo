package agollo

import (
	"sync"
)

var (
	once          sync.Once
	defaultClient Client
)

// Start agollo client with Conf, start will block until fetch all config to local
func Start(conf *Conf, opts ...ClientOption) error {
	once.Do(func() { defaultClient = NewClient(conf, opts...) })
	return defaultClient.Start()
}

// Stop sync config
func Stop() error {
	return defaultClient.Stop()
}

// OnUpdate get all updates
func OnUpdate(handler func(*ChangeEvent)) {
	defaultClient.OnUpdate(handler)
}

// SubscribeToNamespaces fetch namespace config to local and subscribe to updates
func SubscribeToNamespaces(namespaces ...string) error {
	return defaultClient.SubscribeToNamespaces(namespaces...)
}

// GetString get value from given namespace
func GetString(key string, opts ...OpOption) string {
	return defaultClient.GetString(key, opts...)
}

// GetNameSpaceContent get contents of namespace
func GetContent(opts ...OpOption) string {
	return defaultClient.GetContent(opts...)
}

// GetPropertiesContent for properties namespace
func GetPropertiesContent(opts ...OpOption) string {
	return defaultClient.GetPropertiesContent(opts...)
}

// GetAllKeys return all config keys in given namespace
func GetAllKeys(opts ...OpOption) []string {
	return defaultClient.GetAllKeys(opts...)
}

// GetReleaseKey return release key for namespace
func GetReleaseKey(opts ...OpOption) string {
	return defaultClient.GetReleaseKey(opts...)
}
