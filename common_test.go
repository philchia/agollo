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
	target := notificationURL(
		&Conf{
			IP:      "127.0.0.1:8080",
			AppID:   "SampleApp",
			Cluster: "default",
		}, "")
	_, err := url.Parse(target)
	if err != nil {
		t.Error(err)
	}
}

func TestConfigURL(t *testing.T) {
	target := configURL(
		&Conf{
			IP:      "127.0.0.1:8080",
			AppID:   "SampleApp",
			Cluster: "default",
		}, "application", "")
	_, err := url.Parse(target)
	if err != nil {
		t.Error(err)
	}
}
