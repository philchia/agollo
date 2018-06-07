package agollo

import (
	"testing"
	"time"

	"github.com/philchia/agollo/internal/mockserver"
)

func TestMain(m *testing.M) {
	setup()
	defer teardown()
	m.Run()
}

func setup() {
	go mockserver.Run()
}

func teardown() {
	mockserver.Close()
}

func TestAgolloStart(t *testing.T) {
	if err := StartWithConfFile("./testdata/app.properties"); err != nil {
		t.FailNow()
	}

	defer Stop()

	updates := WatchUpdate()

	select {
	case event := <-updates:
		_ = event
	case <-time.After(time.Millisecond * 10):
	}
}
