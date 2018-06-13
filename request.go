package agollo

import (
	"io"
	"io/ioutil"
	"net/http"
)

type requester interface {
	request(url string) ([]byte, error)
}

type httprequester struct {
	client *http.Client
}

func newHTTPRequester(client *http.Client) requester {
	return &httprequester{
		client: client,
	}
}

func (r *httprequester) request(url string) ([]byte, error) {
	resp, err := r.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return ioutil.ReadAll(resp.Body)
	}

	// Diacard all body if status code is not 200
	io.Copy(ioutil.Discard, resp.Body)
	return nil, nil
}
