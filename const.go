package agollo

import (
	"time"
)

const (
	defaultConfName  = "app.properties"
	defaultNamespace = "application"

	signHttpHeaderAuthorization = "Authorization"
	signHttpHeaderTimestamp     = "Timestamp"
	signAuthorizationFormat     = "Apollo %s:%s"
	signDelimiter               = "\n"

	longPollInterval      = time.Second * 2
	longPollTimeout       = time.Second * 90
	queryTimeout          = time.Second * 2
	defaultNotificationID = -1
)
