package agollo

import (
	"testing"
)

func TestWithLogger(t *testing.T) {
	client := new(client)
	WithLogger(nil)(client)
	if client.logger != nil {
		t.Error("err logger is not nil")
		return
	}
}

func TestSkipLocalCache(t *testing.T) {
	client := new(client)
	SkipLocalCache()(client)
	if client.skipLocalCache != true {
		t.Error("skipLocalCache is not true")
		return
	}
}
