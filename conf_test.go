package agollo

import "testing"

func TestNewConf(t *testing.T) {
	conf, err := NewConf("./testdata/app.properties")
	if err != nil {
		t.Error(err)
	}
	_ = conf
}
