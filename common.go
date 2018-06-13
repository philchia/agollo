package agollo

import (
	"fmt"
	"net"
	"net/url"
)

func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, a := range addrs {
		if ip4 := toIP4(a); ip4 != nil {
			return ip4.String()
		}
	}
	return ""
}

func toIP4(addr net.Addr) net.IP {
	if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
		return ipnet.IP.To4()
	}
	return nil
}

func notificationURL(conf *Conf, notifications string) string {
	return fmt.Sprintf("http://%s/notifications/v2?appId=%s&cluster=%s&notifications=%s",
		conf.IP,
		url.QueryEscape(conf.AppID),
		url.QueryEscape(conf.Cluster),
		url.QueryEscape(notifications))
}

func configURL(conf *Conf, namespace, releaseKey string) string {
	return fmt.Sprintf("http://%s/configs/%s/%s/%s?releaseKey=%s&ip=%s",
		conf.IP,
		url.QueryEscape(conf.AppID),
		url.QueryEscape(conf.Cluster),
		url.QueryEscape(namespace),
		releaseKey,
		getLocalIP())
}
