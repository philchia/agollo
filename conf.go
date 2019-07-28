package agollo

import (
	"encoding/json"
	"fmt"
	"os"
)

// Conf ...
type Conf struct {
	AppID          string   `json:"appId,omitempty"`
	Cluster        string   `json:"cluster,omitempty"`
	NameSpaceNames []string `json:"namespaceNames,omitempty"`
	CacheDir       string   `json:"cacheDir,omitempty"`
	IP             string   `json:"ip,omitempty"`
	MetaAddr       string   `json:"meta_addr"`
}

// NewConf create Conf from file
func NewConf(name string) (*Conf, error) {
	f, err := os.Open(name)
	if err != nil {
		fmt.Println("err:", err)
		return nil, err
	}
	defer f.Close()

	var ret Conf
	if err := json.NewDecoder(f).Decode(&ret); err != nil {
		return nil, err
	}

	return &ret, nil
}
