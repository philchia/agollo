package agollo

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

var ErrorStatusNotOK = errors.New("http resp code not ok")

// this is a static check
var _ requester = (*httpRequester)(nil)
var _ requester = (*httpSignRequester)(nil)

type requester interface {
	request(url string) ([]byte, error)
}

type httpRequester struct {
	client  *http.Client
	retries int
}

func newHTTPRequester(client *http.Client, retries int) requester {
	return &httpRequester{
		client:  client,
		retries: retries,
	}
}

func (r *httpRequester) request(url string) ([]byte, error) {
	return r.requestWithRetry(url, r.retries)
}

func (r *httpRequester) requestWithRetry(url string, retries int) ([]byte, error) {
	resp, err := r.client.Get(url)
	if err != nil {
		if retries > 0 {
			return r.requestWithRetry(url, retries-1)
		}
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return ioutil.ReadAll(resp.Body)
	}

	// Discard all body if status code is not 200
	_, _ = io.Copy(ioutil.Discard, resp.Body)
	return nil, fmt.Errorf("apollo return http resp code %d", resp.StatusCode)
}

type httpSignRequester struct {
	signature *signature
	client    *http.Client
	retries   int
}

func newHttpSignRequester(signature *signature, client *http.Client, retries int) requester {
	return &httpSignRequester{
		signature: signature,
		client:    client,
		retries:   retries,
	}
}

func (r *httpSignRequester) request(url string) ([]byte, error) {
	return r.requestWithRetry(url, r.retries)
}

func (r *httpSignRequester) requestWithRetry(url string, retries int) ([]byte, error) {

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	timestamp := r.signature.getTimestamp()
	req.Header.Set(signHttpHeaderAuthorization, fmt.Sprintf(
		signAuthorizationFormat,
		r.signature.AppID,
		r.signature.getAuthorization(url, timestamp),
	))
	req.Header.Set(signHttpHeaderTimestamp, timestamp)

	resp, err := r.client.Do(req)
	if err != nil {
		if retries > 0 {
			return r.requestWithRetry(url, retries-1)
		}
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return ioutil.ReadAll(resp.Body)
	}

	// Discard all body if status code is not 200
	_, _ = io.Copy(ioutil.Discard, resp.Body)
	return nil, fmt.Errorf("apollo return http resp code %d", resp.StatusCode)
}
