package request

import (
	"bytes"
	"context"
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
		},
	}
}

func (r *Requester) Head(ctx context.Context, path string, headers map[string]string, data interface{}) (*http.Response, error) {
	return r.request(ctx, "HEAD", path, headers, data)
}

func (r *Requester) Post(ctx context.Context, path string, headers map[string]string, data interface{}) (*http.Response, error) {
	return r.request(ctx, "POST", path, headers, data)
}

func (r *Requester) Get(ctx context.Context, path string, headers map[string]string, data interface{}) (*http.Response, error) {
	return r.request(ctx, "GET", path, headers, data)
}

func (r *Requester) Patch(ctx context.Context, path string, headers map[string]string, data interface{}) (*http.Response, error) {
	return r.request(ctx, "PATCH", path, headers, data)
}

func (r *Requester) Put(ctx context.Context, path string, headers map[string]string, data interface{}) (*http.Response, error) {
	return r.request(ctx, "PUT", path, headers, data)
}

func (r *Requester) Delete(ctx context.Context, path string, headers map[string]string, data interface{}) (*http.Response, error) {
	return r.request(ctx, "DELETE", path, headers, data)
}

func (r *Requester) request(ctx context.Context, method string, path string, headers map[string]string, data interface{}) (*http.Response, error) {
	payload, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	reader := bytes.NewBuffer(payload)
	req, err := http.NewRequestWithContext(ctx, method, r.ServiceURL+path, reader)
	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		req.Header.Add(k, v)
	}

	req.Header.Add("Content-Type", "application/json")

	return r.Client.Do(req)
}
