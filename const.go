package agollo

import (
	"time"
)

const (
	defaultConfName  = "app.properties"
	defaultNamespace = "application"

	longPollInterval      = time.Second * 2
	longPollTimeout       = time.Second * 90
	queryTimeout          = time.Second * 2
	defaultNotificationID = -1
)
