package request

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"
)

type Requester struct {
	ServiceURL string
	Client     *http.Client
}

func NewRequester(host string) Requester {
	return Requester{
		ServiceURL: host,
		Client: &http.Client{
			Timeout: 20 * time.Second,
			Transport: &http.Transport{
				MaxIdleConnsPerHost: 10,
			},
		}}
}

func (r *Requester) Head(path string, headers map[string]string, data interface{}) (*http.Response, error) {
	return r.request("HEAD", path, headers, data)
}

func (r *Requester) Post(path string, headers map[string]string, data interface{}) (*http.Response, error) {
	return r.request("POST", path, headers, data)
}

func (r *Requester) Get(path string, headers map[string]string, data interface{}) (*http.Response, error) {
	return r.request("GET", path, headers, data)
}

func (r *Requester) Patch(path string, headers map[string]string, data interface{}) (*http.Response, error) {
	return r.request("PATCH", path, headers, data)
}

func (r *Requester) Put(path string, headers map[string]string, data interface{}) (*http.Response, error) {
	return r.request("PUT", path, headers, data)
}

func (r *Requester) Delete(path string, headers map[string]string, data interface{}) (*http.Response, error) {
	return r.request("DELETE", path, headers, data)
}

func (r *Requester) request(method string, path string, headers map[string]string, data interface{}) (*http.Response, error) {
	payload, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	reader := bytes.NewBuffer(payload)
	req, err := http.NewRequest(method, r.ServiceURL+path, reader)
	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		req.Header.Add(k, v)
	}

	req.Header.Add("Content-Type", "application/json")

	return r.Client.Do(req)
}
