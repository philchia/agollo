package agollo

import (
	"net/url"
	"testing"
)

func TestLocalIp(t *testing.T) {
	ip := getLocalIP()
	if ip == "" {
		t.FailNow()
	}
}

func TestNotificationURL(t *testing.T) {
	target := notificationURL("127.0.0.1:8080", "SampleApp", "default", "")
	_, err := url.Parse(target)
	if err != nil {
		t.Error(err)
	}
}

func TestConfigURL(t *testing.T) {
	target := configURL("127.0.0.1:8080", "SampleApp", "default", "application", "")
	_, err := url.Parse(target)
	if err != nil {
		t.Error(err)
	}
}
