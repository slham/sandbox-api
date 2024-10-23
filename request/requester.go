package request

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log/slog"
	"math"
	"net/http"
	"time"
)

const (
	DEFAULT_REQUEST_TIMEOUT               = 20
	DEFAULT_MAX_IDLE_CONNECTIONS_PER_HOST = 10
	DEFAULT_RETRY_COUNT                   = 3
)

type Requester struct {
	ServiceURL string
	Client     *http.Client
}

type retryableTransport struct {
	transport http.RoundTripper
}

func (t *retryableTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Clone the request body
	var bodyBytes []byte
	if req.Body != nil {
		bodyBytes, err := io.ReadAll(req.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read body. %w", err)
		}
		req.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	}
	// Send the request
	resp, err := t.transport.RoundTrip(req)
	slog.Debug("failed round trip", "err", err)
	// Retry logic
	retries := 0
	for shouldRetry(err, resp) && retries < DEFAULT_RETRY_COUNT {
		// Wait for the specified backoff period
		time.Sleep(backoff(retries))
		// We're going to retry, consume any response to reuse the connection.
		drainBody(resp)
		// Clone the request body again
		if req.Body != nil {
			req.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
		}
		// Retry the request
		resp, err = t.transport.RoundTrip(req)
		slog.Debug("failed round trip", "err", err)
		retries++
	}
	// Return the response
	return resp, err
}

func NewRequester(host string) Requester {
	return Requester{
		ServiceURL: host,
		Client:     NewRetryableClient(),
	}
}

func NewRetryableClient() *http.Client {
	transport := &retryableTransport{
		transport: &http.Transport{
			MaxIdleConnsPerHost: DEFAULT_MAX_IDLE_CONNECTIONS_PER_HOST,
		},
	}

	return &http.Client{
		Timeout:   DEFAULT_REQUEST_TIMEOUT * time.Second,
		Transport: transport,
	}
}

func backoff(retries int) time.Duration {
	return time.Duration(math.Pow(2, float64(retries))) * time.Second
}

func shouldRetry(err error, resp *http.Response) bool {
	if err != nil {
		slog.Debug("failed http request", "err", err)
		return true
	}

	return resp.StatusCode == http.StatusBadGateway ||
		resp.StatusCode == http.StatusServiceUnavailable ||
		resp.StatusCode == http.StatusGatewayTimeout
}

func drainBody(resp *http.Response) {
	if resp.Body != nil {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}
}

func (r *Requester) Head(ctx context.Context, path string, headers map[string]string, data any) (*http.Response, error) {
	return r.request(ctx, http.MethodHead, path, headers, data)
}

func (r *Requester) Post(ctx context.Context, path string, headers map[string]string, data any) (*http.Response, error) {
	return r.request(ctx, http.MethodPost, path, headers, data)
}

func (r *Requester) Get(ctx context.Context, path string, headers map[string]string, data any) (*http.Response, error) {
	return r.request(ctx, http.MethodGet, path, headers, data)
}

func (r *Requester) Patch(ctx context.Context, path string, headers map[string]string, data any) (*http.Response, error) {
	return r.request(ctx, http.MethodPatch, path, headers, data)
}

func (r *Requester) Put(ctx context.Context, path string, headers map[string]string, data any) (*http.Response, error) {
	return r.request(ctx, http.MethodPut, path, headers, data)
}

func (r *Requester) Delete(ctx context.Context, path string, headers map[string]string, data any) (*http.Response, error) {
	return r.request(ctx, http.MethodDelete, path, headers, data)
}

func (r *Requester) request(ctx context.Context, method string, path string, headers map[string]string, data any) (*http.Response, error) {
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
