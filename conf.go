package agollo

import (
	"encoding/json"
	"os"
)

// Conf ...
type Conf struct {
	AppID              string   `json:"appId,omitempty"`
	Cluster            string   `json:"cluster,omitempty"`
	NameSpaceNames     []string `json:"namespaceNames,omitempty"`
	CacheDir           string   `json:"cacheDir,omitempty"`
	MetaAddr           string   `json:"meta_addr"`
	AccesskeySecret    string   `json:"accesskey_secret"`
	InsecureSkipVerify bool     `json:"insecure_skip_verify"`
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
	if !strIn(c.NameSpaceNames, defaultNamespace) {
		c.NameSpaceNames = append(c.NameSpaceNames, defaultNamespace)
	}
}
