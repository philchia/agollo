package agollo

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"net/http"
	"sync/atomic"
	"time"
)

// this is a static check
var _ poller = (*longPoller)(nil)

// poller fetch confi updates
type poller interface {
	// start poll updates
	start()
	// preload fetch all config to local cache, and update all notifications
	preload() error
	// stop poll updates
	stop()
	// addNamespaces add new namespace and pump config data
	addNamespaces(namespaces ...string) error
}

// notificationHandler handle namespace update notification
type notificationHandler func(namespace string, notificationId int) error

// longPoller implement poller interface
type longPoller struct {
	conf *Conf

	pollerInterval time.Duration
	ctx            context.Context
	cancel         context.CancelFunc
	version        uint64
	requester      requester

	notifications *notificationRepo
	handler       notificationHandler
}

// newLongPoller create a Poller
func newLongPoller(conf *Conf, interval time.Duration, handler notificationHandler) poller {
	httpClient := &http.Client{
		Timeout: time.Millisecond * time.Duration(conf.PollTimeout),
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: conf.InsecureSkipVerify},
		},
	}
	requester := newHTTPRequester(httpClient, conf.Retry)
	if conf.AccesskeySecret != "" {
		requester = newHttpSignRequester(
			newSignature(conf.AppID, conf.AccesskeySecret),
			httpClient,
			conf.Retry,
		)
	}

	poller := &longPoller{
		conf:           conf,
		pollerInterval: interval,
		requester:      requester,
		notifications:  new(notificationRepo),
		handler:        handler,
	}

	poller.ctx, poller.cancel = context.WithCancel(context.Background())

	for _, namespace := range conf.NameSpaceNames {
		poller.notifications.setNotificationID(namespace, defaultNotificationID)
	}

	return poller
}

func (p *longPoller) start() {
	go p.watchUpdates()
}

func (p *longPoller) preload() error {
	return p.pumpUpdates()
}

// addNamespaces subscribe to new namespaces and pull all config data to local
func (p *longPoller) addNamespaces(namespaces ...string) error {
	var update bool
	for _, namespace := range namespaces {
		if p.notifications.addNotificationID(namespace, defaultNotificationID) {
			update = true
		}
	}

	if update {
		return p.pumpUpdates()
	}

	return nil
}

func (p *longPoller) watchUpdates() {
	timer := time.NewTimer(p.pollerInterval)
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			_ = p.pumpUpdates()
			timer.Reset(p.pollerInterval)

		case <-p.ctx.Done():
			return
		}
	}
}

func (p *longPoller) stop() {
	p.cancel()
}

func (p *longPoller) updateNotificationConf(notification *notification) {
	p.notifications.setNotificationID(notification.NamespaceName, notification.NotificationID)
}

// pumpUpdates fetch updated namespace, handle updated namespace then update notification id
func (p *longPoller) pumpUpdates() error {
	// serialize pumpUpdates request

	version := atomic.AddUint64(&p.version, 1)

	var ret error

	updates, err := p.poll()
	if err != nil {
		return err
	}

	if atomic.LoadUint64(&p.version) != version {
		return nil
	}

	for _, update := range updates {
		if err := p.handler(update.NamespaceName, update.NotificationID); err != nil {
			ret = err
			continue
		}
		p.updateNotificationConf(update)
	}
	return ret
}

// poll until a update or timeout
func (p *longPoller) poll() ([]*notification, error) {
	notifications := p.notifications.toString()
	url := notificationURL(p.conf, notifications)
	bts, err := p.requester.request(url)
	if err != nil || len(bts) == 0 {
		return nil, err
	}
	var ret []*notification
	if err := json.Unmarshal(bts, &ret); err != nil {
		return nil, err
	}
	return ret, nil
}
