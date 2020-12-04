package agollo

import (
	"os"
	"path"
	"testing"
	"time"

	"github.com/philchia/agollo/v4/internal/mockserver"
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func setup() {
	go mockserver.Run()
	// wait for mock server to run
	time.Sleep(time.Millisecond * 5)
}

func teardown() {
	mockserver.Close()
}

func TestAgolloStart(t *testing.T) {
	conf, err := NewConf("./testdata/app.properties")
	if err != nil {
		t.Error(err)
		return
	}

	if err := Start(conf); err != nil {
		t.Errorf("Start with app.properties should return nil, got :%v", err)
		return
	}

	client := defaultClient.(*client)

	f, err := os.Stat(path.Dir(client.getDumpFileName()))
	if err != nil {
		t.Errorf("dump file dir should exists, got err:%v", err)
		return
	}

	if !f.IsDir() {
		t.Errorf("dump file dir should be a dir, got file")
		return
	}

	defer Stop()

	defer os.Remove(client.getDumpFileName())

	if err := client.loadLocal(client.getDumpFileName()); err != nil {
		t.Errorf("loadLocal should return nil, got: %v", err)
		return
	}

	mockserver.Set("application", "key", "value")
	updates := make(chan struct{}, 1)
	defer close(updates)
	OnUpdate(func(event *ChangeEvent) {
		updates <- struct{}{}
	})

	select {
	case <-updates:
	case <-time.After(time.Millisecond * 30000):
	}

	val := GetString("key")
	if val != "value" {
		t.Errorf("GetStringValue of key should = value, got %v", val)
		return
	}

	keys := GetAllKeys()
	if len(keys) != 1 {
		t.Errorf("GetAllKeys should return 1 key")
		return
	}

	releasekey := GetReleaseKey()
	if releasekey != "" {
		t.Errorf("GetReleaseKey return empty release key")
		return
	}

	mockserver.Set("application", "key", "newvalue")
	select {
	case <-updates:
	case <-time.After(time.Millisecond * 30000):
	}

	val = GetString("key")
	if val != "newvalue" {
		t.Errorf("GetStringValue of key should = newvalue, got %v", val)
		return
	}

	content := GetPropertiesContent()
	if content != "key=newvalue\n" {
		t.Errorf("GetPropertiesContent of application = %s, want %v", content, "key=newvalue\n")
		return
	}

	keys = GetAllKeys()
	if len(keys) != 1 {
		t.Errorf("GetAllKeys should return 1 key")
		return
	}

	mockserver.Delete("application", "key")
	select {
	case <-updates:
	case <-time.After(time.Millisecond * 30000):
	}

	val = GetString("key")
	if val != "" {
		t.Errorf("GetStringValue of key should = defaultValue, got %v", val)
		return
	}

	keys = GetAllKeys()
	if len(keys) != 0 {
		t.Errorf("GetAllKeys should return 0 key")
		return
	}

	mockserver.Set("client.json", "content", `{"name":"agollo"}`)
	select {
	case <-updates:
	case <-time.After(time.Millisecond * 30000):
	}

	val = GetContent(WithNamespace("client.json"))
	if val != `{"name":"agollo"}` {
		t.Errorf(`GetStringValue of client.json content should  = {"name":"agollo"}, got %v`, val)
		return
	}

	if err := SubscribeToNamespaces("new_namespace.json"); err != nil {
		t.Error(err)
		return
	}

	mockserver.Set("new_namespace.json", "content", "1")
	select {
	case <-updates:
	case <-time.After(time.Millisecond * 30000):
	}

	val = GetContent(WithNamespace("new_namespace.json"))
	if val != `1` {
		t.Errorf(`GetStringValueWithNameSpace of new_namespace.json content should  = 1, got %v`, val)
		return
	}
}
