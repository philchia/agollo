package agollo

var (
	defaultClient *Client
)

// Start agollo
func Start() error {
	return StartWithConfFile(defaultConfName)
}

// StartWithConfFile run agollo with conf file
func StartWithConfFile(name string) error {
	conf, err := NewConf(name)
	if err != nil {
		return err
	}
	return StartWithConf(conf)
}

// StartWithConf run agollo with Conf
func StartWithConf(conf *Conf) error {
	defaultClient = NewClient(conf)

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

// GetStringValueWithNameSpace get value from given namespace
func GetStringValueWithNameSpace(namespace, key, defaultValue string) string {
	return defaultClient.GetStringValueWithNameSpace(namespace, key, defaultValue)
}

// GetStringValue from default namespace
func GetStringValue(key, defaultValue string) string {
	return GetStringValueWithNameSpace(defaultNamespace, key, defaultValue)
}

// GetNameSpaceContent get contents of namespace
func GetNameSpaceContent(namespace, defaultValue string) string {
	return defaultClient.GetNameSpaceContent(namespace, defaultValue)
}

// GetAllKeys return all config keys in given namespace
func GetAllKeys(namespace string) []string {
	return defaultClient.GetAllKeys(namespace)
}

// GetReleaseKey return release key for namespace
func GetReleaseKey(namespace string) string {
	return defaultClient.GetReleaseKey(namespace)
}
