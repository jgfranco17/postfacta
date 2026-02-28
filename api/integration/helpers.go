package integration_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"time"
)

var httpClient = &http.Client{Timeout: 3 * time.Second}

func toJSONReader(value any) (io.Reader, error) {
	payload, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(payload), nil
}

type HttpTestRunner struct {
	BaseURL string
}

func (r *HttpTestRunner) Do(method string, endpoint string, body io.Reader, headers map[string]string) (*http.Response, []byte, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), 2*time.Second)
	defer cancel()

	fullPath, err := url.JoinPath(r.BaseURL, endpoint)
	if err != nil {
		return nil, nil, err
	}
	request, err := http.NewRequestWithContext(ctx, method, fullPath, body)
	if err != nil {
		return nil, nil, err
	}

	for key, value := range headers {
		request.Header.Set(key, value)
	}

	response, err := httpClient.Do(request)
	if err != nil {
		return nil, nil, err
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return response, nil, err
	}

	return response, responseBody, nil
}
