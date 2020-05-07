package agollo

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	neturl "net/url"
	"strconv"
	"time"
)

type signature struct {
	AppID           string `json:"appId"`
	AccesskeySecret string `json:"accesskey_secret"`
}

func newSignature(appId, accesskeySecret string) *signature {
	return &signature{AppID: appId, AccesskeySecret: accesskeySecret}
}

func (s *signature) getAuthorization(url, timestamp string) string {
	sign := timestamp + signDelimiter + s.getURL2PathWithQuery(url)
	key := []byte(s.AccesskeySecret)
	mac := hmac.New(sha1.New, key)
	mac.Write([]byte(sign))
	ss := mac.Sum(nil)
	return base64.StdEncoding.EncodeToString(ss)
}

func (s *signature) getTimestamp() string {
	t := time.Now().UnixNano() / 1e6
	return strconv.Itoa(int(t))
}

func (s *signature) getURL2PathWithQuery(rawurl string) string {
	url, _ := neturl.Parse(rawurl)
	return url.RequestURI()
}
