package agollo

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

// this is a static check
var _ Poller = (*longPoller)(nil)

// NotificationHandler handle namespace update notification
type NotificationHandler func(*notification)

// Poller fetch confi updates
type Poller interface {
	Start(handler NotificationHandler)
	Stop()
	UpdateNotification(notification *notification)
}

type longPoller struct {
	appID   string
	cluster string
	ip      string

	pollerInterval time.Duration
	ctx            context.Context
	cancel         context.CancelFunc

	client http.Client

	notifications *notificatonRepo
	handler       NotificationHandler
}

// newLongPoller create a Poller
func newLongPoller(conf *Conf, interval time.Duration) Poller {
	poller := &longPoller{
		appID:          conf.AppID,
		cluster:        conf.Cluster,
		ip:             conf.IP,
		pollerInterval: interval,
		client:         http.Client{Timeout: longPoolTimeout},
		notifications:  new(notificatonRepo),
	}
	for _, namespace := range conf.NameSpaceNames {
		poller.notifications.SetNotificationID(namespace, defaultNotificationID)
	}

	return poller
}

func (p *longPoller) Start(handler NotificationHandler) {
	p.handler = handler
	go p.watchUpdates()
}

func (p *longPoller) UpdateNotification(notification *notification) {
	p.updateNotificationConf(notification)
}

func (p *longPoller) watchUpdates() {

	p.ctx, p.cancel = context.WithCancel(context.Background())
	defer p.cancel()

	tick := time.NewTicker(p.pollerInterval)

	for {
		select {
		case <-tick.C:
			if updates := p.fetch(); len(updates) > 0 {
				for _, update := range updates {
					p.handler(update)
				}
			}
		case <-p.ctx.Done():
			return
		}
	}
}

func (p *longPoller) Stop() {
	p.cancel()
}

func (p *longPoller) updateNotificationConfs(notifications []*notification) {
	for _, noti := range notifications {
		p.updateNotificationConf(noti)
	}
}

func (p *longPoller) updateNotificationConf(notification *notification) {
	p.notifications.SetNotificationID(notification.NamespaceName, notification.NotificationID)
}

func (p *longPoller) fetch() []*notification {
	notifications := p.notifications.AllNotifications()
	url := notificationURL(p.ip, p.appID, p.cluster, notifications)
	bts, err := p.request(url)
	if err != nil || len(bts) == 0 {
		return nil
	}
	var ret []*notification
	if err := json.Unmarshal(bts, &ret); err != nil {
		return nil
	}
	return ret
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
	return nil, nil
}
