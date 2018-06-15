package agollo

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/philchia/agollo/internal/mockserver"
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
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
	if err := Start(); err == nil {
		t.FailNow()
	}

	if err := StartWithConfFile("fake.properties"); err == nil {
		t.FailNow()
	}
	if err := StartWithConfFile("./testdata/app.properties"); err != nil {
		t.FailNow()
	}
	defer Stop()
	defer os.Remove(defaultDumpFile)
	if err := defaultClient.loadLocal(defaultDumpFile); err != nil {
		t.FailNow()
	}

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
