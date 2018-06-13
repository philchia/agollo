package mockserver

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"sync"
	"time"
)

type mockServer struct {
	once   sync.Once
	server http.Server

	lock          sync.Mutex
	notifications map[string]int
	config        map[string]map[string]string
}

func (s *mockServer) NotificationHandler(rw http.ResponseWriter, req *http.Request) {
	s.lock.Lock()
	defer s.lock.Unlock()

	req.ParseForm()
	query := req.FormValue("notifications")
	var notifications []struct {
		NamespaceName  string `json:"namespaceName,omitempty"`
		NotificationID int    `json:"notificationId,omitempty"`
	}

	if err := json.Unmarshal([]byte(query), &notifications); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	var changes []struct {
		NamespaceName  string `json:"namespaceName,omitempty"`
		NotificationID int    `json:"notificationId,omitempty"`
	}

	for _, noti := range notifications {
		if currentID := s.notifications[noti.NamespaceName]; currentID != noti.NotificationID {
			changes = append(changes, struct {
				NamespaceName  string `json:"namespaceName,omitempty"`
				NotificationID int    `json:"notificationId,omitempty"`
			}{
				NamespaceName:  noti.NamespaceName,
				NotificationID: currentID,
			})
		}
	}

	if len(changes) == 0 {
		rw.WriteHeader(http.StatusNotModified)
		return
	}

	bts, err := json.Marshal(&changes)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.Write(bts)
}

func (s *mockServer) ConfigHandler(rw http.ResponseWriter, req *http.Request) {
	req.ParseForm()

	var appid, cluster, namespace, releaseKey, ip string

	strs := strings.Split(req.RequestURI, "/")

	appid = strs[2]
	cluster = strs[3]
	namespace = strings.Split(strs[4], "?")[0]
	releaseKey = req.FormValue("releaseKey")
	ip = req.FormValue("ip")
	_ = ip
	config := s.Get(namespace)

	var result = struct {
		AppID          string            `json:"appId"`
		Cluster        string            `json:"cluster"`
		NamespaceName  string            `json:"namespaceName"`
		Configurations map[string]string `json:"configurations"`
		ReleaseKey     string            `json:"releaseKey"`
	}{}

	result.AppID = appid
	result.Cluster = cluster
	result.NamespaceName = namespace
	result.Configurations = config
	result.ReleaseKey = releaseKey

	bts, err := json.Marshal(&result)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.Write(bts)
	return
}

var server *mockServer

func (s *mockServer) Set(namespace, key, value string) {
	server.lock.Lock()
	defer server.lock.Unlock()

	notificationID := s.notifications[namespace]
	notificationID++
	s.notifications[namespace] = notificationID

	if kv, ok := s.config[namespace]; ok {
		kv[key] = value
		return
	}
	kv := map[string]string{key: value}
	s.config[namespace] = kv
}

func (s *mockServer) Get(namespace string) map[string]string {
	server.lock.Lock()
	defer server.lock.Unlock()

	return s.config[namespace]
}

func (s *mockServer) GetValue(namespace, key string) string {
	server.lock.Lock()
	defer server.lock.Unlock()

	if kv, ok := s.config[namespace]; ok {
		return kv[key]
	}

	return ""
}

func (s *mockServer) Delete(namespace, key string) {
	server.lock.Lock()
	defer server.lock.Unlock()

	if kv, ok := s.config[namespace]; ok {
		delete(kv, key)
	}

	notificationID := s.notifications[namespace]
	notificationID++
	s.notifications[namespace] = notificationID
}

// Set namespace's key value
func Set(namespace, key, value string) {
	server.Set(namespace, key, value)
}

// Delete namespace's key
func Delete(namespace, key string) {
	server.Delete(namespace, key)
}

// Run mock server
func Run() error {
	initServer()
	return server.server.ListenAndServe()
}

func initServer() {
	server = &mockServer{
		notifications: map[string]int{},
		config:        map[string]map[string]string{},
	}
	server.once.Do(func() {
		mux := http.NewServeMux()
		mux.Handle("/notifications/", http.HandlerFunc(server.NotificationHandler))
		mux.Handle("/configs/", http.HandlerFunc(server.ConfigHandler))
		server.server.Handler = mux
		server.server.Addr = ":8080"
	})
}

// Close mock server
func Close() error {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Second))
	defer cancel()

	return server.server.Shutdown(ctx)
}
