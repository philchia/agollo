package mockserver

import (
	"context"
	"net/http"
	"sync"
	"time"
)

type mockServer struct {
	once   sync.Once
	server http.Server

	lock          sync.Mutex
	notifications map[string]int
	config        map[string]string
}

func (s *mockServer) notificationHandler(rw http.ResponseWriter, req *http.Request) {

}

func (s *mockServer) configHandler(rw http.ResponseWriter, req *http.Request) {

}

var server mockServer

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
