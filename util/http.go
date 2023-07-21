package util

import (
	"context"
	log "github.com/sirupsen/logrus"
	"net/http"
	"sync"
)

type HttpClient struct {
	Client *http.Client
}

var httpClientOnce sync.Once
var httpClient *HttpClient

func GetHttpClient() *HttpClient {
	httpClientOnce.Do(func() {
		httpClient = &HttpClient{
			Client: &http.Client{},
		}
	})
	return httpClient
}

func (h *HttpClient) HttpGet(ctx context.Context, url string, headerMap map[string]string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Errorf("HttpGet error, err: %v", err)
		return nil, err
	}
	for key, value := range headerMap {
		req.Header.Add(key, value)
	}
	resp, err := h.Client.Do(req)
	if err != nil {
		log.Errorf("HttpGet error, err: %v", err)
		return nil, err
	}
	return resp, nil
}
