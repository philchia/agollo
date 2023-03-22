package agollo

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
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
