package agollo

import (
	"encoding/json"
	"os"
)

// Conf ...
type Conf struct {
	AppID              string   `json:"app_id"`
	Cluster            string   `json:"cluster,omitempty"`
	NameSpaceNames     []string `json:"namespace_names,omitempty"`
	CacheDir           string   `json:"cache_dir,omitempty"`
	MetaAddr           string   `json:"meta_addr,omitempty"`
	AccesskeySecret    string   `json:"accesskey_secret,omitempty"`
	InsecureSkipVerify bool     `json:"insecure_skip_verify,omitempty"`
	Retry              int      `json:"retry,omitempty"`        // retry count. 0 means no retry
	SyncTimeout        int64    `json:"sync_timeout,omitempty"` // sync request timeout, unit:ms
	PollTimeout        int64    `json:"poll_timeout,omitempty"` // poll request timeout, unit:ms
}

// NewConf create Conf from file
func NewConf(name string) (*Conf, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var ret Conf
	if err := json.NewDecoder(f).Decode(&ret); err != nil {
		return nil, err
	}

	return &ret, nil
}

func (c *Conf) normalize() {
	if c.Cluster == "" {
		c.Cluster = defaultCluster
	}

	if !strIn(c.NameSpaceNames, defaultNamespace) &&
		!strIn(c.NameSpaceNames, nomalizeNamespace(defaultNamespace)) {
		c.NameSpaceNames = append(c.NameSpaceNames, defaultNamespace)
	}

	if c.SyncTimeout == 0 {
		c.SyncTimeout = 2000
	}

	if c.PollTimeout == 0 {
		c.PollTimeout = 90000
	}
}
