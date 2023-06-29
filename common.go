package agollo

import (
	"encoding/json"
	"fmt"
	"net"
	"net/url"
	"strings"
)

type messages struct {
	Details map[string]interface{} `json:"details,omitempty"`
}

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

func httpurl(ipOrAddr string) string {
	if strings.HasPrefix(ipOrAddr, "http://") || strings.HasPrefix(ipOrAddr, "https://") {
		return ipOrAddr
	}

	return fmt.Sprintf("http://%s", ipOrAddr)
}

func nomalizeNamespace(namespace string) string {
	return strings.TrimSuffix(namespace, propertiesSuffix)
}

func notificationURL(conf *Conf, notifications string) string {
	var addr = conf.MetaAddr
	return fmt.Sprintf("%s/notifications/v2?appId=%s&cluster=%s&notifications=%s",
		httpurl(addr),
		url.QueryEscape(conf.AppID),
		url.QueryEscape(conf.Cluster),
		url.QueryEscape(notifications))
}

func configURL(conf *Conf, namespace string, notificationId int) string {
	var addr = conf.MetaAddr
	message := getMessages(conf, namespace, notificationId)
	return fmt.Sprintf("%s/configs/%s/%s/%s?releaseKey=&ip=%s&messages=%s",
		httpurl(addr),
		url.QueryEscape(conf.AppID),
		url.QueryEscape(conf.Cluster),
		url.QueryEscape(nomalizeNamespace(namespace)),
		getLocalIP(),
		url.QueryEscape(message))
}

func getMessages(conf *Conf, namespace string, notificationId int) string {
	key := fmt.Sprintf("%s+%s+%s", conf.AppID, conf.Cluster, nomalizeNamespace(namespace))
	d := make(map[string]interface{})
	d[key] = notificationId
	m := &messages{}
	m.Details = d
	message, err := json.Marshal(m)
	if err != nil {
		fmt.Println("get messaget err:", err)
		return ""
	}
	return string(message)
}

func strIn(slice []string, target string) bool {
	for _, v := range slice {
		if v == target {
			return true
		}
	}

	return false
}
