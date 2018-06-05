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
	client, err := NewClient(conf)
	if err != nil {
		return err
	}
	defaultClient = client

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
