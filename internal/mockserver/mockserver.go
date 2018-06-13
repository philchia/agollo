package mockserver

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

func init() {
	server = &mockServer{
		notifications: map[string]int{},
		config:        map[string]map[string]string{},
	}
}

type mockServer struct {
	once   sync.Once
	server http.Server

	lock          sync.Mutex
	notifications map[string]int
	config        map[string]map[string]string
}

func (s *mockServer) notificationHandler(rw http.ResponseWriter, req *http.Request) {
	s.lock.Lock()
	defer s.lock.Unlock()

	req.ParseForm()
	query := req.FormValue("notifications")
	fmt.Println(query)
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

func (s *mockServer) configHandler(rw http.ResponseWriter, req *http.Request) {

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

// Set namespace's key value
func Set(namespace, key, value string) {
	server.Set(namespace, key, value)
}

// Run mock server
func Run() error {
	initServer()
	return server.server.ListenAndServe()
}

func initServer() {
	server.once.Do(func() {
		mux := http.NewServeMux()
		mux.Handle("/notifications", http.HandlerFunc(server.notificationHandler))
		mux.Handle("/configs", http.HandlerFunc(server.configHandler))
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
