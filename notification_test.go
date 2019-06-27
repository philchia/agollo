package agollo

import "testing"

func TestNotification(t *testing.T) {
	repo := new(notificationRepo)

	repo.setNotificationID("namespace", 1)
	if id, ok := repo.getNotificationID("namespace"); !ok || id != 1 {
		t.FailNow()
	}

	repo.setNotificationID("namespace", 2)
	if id, ok := repo.getNotificationID("namespace"); !ok || id != 2 {
		t.FailNow()
	}

	if id, ok := repo.getNotificationID("null"); ok || id != defaultNotificationID {
		t.FailNow()
	}

	if str := repo.toString(); str == "" {
		t.FailNow()
	}
}
