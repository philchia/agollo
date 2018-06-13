package agollo

import (
	"log"
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
	go func() {
		if err := mockserver.Run(); err != nil {
			log.Fatal(err)
		}
	}()
	// wait for mock server to run
	time.Sleep(time.Millisecond * 10)
}

func teardown() {
	mockserver.Close()
}

func TestAgolloStart(t *testing.T) {
	if err := StartWithConfFile("./testdata/app.properties"); err != nil {
		t.FailNow()
	}

	defer Stop()

	mockserver.Set("application", "key", "value")

	updates := WatchUpdate()

	select {
	case <-updates:
	case <-time.After(time.Millisecond * 30000):
	}

	val := GetStringValue("key", "defaultValue")
	if val != "value" {
		t.FailNow()
	}

	mockserver.Set("application", "key", "newvalue")
	select {
	case <-updates:
	case <-time.After(time.Millisecond * 30000):
	}

	val = defaultClient.GetStringValue("key", "defaultValue")
	if val != "newvalue" {
		t.FailNow()
	}

	mockserver.Delete("application", "key")
	select {
	case <-updates:
	case <-time.After(time.Millisecond * 30000):
	}

	val = GetStringValue("key", "defaultValue")
	if val != "defaultValue" {
		t.FailNow()
	}
}
