package agollo

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func TestRequest(t *testing.T) {
	request := newHTTPRequester(&http.Client{}, 3)

	serv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte("test"))
	}))

	bts, err := request.request(serv.URL)
	if err != nil {
		t.Error(err)
	}

	if !bytes.Equal(bts, []byte("test")) {
		t.FailNow()
	}

	serv = httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusInternalServerError)
	}))
	bts, err = request.request(serv.URL)
	if err != nil && !strings.Contains(err.Error(), "apollo return http resp code") {
		t.Error(err)
	}

	if len(bts) != 0 {
		t.FailNow()
	}

	serv = httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusInternalServerError)
	}))
	serv.Close()
	_, err = request.request(serv.URL)
	if err == nil {
		t.FailNow()
	}
}

func TestRequestWithSign(t *testing.T) {
	request := newHttpSignRequester(newSignature("appid", "secret"), &http.Client{}, 3)

	serv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		auth := req.Header.Get(signHttpHeaderAuthorization)
		t.Log(auth)
		rw.Write([]byte("test"))
	}))

	bts, err := request.request(serv.URL)
	if err != nil {
		t.Error(err)
	}

	if !bytes.Equal(bts, []byte("test")) {
		t.FailNow()
	}

	serv = httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusInternalServerError)
	}))
	bts, err = request.request(serv.URL)
	if err != nil && !strings.Contains(err.Error(), "apollo return http resp code") {
		t.Error(err)
	}

	if len(bts) != 0 {
		t.FailNow()
	}

	serv = httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusInternalServerError)
	}))
	serv.Close()
	_, err = request.request(serv.URL)
	if err == nil {
		t.FailNow()
	}
}

func TestRequestWithRetry(t *testing.T) {
	var retries int64 = 3
	request := newHTTPRequester(&http.Client{Timeout: time.Millisecond}, int(retries))

	done := make(chan struct{})
	serv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		atomic.AddInt64(&retries, -1)
		<-done
	}))
	_, err := request.request(serv.URL)
	if err == nil {
		t.Errorf("request must be error")
	}

	// wait close
	done <- struct{}{}
	close(done)

	// retry three times, equal query server four times.
	if retries != -1 {
		t.Errorf("must retry three times")
	}
}

func TestRequestWithSignWithRetry(t *testing.T) {
	var retries int64 = 3
	request := newHttpSignRequester(newSignature("appid", "secret"), &http.Client{Timeout: time.Millisecond}, int(retries))

	done := make(chan struct{})
	serv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		atomic.AddInt64(&retries, -1)
		<-done
	}))

	_, err := request.request(serv.URL)
	if err == nil {
		t.Errorf("request must be error")
	}

	// wait close
	done <- struct{}{}
	close(done)

	// retry three times, equal query server four times.
	if retries != -1 {
		t.Errorf("must retry three times")
	}
}
