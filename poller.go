package agollo

import (
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"sync"
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
}

// notificationHandler handle namespace update notification
type notificationHandler func(namespace string) error

// longPoller implement poller interface
type longPoller struct {
	appID   string
	cluster string
	ip      string

	pollerInterval time.Duration
	ctx            context.Context
	cancel         context.CancelFunc

	client http.Client

	notifications *notificatonRepo
	handler       notificationHandler

	preloadonce sync.Once
}

// newLongPoller create a Poller
func newLongPoller(conf *Conf, interval time.Duration, handler notificationHandler) poller {
	poller := &longPoller{
		appID:          conf.AppID,
		cluster:        conf.Cluster,
		ip:             conf.IP,
		pollerInterval: interval,
		client:         http.Client{Timeout: longPoolTimeout},
		notifications:  new(notificatonRepo),
		handler:        handler,
	}
	for _, namespace := range conf.NameSpaceNames {
		poller.notifications.setNotificationID(namespace, defaultNotificationID)
	}

	return poller
}

func (p *longPoller) start() {
	go p.watchUpdates()
}

func (p *longPoller) preload() error {
	var err error
	p.preloadonce.Do(func() { err = p.pumpUpdates() })
	return err
}

func (p *longPoller) watchUpdates() {

	p.ctx, p.cancel = context.WithCancel(context.Background())
	defer p.cancel()

	timer := time.NewTimer(p.pollerInterval)
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			p.pumpUpdates()
			timer.Reset(p.pollerInterval)

		case <-p.ctx.Done():
			return
		}
	}
}

func (p *longPoller) stop() {
	p.cancel()
}

func (p *longPoller) updateNotificationConfs(notifications []*notification) {
	for _, noti := range notifications {
		p.updateNotificationConf(noti)
	}
}

func (p *longPoller) updateNotificationConf(notification *notification) {
	p.notifications.setNotificationID(notification.NamespaceName, notification.NotificationID)
}

// pumpUpdates fetch updated namespace, handle updated namespace then update notification id
func (p *longPoller) pumpUpdates() error {
	var ret error

	updates, err := p.poll()
	if err != nil {
		return err
	}

	for _, update := range updates {
		if err := p.handler(update.NamespaceName); err != nil {
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
	url := notificationURL(p.ip, p.appID, p.cluster, notifications)
	bts, err := p.request(url)
	if err != nil || len(bts) == 0 {
		return nil, err
	}
	var ret []*notification
	if err := json.Unmarshal(bts, &ret); err != nil {
		return nil, err
	}
	return ret, nil
}

func (p *longPoller) request(url string) ([]byte, error) {
	resp, err := p.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return ioutil.ReadAll(resp.Body)
	}
	io.Copy(ioutil.Discard, resp.Body)
	return nil, nil
}
