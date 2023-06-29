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
			MetaAddr: "http://127.0.0.1:8080",
			AppID:    "SampleApp",
			Cluster:  "default",
		}, "")
	_, err := url.Parse(target)
	if err != nil {
		t.Error(err)
	}
}

func TestConfigURL(t *testing.T) {
	target := configURL(
		&Conf{
			MetaAddr: "127.0.0.1:8080",
			AppID:    "SampleApp",
			Cluster:  "default",
		}, "application", -1)
	_, err := url.Parse(target)
	if err != nil {
		t.Error(err)
	}
}

func TestStrIn(t *testing.T) {
	cases := []struct {
		slice  []string
		target string
		ok     bool
	}{
		{
			slice:  []string{"a", "b", "c"},
			target: "a",
			ok:     true,
		},

		{
			slice:  []string{"a", "b", "c"},
			target: "d",
			ok:     false,
		},
	}

	for _, c := range cases {
		if strIn(c.slice, c.target) != c.ok {
			t.Fatal("target", c.target, "should be in slice:", c.slice, c.ok)
		}
	}
}
