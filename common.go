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

func notificationURL(remoteIP, appid, cluster, notifications string) string {
	return fmt.Sprintf("http://%s/notifications/v2?appId=%s&cluster=%s&notifications=%s",
		remoteIP,
		url.QueryEscape(appid),
		url.QueryEscape(cluster),
		url.QueryEscape(notifications))
}

func configURL(remoteIP, appid, cluster, namespace, releaseKey string) string {
	return fmt.Sprintf("http://%s/configs/%s/%s/%s?releaseKey=%s&ip=%s",
		remoteIP,
		url.QueryEscape(appid),
		url.QueryEscape(cluster),
		url.QueryEscape(namespace),
		releaseKey,
		getLocalIP())
}
